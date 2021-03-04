// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package stagebazelgotoolchain

import (
	"os"
	"path"
	"path/filepath"
	
	"github.com/cockroachdb/cockroach/pkg/build/bazel"
)

// When imported, this package stages the Go binary from the Bazel toolchain
// in the `PATH` and sets an appropriate `GOROOT`. Does nothing if the library
// was not built with Bazel. Useful for tests that need access to a Go toolchain
// (lints, etc.)

func init() {
	if !bazel.BuiltWithBazel() {
		return
	}
	go_bin, err := bazel.Runfile("bin/go")
	if err != nil {
		panic(err)
	}
	path_env := os.Getenv("PATH")
	if err := os.Setenv("PATH", filepath.Dir(go_bin) + string(os.PathListSeparator) + path_env); err != nil {
		panic(err)
	}
	if err := os.Setenv("GOROOT", filepath.Dir(filepath.Dir(go_bin))); err != nil {
		panic(err)
	}
	if err := os.Setenv("GOCACHE", path.Join(bazel.TestTmpDir(), ".gocache")); err != nil {
		panic(err)
	}
}
