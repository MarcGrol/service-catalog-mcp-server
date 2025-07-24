//go:build tools
// +build tools

package main

import (
	_ "github.com/golang/mock/mockgen"
	_ "golang.org/x/lint/golint"
)
