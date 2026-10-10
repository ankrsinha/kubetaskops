//go:build tools
// +build tools

// Package hack tracks build and code generation tool dependencies.
package hack

import (
	_ "sigs.k8s.io/controller-tools/pkg/version"
)
