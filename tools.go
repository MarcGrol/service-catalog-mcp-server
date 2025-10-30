//go:build tools

package main

import (
	_ "go.uber.org/mock/mockgen"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/goimports"
)
