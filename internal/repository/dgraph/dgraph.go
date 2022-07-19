package dgraph

import (
	"losh/crawler/logging"
	"losh/internal/net"
	"losh/internal/net/request"
	"losh/internal/repository/dgraph/dgclient"
	"net/http"
	"time"

	"github.com/aisbergg/go-errors/pkg/errors"

	gql "github.com/Yamashou/gqlgenc/clientv2"
	"github.com/jinzhu/copier"
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

	// Type converters are used in conjunction with copier to convert between
	// Dgraph input/output models and regular data models. Using copier saves me
	// a lot of manual coding and duplicated code.
	convertersForGet  []copier.TypeConverter
	convertersForSave []copier.TypeConverter
}

// NewDgraphRepository creates a new DgraphRepository.
func NewDgraphRepository(address string) *DgraphRepository {
	log := logging.NewLogger("repo-dgraph")
	timeout := 30 * time.Second
	httpClient := &http.Client{Timeout: timeout}
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
	}
	dgraphRepo.initializeConverters()
	return dgraphRepo
}

// IsReachable indicates whether the Dgraph repository is reachable.
func (dr *DgraphRepository) IsReachable() bool {
	return net.CheckHTTPConnection(dr.httpClient, dr.address)
}
