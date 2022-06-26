//go:build tools
// +build tools

// This package contains the tool dependencies of the project.
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package tools

import (
	_ "github.com/Yamashou/gqlgenc"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
