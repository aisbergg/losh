package dgraph

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"strconv"
	"strings"
	"time"

	"losh/internal/infra/dgraph/dgclient"
	"losh/internal/lib/net"
	"losh/internal/lib/net/request"
	"losh/internal/lib/util/pathutil"

	"github.com/aisbergg/go-errors/pkg/errors"

	gql "github.com/Yamashou/gqlgenc/clientv2"
	"github.com/aisbergg/go-copier/pkg/copier"
	"go.uber.org/zap"
)

// DgraphRepository is the license repository that uses Dgraph as the backend.
type DgraphRepository struct {
	// address is the address of the Dgraph GraphQL endpoint.
	address string

	// used internally
	httpClient *http.Client
	client     dgclient.DgraphGraphQLClient
	log        *zap.SugaredLogger

	dataCopier *copier.Copier

	// Type converters are used in conjunction with copier to convert between
	// Dgraph input/output models and regular data models. Using copier saves me
	// a lot of manual coding and duplicated code.
	convertersForGet  []copier.TypeConverter
	convertersForSave []copier.TypeConverter
}

// NewDgraphRepository creates a new DgraphRepository.
func NewDgraphRepository(dbConfig Config) (*DgraphRepository, error) {
	log := logging.NewLogger("repo-dgraph")
	timeout := 30 * time.Second

	// create HTTP client
	httpClient := &http.Client{Timeout: timeout}
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
		SetMaxWaitTime(timeout)
	client := dgclient.NewClient(graphQLRequester)
	dgraphRepo := &DgraphRepository{
		client:     client,
		httpClient: httpClient,
		address:    address,
		log:        log,
		dataCopier: copier.New(copier.Options{AutoConvert: true, IgnoreEmpty: true}),
	}
	dgraphRepo.initializeConverters()
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
