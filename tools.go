// +build tools

// Package tools manages development tool versions through the module system.
//
// See https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package tools

import (
	_ "github.com/kisielk/errcheck"
	_ "github.com/mitchellh/gox"
	_ "golang.org/x/lint/golint"
)
