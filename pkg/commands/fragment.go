/*
Copyright 2022-2023 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"

	"github.com/spf13/cobra"
	cli "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
)

func FragmentCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fragment",
		Short:   "Fragment commands",
		Long:    "Commands to manage accelerator fragments",
		Example: "tanzu accelerator fragment --help",
		Aliases: []string{"frag"},
	}
	cmd.AddCommand(FragmentListCmd(ctx, c))
	cmd.AddCommand(FragmentCreateCmd(ctx, c))
	cmd.AddCommand(FragmentGetCmd(ctx, c))
	cmd.AddCommand(FragmentUpdateCmd(ctx, c))
	cmd.AddCommand(FragmentDeleteCmd(ctx, c))

	return cmd
}
