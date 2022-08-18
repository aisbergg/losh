package assets

import (
	"embed"
	"net/http"

	"github.com/spf13/afero"
)

//go:embed icons
//go:embed static
//go:embed templates
var embeddedAssets embed.FS

// AssetsHTTP holds the assets of the application for the http.Filesystem.
var AssetsHTTP = http.FS(embeddedAssets)

// AssetsAfero holds the assets of the application in form of afero.Fs.
var AssetsAfero = afero.FromIOFS{FS: embeddedAssets}
