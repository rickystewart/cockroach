// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package skip

// SkippableTest is a testing.TB with Skip methods.
type SkippableTest interface {
	Helper()
	Skip(...interface{})
	Skipf(string, ...interface{})
}

// IgnoreLint skips this test, explicitly marking it as not a test that
// should be tracked as a "skipped test" by external tools. You should use this
// if, for example, your test should only be run in Race mode.
func IgnoreLint(t SkippableTest, args ...interface{}) {
	t.Helper()
	t.Skip(args...)
}
