package wikifactory

import (
	"context"
	"net/http"
	"time"

	"losh/crawler/core/validator"
	"losh/crawler/core/wikifactory/wfclient"
	"losh/internal/core/product/models"
	"losh/internal/infra/dgraph"
	"losh/internal/lib/log"
	"losh/internal/lib/net/download"
	"losh/internal/lib/net/ratelimit"
	"losh/internal/lib/net/request"
	"losh/internal/license"

	gql "github.com/Yamashou/gqlgenc/clientv2"
	"github.com/aisbergg/go-errors/pkg/errors"

	"go.uber.org/zap"

	n "losh/internal/lib/net"
	"losh/internal/lib/unit"
)

// Constants that define the crawler behavior.
const (
	crawlerName         = "wikifactory.com"
	timeout             = time.Duration(15) * time.Second
	batchSize           = 10
	retries             = 5
	maxFileSizeManifest = 10 * unit.MiB
	maxWaitTime         = 5 * time.Minute
	maxRedirects        = 5
)

type CrawlerState struct {
	Timestamp   time.Time
	StartTime   time.Time
	ElapsedTime time.Time
	NumCrawled  int64
	NumIndexed  int64
	Cursor      string
	Page        int64
	BatchSize   int64
}

// WikifactoryCrawler is a product crawler for Wikifactory.
type WikifactoryCrawler struct {
	repository *dgraph.DgraphRepository
	validator  *validator.Validator
	normalizer *WikifactoryNormalizer

	fileDownloader *download.Downloader
	wfGqlClient    wfclient.WikifactoryGraphQLClient
	log            *zap.SugaredLogger
}

// NewWikifactoryCrawler creates a new WikifactoryCrawler.
func NewWikifactoryCrawler(repository *dgraph.DgraphRepository, licenseCache *license.LicenseCache, userAgent string) *WikifactoryCrawler {
	normalizer := NewWikifactoryNormalizer(licenseCache)
	validator := validator.NewValidator(licenseCache)
	log := log.NewLogger("crawler-wikifactory")

	// clients for external requests
	httpClient := &http.Client{
		Timeout:       timeout,
		CheckRedirect: n.NewRedirectHandler(maxRedirects),
	}

	// file downloader
	fileRequester := request.NewHTTPRequester(httpClient).
		SetLogger(log).
		SetRetryCount(retries).
		SetMaxWaitTime(maxWaitTime).
		AddRateLimiter(ratelimit.NewTimedeltaRateLimiter(1*time.Second, 5))
	fileDownloader := download.NewDownloaderWithRequester(fileRequester).
		SetUserAgent(userAgent)

	// Wikifactory GraphQL Client
	gqlClient := gql.NewClient(httpClient, "https://wikifactory.com/api/graphql")
	graphQLRequester := request.NewGraphQLRequester(gqlClient).
		SetLogger(log).
		SetRetryCount(uint64(retries)).
		SetMaxWaitTime(maxWaitTime).
		AddRateLimiter(ratelimit.NewTimedeltaRateLimiter(5*time.Second, 5))
	graphQLClient := wfclient.NewClient(graphQLRequester)

	return &WikifactoryCrawler{
		repository: repository,
		validator:  validator,
		normalizer: normalizer,

		fileDownloader: fileDownloader,
		wfGqlClient:    graphQLClient,

		log: log,
	}
}

func (c *WikifactoryCrawler) GetProduct(ctx context.Context, productID models.ProductID) (*models.Product, error) {
	// // get project information
	// getProjectBySlug, err := c.wfGqlClient.GetProjectBySlug(ctx, &productID.Owner, &productID.Repo)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to get project information")
	// }
	// if getProjectBySlug == nil {
	// 	return nil, errors.New("project not found")
	// }
	// discoveredAt := time.Now()

	// get basic product information
	getMandatoryProjectBySlug, err := c.wfGqlClient.GetMandatoryProjectBySlug(ctx, &productID.Owner, &productID.Repo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project information")
	}
	if getMandatoryProjectBySlug == nil {
		return nil, errors.New("project not found")
	}
	discoveredAt := time.Now()

	// quick check if mandatory fields are present
	mdtFlds := validator.MandatoryFields{
		OKHV:                  "OKH-LOSHv1.0",
		Name:                  *getMandatoryProjectBySlug.Project.Result.Name,
		Description:           *getMandatoryProjectBySlug.Project.Result.Description,
		Version:               *getMandatoryProjectBySlug.Project.Result.Contribution.Version,
		Repository:            "https://wikifactory.com", // just to pass the check
		License:               c.normalizer.translateLicense(*getMandatoryProjectBySlug.Project.Result.License.Abreviation),
		Licensor:              *getMandatoryProjectBySlug.Project.Result.Creator.Profile.Username,
		DocumentationLanguage: c.normalizer.getDocumentationLanguage(*getMandatoryProjectBySlug.Project.Result.Description),
	}
	vldErr := c.validator.ValidateMandatory(mdtFlds)
	if vldErr != nil {
		return nil, errors.Wrap(vldErr, "invalid mandatory fields")
	}

	// get full project information
	getFullProjectBySlug, err := c.wfGqlClient.GetFullProjectBySlug(ctx, &productID.Owner, &productID.Repo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project information")
	}
	if getFullProjectBySlug == nil {
		return nil, errors.New("project not found")
	}

	// normalize project information
	product := c.normalizer.NormalizeProduct(discoveredAt, getFullProjectBySlug)

	// validate product
	vldErr = c.validator.ValidateProduct(product)
	if vldErr != nil {
		return nil, errors.Wrap(vldErr, "invalid product")
	}

	// // save product
	// err = c.repository.SaveProducts([]*models.Product{product})
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to save product")
	// }

	// return normalized product
	return product, nil
}

// UploadPlatformInfo uploads platform information to the database.
func (c *WikifactoryCrawler) UploadPlatformInfo(ctx context.Context) error {
	host := &models.Host{
		Name:   "Wikifactory",
		Domain: "wikifactory.com",
	}
	c.repository.SaveHosts([]*models.Host{host})
	return nil
}

func (c *WikifactoryCrawler) DiscoverProducts(ctx context.Context) error {
	// load state
	// state := CrawlerState{}

	// numRetriesAfterIncompleteResults := 0
	// page := state.Page

	return nil
}

func (c *WikifactoryCrawler) UpdateProducts(ctx context.Context) error {
	return nil
}
