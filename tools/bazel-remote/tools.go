// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

// import-tools adds a blank import to tools we use such that `go mod tidy`
// doesn't clean up needed dependencies when running `go install`.

//go:build tools
// +build tools

package main

import (
	_ "github.com/buchgr/bazel-remote"
)

func main() {}
