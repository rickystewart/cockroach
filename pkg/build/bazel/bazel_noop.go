// Copyright 2015 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package bazel

import (
	"fmt"
	"os"
)

// Shim no-op implementations of all the functions in bazel.go.in.
// Bazel builds will compile that file into the library instead.

func BuiltWithBazel() bool {
	return false
}

func FindBinary(pkg, name string) (string, bool) {
	return "", false
}

func Runfile(path string) (string, error) {
	return "", fmt.Errorf("Library not built with Bazel")
}

func RunfilesPath() (string, error) {
	return "", fmt.Errorf("Library not built with Bazel")
}

func TestTmpDir() string {
	return os.TempDir()
}
