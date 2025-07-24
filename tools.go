//go:build tools
// +build tools

package main

import (
	_ "golang.org/x/lint/golint"
	_ "go.uber.org/mock/mockgen"
)
