// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package main

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
)

const homeDirFlag = "home-dir"

// MakeBuilderCmd constructs the subcommand used to run
func makeBuilderCmd(runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	builderCmd := &cobra.Command{
		Use:     "builder",
		Short:   "Run the Bazel builder image.",
		Long:    "Run the Bazel builder image.",
		Example: `dev builder`,
		Args:    cobra.ExactArgs(0),
		RunE:    runE,
	}
	builderCmd.Flags().String(homeDirFlag, "",
		"directory to mount as the home directory in the container")
	return builderCmd
}

func (d *dev) builder(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	homeDir := mustGetFlagString(cmd, homeDirFlag)
	args, err := d.getDockerRunArgs(ctx, homeDir, true)
	if err != nil {
		return err
	}
	return d.exec.CommandContextInheritingStdStreams(ctx, "docker", args...)
}

func (d *dev) getDockerRunArgs(
	ctx context.Context, homeDir string, tty bool,
) (args []string, err error) {
	err = d.ensureBinaryInPath("docker")
	if err != nil {
		return
	}
	if homeDir == "" {
		var cacheDir string
		cacheDir, err = os.UserCacheDir()
		if err != nil {
			return
		}
		homeDir = filepath.Join(cacheDir, "cockroachbuild", "bzlhome")
	}
	err = d.os.MkdirAll(homeDir)
	if err != nil {
		return
	}

	args = append(args, "run", "--rm")
	if tty {
		args = append(args, "-it")
	} else {
		args = append(args, "-i")
	}
	workspace, err := d.getWorkspace(ctx)
	if err != nil {
		return
	}
	args = append(args, "-v", workspace+":/cockroach:ro")
	args = append(args, "--workdir=/cockroach")
	// Create the artifacts directory.
	artifacts := filepath.Join(workspace, "artifacts")
	err = d.os.MkdirAll(artifacts)
	if err != nil {
		return
	}
	args = append(args, "-v", artifacts+":/artifacts")
	// The `delegated` switch ensures that the container's view of the cache
	// is authoritative. This can result in writes to the actual underlying
	// filesystem to be lost, but it's a cache so we don't care about that.
	args = append(args, "-v", homeDir+":/home/roach:delegated")

	// Apply the same munging for the UID/GID that we do in build/builder.sh.
	// Quoth a comment from there:
	// Attempt to run in the container with the same UID/GID as we have on the host,
	// as this results in the correct permissions on files created in the shared
	// volumes. This isn't always possible, however, as IDs less than 100 are
	// reserved by Debian, and IDs in the low 100s are dynamically assigned to
	// various system users and groups. To be safe, if we see a UID/GID less than
	// 500, promote it to 501. This is notably necessary on macOS Lion and later,
	// where administrator accounts are created with a GID of 20. This solution is
	// not foolproof, but it works well in practice.
	currentUser, err := user.Current()
	if err != nil {
		return
	}
	uid := currentUser.Uid
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return
	}
	if uidInt < 500 {
		uid = "501"
	}
	gid := currentUser.Gid
	gidInt, err := strconv.Atoi(gid)
	if err != nil {
		return
	}
	if gidInt < 500 {
		gid = uid
	}
	args = append(args, "-u", fmt.Sprintf("%s:%s", uid, gid))
	
	// Read the Docker image from build/teamcity-bazel-support.sh.
	buf, err := d.os.ReadFile(filepath.Join(workspace, "build/teamcity-bazel-support.sh"))
	if err != nil {
		return
	}
	var bazelImage string
	for _, line := range strings.Split(buf, "\n") {
		if strings.HasPrefix(line, "BAZEL_IMAGE=") {
			bazelImage = strings.Trim(strings.TrimPrefix(line, "BAZEL_IMAGE="), "\n ")
		}
	}
	if bazelImage == "" {
		err = errors.New("Could not find BAZEL_IMAGE in build/teamcity-bazel-support.sh")
		return
	}
	args = append(args, bazelImage)
	return
}
