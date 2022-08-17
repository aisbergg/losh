package binding

import (
	"losh/web/build/assets"
	"losh/web/core/config"

	"github.com/aisbergg/go-pathlib/pkg/pathlib"
	orderedmap "github.com/wk8/go-ordered-map"
)

type ColorsItem struct {
	Class string `json:"class" liquid:"class"`
	Hex   string `json:"hex" liquid:"hex"`
	Title string `json:"title" liquid:"title"`
}

// Site contains the site configuration, which is used in the tabler.io
// HTML templates.
type Site struct {
	Debug            bool                   `json:"debug" liquid:"debug"`
	UseIconfont      bool                   `json:"useIconfont" liquid:"use-iconfont"`
	Title            string                 `json:"title" liquid:"title"`
	CopyRight        string                 `json:"copyRight" liquid:"copy-right"`
	Description      string                 `json:"description" liquid:"description"`
	IssueURL         string                 `json:"issueUrl" liquid:"issue-url"`
	LayoutDark       bool                   `json:"layoutDark" liquid:"layout-dark"`
	TablerCSSPlugins []interface{}          `json:"tablerCssPlugins" liquid:"tabler-css-plugins"`
	Data             map[string]interface{} `json:"data" liquid:"data"`

	// for charts and such

	MonthsShort []string              `json:"monthsShort" liquid:"months-short"`
	MonthsLong  []string              `json:"monthsLong" liquid:"months-long"`
	Colors      map[string]ColorsItem `json:"colors" liquid:"colors"`
}

// newSiteBinding creates a new Site.
func newSiteBinding(config *config.Config) *Site {
	icnDir := pathlib.NewPosixPathWithFS(assets.AssetsAfero, "icons")
	icnPth, err := icnDir.ReadDir()
	if err != nil {
		panic(err)
	}
	icons := make(map[string]string, len(icnPth))
	for _, p := range icnPth {
		icnCnt, err := p.ReadFile()
		if err != nil {
			panic(err)
		}
		icons[p.Stem()] = string(icnCnt)
	}

	return &Site{
		Debug:            config.Debug.Enabled,
		UseIconfont:      false, // TODO: use built-in
		Title:            "LOSH",
		CopyRight:        "André Lehmann",
		Description:      "Library of Open Source Hardware",
		IssueURL:         "https://github.com/aisbergg/losh/issues",
		LayoutDark:       false,
		TablerCSSPlugins: []interface{}{},
		Data: map[string]interface{}{
			"menu": orderedmap.NewWithPairs(
				"search", map[string]interface{}{
					"url":   "search",
					"icon":  "search",
					"title": "Search",
				},
				"about", map[string]interface{}{
					"icon":  "info-circle",
					"title": "About",
					"children": orderedmap.NewWithPairs(
						"project", map[string]interface{}{
							"url":   "about/project",
							"title": "The Project",
							"icon":  "dna-2",
						},
						"faq", map[string]interface{}{
							"url":   "about/faq",
							"title": "FAQ",
							"icon":  "question-mark",
						},
					),
				}),
			"icons": icons,
		},
		MonthsShort: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
		MonthsLong:  []string{"January", "Febuary", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
		Colors: map[string]ColorsItem{
			"blue": ColorsItem{
				Class: "blue",
				Hex:   "#206bc4",
				Title: "Blue",
			},
			"azure": ColorsItem{
				Class: "azure",
				Hex:   "#45aaf2",
				Title: "Azure",
			},
			"indigo": ColorsItem{
				Class: "indigo",
				Hex:   "#6574cd",
				Title: "Indigo",
			},
			"purple": ColorsItem{
				Class: "purple",
				Hex:   "#a55eea",
				Title: "Purple",
			},
			"pink": ColorsItem{
				Class: "pink",
				Hex:   "#f66d9b",
				Title: "Pink",
			},
			"red": ColorsItem{
				Class: "red",
				Hex:   "#fa4654",
				Title: "Red",
			},
			"orange": ColorsItem{
				Class: "orange",
				Hex:   "#fd9644",
				Title: "Orange",
			},
			"yellow": ColorsItem{
				Class: "yellow",
				Hex:   "#f1c40f",
				Title: "Yellow",
			},
			"lime": ColorsItem{
				Class: "lime",
				Hex:   "#7bd235",
				Title: "Lime",
			},
			"green": ColorsItem{
				Class: "green",
				Hex:   "#5eba00",
				Title: "Green",
			},
			"teal": ColorsItem{
				Class: "teal",
				Hex:   "#2bcbba",
				Title: "Teal",
			},
			"cyan": ColorsItem{
				Class: "cyan",
				Hex:   "#17a2b8",
				Title: "Cyan",
			},
		},
	}
}

type TemplateBindingProvider struct {
	site *Site
}

// NewTemplateBindingProvider creates a new TemplateContextProvider.
func NewTemplateBindingProvider(config *config.Config) TemplateBindingProvider {
	return TemplateBindingProvider{
		site: newSiteBinding(config),
	}
}

// Get returns a new context for templating use.
func (p TemplateBindingProvider) Get() map[string]interface{} {
	bnd := make(map[string]interface{}, 3)
	bnd["site"] = p.site
	bnd["page"] = make(map[string]interface{}, 10)
	return bnd
}
