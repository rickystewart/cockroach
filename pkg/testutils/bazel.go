// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package testutils

import (
	"os"
	"path"
	"strings"

	"github.com/cockroachdb/cockroach/pkg/build/bazel"
)

//
// Utilities in this file are intended to provide helpers for tests
// to access system resources in a correct way, when executing under bazel.
//
// See  https://docs.bazel.build/versions/master/test-encyclopedia.html for more details
// on the directory layout.
//
// When test is compiled with bazel, bazel will create a directory, called RUNFILES directory,
// containing all of the resources required to execute such test (libraries, final binary,
// data resources, etc).
//
// Bazel also sets up various environmental variables that point to the location(s) of
// those resources.

// Name of the environment variable containing the bazel target path (//pkg/cmd/foo:bar).
const testTargetEnv = "TEST_TARGET"

// bazeRelativeTargetPath returns relative path to the package
// of the current test.
func bazelRelativeTargetPath() string {
	target := os.Getenv(testTargetEnv)
	if target == "" {
		return ""
	}

	// Drop target name.
	if last := strings.LastIndex(target, ":"); last > 0 {
		target = target[:last]
	}
	return strings.TrimPrefix(target, "//")
}

// TestDataPath returns a path to an asset in the testdata directory.
//
// For example, if there is a file testdata/a.txt, you can get a path to that
// file with TestDataPath("a.txt").
func TestDataPath(relative ...string) (string, error) {
	relative = append([]string{"testdata"}, relative...)
	if bazel.BuiltWithBazel() {
		runfiles, err := bazel.RunfilesPath()
		if err != nil {
			return "", err
		}
		return path.Join(runfiles, bazelRelativeTargetPath(), path.Join(relative...)), nil
	}
	// If we're not running in Bazel, we're in the package directory and can
	// just return a relative path.
	return path.Join(relative...), nil
}
