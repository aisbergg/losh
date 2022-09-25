// Copyright 2022 André Lehmann
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dgraph

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	goerrors "errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"losh/internal/infra/dgraph/dgclient"
	"losh/internal/lib/log"
	"losh/internal/lib/net"
	"losh/internal/lib/net/request"
	"losh/internal/lib/util/pathutil"

	"github.com/aisbergg/go-copier/pkg/copier"
	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-retry/pkg/retry"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/fenos/dqlx"
	"google.golang.org/grpc"

	gql "github.com/Yamashou/gqlgenc/clientv2"
	"go.uber.org/zap"
)

// DgraphRepository is the license repository that uses Dgraph as the backend.
type DgraphRepository struct {
	// address is the address of the Dgraph GraphQL endpoint.
	address string

	// used internally
	timeout    time.Duration
	httpClient *http.Client
	requester  *request.GraphQLRequester
	client     *dgclient.Client
	copier     *copier.Copier
	dqlCopier  *copier.Copier
	// copierFull *copier.Copier
	log *zap.SugaredLogger

	// dqlxClient
	dqlxClient   dqlx.DB
	dgraphClient *dgo.Dgraph
}

// NewDgraphRepository creates a new DgraphRepository.
func NewDgraphRepository(dbConfig Config) (*DgraphRepository, error) {
	log := log.NewLogger("repo-dgraph")

	// create HTTP client
	httpClient := &http.Client{Timeout: dbConfig.Timeout}
	if dbConfig.TLS.Enabled {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: !dbConfig.TLS.Verify,
		}
		if dbConfig.TLS.Certificate != "" {
			if dbConfig.TLS.Key != "" {
				return nil, errors.New("missing TLS client key")
			}
			tlsCrtPth, err := pathutil.GetValidFilePath(dbConfig.TLS.Certificate)
			if err != nil {
				return nil, errors.Wrap(err, "failed to load TLS client certificate")
			}
			tlsKeyPth, err := pathutil.GetValidFilePath(dbConfig.TLS.Key)
			if err != nil {
				return nil, errors.Wrap(err, "failed to load TLS client key")
			}

			cert, err := tls.LoadX509KeyPair(tlsCrtPth.String(), tlsKeyPth.String())
			if err != nil {
				errors.Wrap(err, "failed to load client-key pair")
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}
		if len(dbConfig.TLS.CACertificates) > 0 {
			caCertPool := x509.NewCertPool()
			tlsConfig.RootCAs = caCertPool
			for _, crt := range dbConfig.TLS.CACertificates {
				caCrtPth, err := pathutil.GetValidFilePath(crt)
				if err != nil {
					return nil, errors.Wrap(err, "failed to load CA certificate")
				}
				crtCnt, err := caCrtPth.ReadFile()
				if err != nil {
					return nil, errors.Wrap(err, "failed to read CA certificate")
				}
				caCertPool.AppendCertsFromPEM(crtCnt)
			}
		}
		httpClient.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	scheme := "http"
	if dbConfig.TLS.Enabled {
		scheme = "https"
	}
	address := scheme + "://" + strings.TrimSpace(dbConfig.Host) + ":" + strconv.FormatUint(dbConfig.Port, 10) + "/graphql"
	gqlClient := gql.NewClient(httpClient, address)
	graphQLRequester := request.NewGraphQLRequester(gqlClient).
		SetRetryCount(5).
		SetMaxWaitTime(dbConfig.Timeout)
	// client := dgclient.NewClient(graphQLRequester)
	client := &dgclient.Client{
		Requester: graphQLRequester,
	}
	// copierFull := copier.New(copier.Options{
	// 	AutoConvert:    true,
	// 	CopyUnexported: true,
	// 	IgnoreEmpty:    false,
	// 	Converters:     copierConverters(),
	// })
	dqlCopier := copier.New(copier.Options{
		AutoConvert:    true,
		CopyUnexported: true,
		IgnoreEmpty:    false,
		Converters:     dqlCopierConverters(),
		ConsiderTags:   []string{"dql"},
	})
	copier := copier.New(copier.Options{
		AutoConvert:    true,
		CopyUnexported: true,
		IgnoreEmpty:    false,
		Converters:     copierConverters(),
		ConsiderTags:   []string{"json"},
	})

	// -------------------------------------------------------------------------
	// dqlxClient
	// dgraphClient := api.DgraphClient, len(addresses))
	// TODO: TLS stuff
	dial, err := grpc.Dial(strings.TrimSpace(dbConfig.Host)+":9080", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	dgraphClient := dgo.NewDgraphClient(api.NewDgraphClient(dial))
	dqlxClient := dqlx.FromClient(dgraphClient)

	// TODO: access to cluster
	// func Connect(addresses ...string) (DB, error) {
	// 	clients := make([]api.DgraphClient, len(addresses))

	// 	for index, address := range addresses {
	// 		dial, err := grpc.Dial(address, grpc.WithInsecure())
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		clients[index] = api.NewDgraphClient(dial)
	// 	}

	// 	dgraph := dgo.NewDgraphClient(clients...)

	// 	return FromClient(dgraph), nil
	// }

	// -------------------------------------------------------------------------
	dgraphRepo := &DgraphRepository{
		timeout:    dbConfig.Timeout,
		httpClient: httpClient,
		requester:  graphQLRequester,
		client:     client,
		copier:     copier,
		dqlCopier:  dqlCopier,
		// copierFull: copierFull,
		address: address,
		log:     log,

		dqlxClient:   dqlxClient,
		dgraphClient: dgraphClient,
	}
	return dgraphRepo, nil
}

// IsReachable indicates whether the Dgraph repository is reachable.
func (dr *DgraphRepository) IsReachable() bool {
	// can a HTTP connection be established?
	canCon := net.CheckHTTPConnection(dr.httpClient, dr.address)
	if !canCon {
		return false
	}

	// can a GraphQL query be issued?
	_, err := dr.client.TestConnection(context.Background())
	if err != nil {
		return false
	}

	return true
}

// WaitUntilReachable waits until the Dgraph repository is reachable. It returns
// an error if the database could not be reached after the configured timeout.
func (dr *DgraphRepository) WaitUntilReachable() error {
	b := retry.NewConstant(5 * time.Second)
	b = request.WithRetryable(b)
	if dr.timeout > 0 {
		b = request.WithDelayLimit(dr.timeout, b)
	}
	b = request.WithHook(func(delay time.Duration, err error) (time.Duration, error) {
		dr.log.Errorf("failed to reach database at %s, waiting %s and try again: %s", dr.address, delay, err.Error())
		return delay, nil
	}, b)

	// create wrapper for request execution
	retryFunc := func(ctx context.Context) (err error) {
		// check if the request was cancelled
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// can a HTTP connection be established?
		canCon := net.CheckHTTPConnection(dr.httpClient, dr.address)
		if !canCon {
			return request.NewRetryableError(goerrors.New("no network connection"))
		}

		// can a GraphQL query be issued?
		_, err = dr.client.TestConnection(context.Background())
		if err != nil {
			return request.NewRetryableError(goerrors.New("GraphQL request failed"))
		}

		// stop on success
		return nil
	}

	// execute the request with retries
	err := retry.Do(context.Background(), b, retryFunc)
	if err != nil {
		return err
	}

	return nil
}
