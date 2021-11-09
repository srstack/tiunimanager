// Copyright 2021 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"path"

	"github.com/pingcap-inc/tiem/tiup/manager"
	operator "github.com/pingcap-inc/tiem/tiup/operation"
	"github.com/pingcap-inc/tiem/tiup/spec"
	"github.com/pingcap-inc/tiem/tiup/task"
	"github.com/pingcap/tiup/pkg/tui"
	"github.com/pingcap/tiup/pkg/utils"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

func newDeployCmd() *cobra.Command {
	opt := manager.DeployOptions{
		IdentityFile: path.Join(utils.UserHome(), ".ssh", "id_rsa"),
	}
	cmd := &cobra.Command{
		Use:          "deploy <cluster-name> <version> <topology.yaml>",
		Short:        "Deploy a TiEM cluster",
		Long:         "Deploy a TiEM cluster. SSH connection will be used to deploy files, as well as creating system users for running the service.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			shouldContinue, err := tui.CheckCommandArgsAndMayPrintHelp(cmd, args, 3)
			if err != nil {
				return err
			}
			if !shouldContinue {
				return nil
			}

			clusterName := args[0]
			version, err := utils.FmtVer(args[1])
			if err != nil {
				return err
			}
			topoFile := args[2]

			if err := supportVersion(version); err != nil {
				return err
			}

			return cm.Deploy(clusterName, version, topoFile, opt, postDeployHook, skipConfirm, gOpt)
		},
	}

	cmd.Flags().StringVarP(&opt.User, "user", "u", utils.CurrentUser(), "The user name to login via SSH. The user must has root (or sudo) privilege.")
	cmd.Flags().StringVarP(&opt.IdentityFile, "identity_file", "i", opt.IdentityFile, "The path of the SSH identity file. If specified, public key authentication will be used.")
	cmd.Flags().BoolVarP(&opt.UsePassword, "password", "p", false, "Use password of target hosts. If specified, password authentication will be used.")

	return cmd
}

func supportVersion(vs string) error {
	if !semver.IsValid(vs) {
		return nil
	}

	return nil
}

func postDeployHook(builder *task.Builder, topo spec.Topology) {
	enableTask := task.NewBuilder().Func("Setting service auto start on boot", func(ctx context.Context) error {
		return operator.Enable(ctx, topo, operator.Options{}, true)
	}).BuildAsStep("Enable service").SetHidden(true)

	builder.Parallel(false, enableTask)
}