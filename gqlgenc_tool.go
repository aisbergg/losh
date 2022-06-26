//go:build tools
// +build tools

// This package contains the dependencies for the gqlgenc tool, which is used to
// generate GraphQL clients. The tool loads models dynamically from go packages
// such as the 'github.com/99designs/gqlgen' one. Therefore this package needs
// to be inside this module (go.mod).
package tools

import (
	_ "github.com/99designs/gqlgen"
)
