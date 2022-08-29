package wikifactory

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"losh/crawler/core/validator"
	"losh/crawler/core/wikifactory/wfclient"
	"losh/internal/core/product/models"
	"losh/internal/core/product/services"
	lerrors "losh/internal/lib/errors"
	"losh/internal/lib/log"
	"losh/internal/lib/net/download"
	"losh/internal/lib/net/ratelimit"
	"losh/internal/lib/net/request"

	gql "github.com/Yamashou/gqlgenc/clientv2"
	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"

	"go.uber.org/zap"

	n "losh/internal/lib/net"
	"losh/internal/lib/unit"
)

// Constants that define the crawler behavior.
const (
	crawlerName         = "wikifactory.com"
	timeout             = time.Duration(30) * time.Second
	batchSize           = 10
	retries             = 5
	maxFileSizeManifest = 10 * unit.MiB
	maxWaitTime         = 5 * time.Minute
	maxRedirects        = 5
)

type CrawlerState struct {
	StartTime   time.Time     `json:"startTime"`
	ElapsedTime time.Duration `json:"elapsedTime"`
	NumCrawled  int64         `json:"numCrawled"`
	NumIndexed  int64         `json:"numIndexed"`
	Cursor      string        `json:"cursor"`
	Page        int64         `json:"page"`
}

// WikifactoryCrawler is a product crawler for Wikifactory.
type WikifactoryCrawler struct {
	productService *services.Service
	validator      *validator.Validator

	fileDownloader *download.Downloader
	wfClient       wfclient.WikifactoryGraphQLClient
	log            *zap.SugaredLogger
}

// NewWikifactoryCrawler creates a new WikifactoryCrawler.
func NewWikifactoryCrawler(productService *services.Service, userAgent string) *WikifactoryCrawler {
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
		AddRateLimiter(ratelimit.NewTimedeltaRateLimiter(5*time.Second, 20))
	wfClient := wfclient.NewClient(graphQLRequester)

	validator := validator.NewValidator(productService)

	return &WikifactoryCrawler{
		productService: productService,
		validator:      validator,

		fileDownloader: fileDownloader,
		wfClient:       wfClient,

		log: log,
	}
}

func (c *WikifactoryCrawler) DiscoverProducts(ctx context.Context) error {
	c.log.Infof("discovering products on %s", crawlerName)

	// TODO: make state file configurable
	stateFilePath := pathlib.NewPath("crawler-state.json")
	state, err := c.loadState(stateFilePath)
	if err != nil {
		return lerrors.NewAppErrorWrap(err, "failed to load state")
	}

	runStartedAt := time.Now()
	prevElapsedTime := state.ElapsedTime
	page := state.Page
	cursor := state.Cursor
	hasNextPage := true
	for hasNextPage {
		c.log.Debugf("getting %d results from page %d (cursor: %s)", batchSize, page, cursor)

		// get project information
		queryProjects, err := c.wfClient.QueryProjects(ctx, batchSize, cursor)
		if err != nil {
			return lerrors.NewAppErrorWrap(err, "failed to get project information")
		}
		discoveredAt := time.Now()
		wfPrjsInfo := queryProjects.Projects.Result.Edges

		// check each project and index if compliant
		for _, edge := range wfPrjsInfo {
			wfPrjInfo := edge.Node
			productID := models.NewProductID(crawlerName, *wfPrjInfo.ParentSlug, *wfPrjInfo.Slug, "")

			// check mandatory fields for compliance
			err = c.checkMandatory(wfPrjInfo)
			if err != nil {
				c.log.Debugf("skipping (%s): %s", productID.String(), err.Error())
				continue
			}

			// get full product information
			c.log.Infof("indexing product (%s)", productID.String())
			prd, err := c.getProduct(ctx, productID, discoveredAt)
			if err != nil {
				if vldErr, ok := err.(*validator.ValidationError); ok {
					c.log.Debugf("skipping (%s): %s", productID.String(), vldErr.Error())
					continue
				}
				return lerrors.NewAppErrorWrap(err, "failed to get product").Add("crawlerProductID", productID.String())
			}

			// save product
			c.log.Debugf("saving product (%s)", productID.String())
			err = c.productService.SaveNode(ctx, prd)
			if err != nil {
				return lerrors.NewAppErrorWrap(err, "failed to save product").Add("crawlerProductID", productID.String())
			}

			state.NumIndexed++
		}

		pageInfo := queryProjects.Projects.Result.PageInfo
		hasNextPage = pageInfo.HasNextPage
		page++
		cursor = *pageInfo.EndCursor

		// save state
		state.Cursor = cursor
		state.Page = page
		state.NumCrawled += int64(len(wfPrjsInfo))
		state.ElapsedTime = prevElapsedTime + time.Now().Sub(runStartedAt)
		err = c.saveState(stateFilePath, state)
		if err != nil {
			return lerrors.NewAppErrorWrap(err, "failed to save state")
		}
		c.log.Debugf("indexed %d of %d products on %s this far", state.NumIndexed, state.NumCrawled, crawlerName)
	}

	// TODO
	// remove state file

	return nil
}

func (c *WikifactoryCrawler) UpdateProducts(ctx context.Context) error {
	return nil
}

func (c *WikifactoryCrawler) getProduct(ctx context.Context, prdID models.ProductID, discoveredAt time.Time) (*models.Product, error) {
	// get full project information
	getProjectFullBySlug, err := c.wfClient.GetProjectFullBySlug(ctx, prdID.Owner, prdID.Repo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Wikifactory project information")
	}
	if getProjectFullBySlug == nil {
		return nil, errors.New("Wikifactory project not found")
	}

	product, err := c.NormalizeProduct(ctx, discoveredAt, getProjectFullBySlug.Project.Result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to normalize product information")
	}

	// validate product
	err = c.validator.ValidateProduct(product)
	if err != nil {
		return nil, err
	}

	// return normalized product
	return product, nil
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
	getProjectMandatoryBySlug, err := c.wfClient.GetProjectMandatoryBySlug(ctx, productID.Owner, productID.Repo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project information")
	}
	projectInfo := getProjectMandatoryBySlug.Project.Result
	if projectInfo == nil {
		return nil, errors.New("project not found")
	}
	discoveredAt := time.Now()

	// check mandatory fields
	err = c.checkMandatory(projectInfo)
	if err != nil {
		return nil, err
	}

	// get full project information
	getProjectFullBySlug, err := c.wfClient.GetProjectFullBySlug(ctx, productID.Owner, productID.Repo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project information")
	}
	if getProjectFullBySlug == nil {
		return nil, errors.New("project not found")
	}

	product, err := c.NormalizeProduct(ctx, discoveredAt, getProjectFullBySlug.Project.Result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to normalize project information")
	}

	// validate product
	vldErr := c.validator.ValidateProduct(product)
	if vldErr != nil {
		return nil, errors.Wrap(vldErr, "invalid product")
	}

	// return normalized product
	return product, nil
}

// loadState loads the crawler state from the state file.
func (c *WikifactoryCrawler) loadState(path pathlib.Path) (CrawlerState, error) {
	// load state
	exists, err := path.Exists()
	if err != nil {
		return CrawlerState{}, errors.Wrap(err, "failed to check file existence")
	}
	// new state
	state := CrawlerState{
		StartTime: time.Now(),
		Page:      1,
		Cursor:    "",
	}
	if !exists {
		return state, nil
	}
	fileContent, err := path.ReadFile()
	if err != nil {
		return CrawlerState{}, errors.Wrap(err, "failed to read state file")
	}
	err = json.Unmarshal(fileContent, &state)
	if err != nil {
		return CrawlerState{}, errors.Wrap(err, "failed to unmarshal state file")
	}
	return state, nil
}

// saveState saves the crawler state to the state file.
func (c *WikifactoryCrawler) saveState(path pathlib.Path, state CrawlerState) error {
	fileContent, err := json.Marshal(state)
	if err != nil {
		return errors.Wrap(err, "failed to marshal state")
	}
	err = path.WriteFile(fileContent)
	if err != nil {
		return errors.Wrap(err, "failed to write state file")
	}
	return nil
}

// checkMandatory checks if mandatory fields are present.
func (c *WikifactoryCrawler) checkMandatory(projectInfo *wfclient.ProjectMandatoryFragment) error {
	var (
		version   string
		license   string
		filePaths []string
	)
	if projectInfo.License != nil {
		license = *projectInfo.License.Abreviation
	}
	if projectInfo.Contribution != nil {
		version = *projectInfo.Contribution.Version
		filePaths = make([]string, 0, len(projectInfo.Contribution.Files))
		for _, file := range projectInfo.Contribution.Files {
			if file == nil || file.File == nil {
				continue
			}
			if file.Dirname != nil && *file.Dirname != "" {
				filePaths = append(filePaths, *file.Dirname+"/"+file.File.Filename)
			} else {
				filePaths = append(filePaths, file.File.Filename)
			}
		}
	}

	// quick check if mandatory fields are present
	mdtFlds := validator.MandatoryFields{
		OKHV:                  "OKH-LOSHv1.0",
		Name:                  stringOrEmpty(projectInfo.Name),
		Description:           stringOrEmpty(projectInfo.Description),
		Version:               version,
		Repository:            "https://wikifactory.com", // just to pass the check
		License:               c.translateLicense(license),
		Licensor:              stringOrEmpty(projectInfo.Creator.Profile.Username),
		DocumentationLanguage: *c.normDocumentationLanguage(stringOrEmpty(projectInfo.Description)),
		FilePaths:             filePaths,
	}
	return c.validator.ValidateMandatory(mdtFlds)
}
