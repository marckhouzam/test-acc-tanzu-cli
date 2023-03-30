/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"

	"github.com/spf13/cobra"
	cli "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
)

func PushCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := PushOptions{}
	cmd := &cobra.Command{
		Use:     "push",
		Short:   "(DEPRECTAED) Push local path to source image",
		Long:    "(DEPRECATED) Push source code from local path to source image used by an accelerator",
		Example: "tanzu accelerator push --local-path <local path> --source-image <image>",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			if opts.LocalPath != "" {
				if err := opts.PublishLocalSource(ctx, c); err != nil {
					return err
				}
			}
			return nil

		},
	}
	opts.DefineFlags(ctx, cmd, c)
	return cmd
}
