package lib

import (
	"log"

	"gopkg.in/yaml.v2"
)

const siteConfig string = `
debug: true # TODO: make configurable

use-iconfont: true  # TODO: use built-in

title: LOSH
# base: 127.0.0.1:8080  # used in footer
copy-right: Library of Open Source Hardware # used in footer
description: Library of Open Source Hardware
issue-url: XXX

layout-dark: false

tabler-css-plugins: []


colors:
  blue:
    class: blue
    hex: "#206bc4"
    title: Blue
  azure:
    class: azure
    hex: "#45aaf2"
    title: Azure
  indigo:
    class: indigo
    hex: "#6574cd"
    title: Indigo
  purple:
    class: purple
    hex: "#a55eea"
    title: Purple
  pink:
    class: pink
    hex: "#f66d9b"
    title: Pink
  red:
    class: red
    hex: "#fa4654"
    title: Red
  orange:
    class: orange
    hex: "#fd9644"
    title: Orange
  yellow:
    class: yellow
    hex: "#f1c40f"
    title: Yellow
  lime:
    class: lime
    hex: "#7bd235"
    title: Lime
  green:
    class: green
    hex: "#5eba00"
    title: Green
  teal:
    class: teal
    hex: "#2bcbba"
    title: Teal
  cyan:
    class: cyan
    hex: "#17a2b8"
    title: Cyan

data:
  menu:
    search:
      url: "search"
      icon: search
      title: Search

    explore:
      url: "explore"
      icon: compass
      title: Explore

    about:
      icon: info-circle
      title: About
      children:
        project:
          url: "about/project"
          title: The Project
        faq:
          url: "about/faq"
          title: FAQ
          icon: question-mark

months-short: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]
months-long: ["January", "Febuary", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]
`

// site defaults
// site dynamic by config
// page defaults
// page dynamic by config
// page dynamic by handler

// XXX: just for testing
func GetBindings() map[string]interface{} {
	// siteCofigContent, err := ioutil.ReadFile("resources/views/_config.yml")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	siteCofigContent := siteConfig
	siteBindings := yaml.MapSlice{} // preserves the order
	err := yaml.Unmarshal([]byte(siteCofigContent), &siteBindings)
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
