package lib

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type X struct {
	Title string `json:"title"`
}

// site defaults
// site dynamic by config
// page defaults
// page dynamic by config
// page dynamic by handler

// XXX: just for testing
func GetBindings() map[string]interface{} {
	siteCofigFile, err := ioutil.ReadFile("resources/views/_config.yml")
	if err != nil {
		log.Fatal(err)
	}
	siteBindings := yaml.MapSlice{} // preserve order
	// siteBindings := make(map[interface{}]interface{})
	err = yaml.Unmarshal(siteCofigFile, &siteBindings)
	if err != nil {
		log.Fatal(err)
	}
	bindings := map[string]interface{}{
		"page": map[string]interface{}{
			"title":                   "Blank page",
			"menu":                    "second",
			"layout-navbar-condensed": true,
			"container-centered":      true,
		},
		"site": siteBindings,
	}
	return bindings
}
