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
	"path"
	"path/filepath"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/build/bazel"
	"github.com/cockroachdb/cockroach/pkg/util/envutil"
	"github.com/stretchr/testify/require"
)

// TestDataPath returns a path to an asset in the testdata directory. It knows
// to access accesses the right path when executing under bazel.
//
// For example, if there is a file testdata/a.txt, you can get a path to that
// file using TestDataPath(t, "a.txt").
//
// This function has to handle a few distinct cases:
// 1. "Relative testdata". In this case the test is assumed to be running in the
//    package directory or something like the package directory, in which case
//    testdata can be found in the current working directory at ./testdata. This
//    is the case for non-Bazel tests as well as Bazel tests that have
//    COCKROACH_RELATIVE_TESTDATA set. The latter of these is true for
//    roachprod-stress tests where we use Bazel to build the binary but the test
//    won't have access to the runfiles.
// 2. "COCKROACH_WORKSPACE"-relative testdata. In this case we pass the full
//    path to the root of the cockroach workspace as an environment variable and
//    infer the location of the *in-workspace* testdata accordingly. `dev` uses
//    this for `--rewrite` support -- we won't have permission to rewrite inside
//    the Bazel runfiles tree so we need the path to the actual workspace.
// 3. In any other case we just use the bazel library to find the path to the
// runfiles tree and get the path accordingly.
func TestDataPath(t testing.TB, relative ...string) string {
	relative = append([]string{"testdata"}, relative...)
	useRelativeTestData := envutil.EnvOrDefaultBool("COCKROACH_RELATIVE_TESTDATA", false)
	if bazel.BuiltWithBazel() && !useRelativeTestData {
		//lint:ignore SA4006 apparently a linter bug.
		cockroachWorkspace, set := envutil.EnvString("COCKROACH_WORKSPACE", 0)
		// dev notifies the library that the test is running in a subdirectory of the
		// workspace with the environment COCKROACH_WORKSPACE.
		if set {
			return path.Join(cockroachWorkspace, bazel.RelativeTestTargetPath(), path.Join(relative...))
		}
		runfiles, err := bazel.RunfilesPath()
		require.NoError(t, err)
		return path.Join(runfiles, bazel.RelativeTestTargetPath(), path.Join(relative...))
	}

	// Otherwise we're in the package directory and can just return a relative path.
	ret := path.Join(relative...)
	ret, err := filepath.Abs(ret)
	require.NoError(t, err)
	return ret
}
