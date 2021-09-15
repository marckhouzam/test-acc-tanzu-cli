/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu-private/tanzu-cli-apps-plugins/pkg/cli-runtime"
)

type CreateOptions struct {
	Namespace   string
	DisplayName string
	Description string
	IconUrl     string
	GitRepoUrl  string
	GitBranch   string
	Tags        []string
}

func (co *CreateOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cli.NamespaceFlag(ctx, cmd, c, &co.Namespace)
	cmd.Flags().StringVar(&co.Description, "description", "", "Accelerator description")
	cmd.Flags().StringVar(&co.DisplayName, "display-name", "", "Accelerator display name")
	cmd.Flags().StringVar(&co.IconUrl, "icon-url", "", "Accelerator icon location")
	cmd.Flags().StringSliceVar(&co.Tags, "tags", []string{}, "Accelerator Tags")
	cmd.Flags().StringVar(&co.GitRepoUrl, "git-repository", "", "Accelerator repo URL")
	cmd.Flags().StringVar(&co.GitBranch, "git-branch", "main", "Accelerator repo branch")

	cmd.MarkFlagRequired("git-repository")
}

type UpdateOptions struct {
	Namespace   string
	DisplayName string
	Description string
	IconUrl     string
	GitRepoUrl  string
	GitBranch   string
	Tags        []string
	Reconcile   bool
}

func (uo *UpdateOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cli.NamespaceFlag(ctx, cmd, c, &uo.Namespace)
	cmd.Flags().StringVar(&uo.Description, "description", "", "Accelerator description")
	cmd.Flags().StringVar(&uo.DisplayName, "display-name", "", "Accelerator display name")
	cmd.Flags().StringVar(&uo.IconUrl, "icon-url", "", "Accelerator icon location")
	cmd.Flags().StringSliceVar(&uo.Tags, "tags", []string{}, "Accelerator Tags")
	cmd.Flags().StringVar(&uo.GitRepoUrl, "git-repository", "", "Accelerator repo URL")
	cmd.Flags().StringVar(&uo.GitBranch, "git-branch", "main", "Accelerator repo branch")
	cmd.Flags().BoolVar(&uo.Reconcile, "reconcile", false, "Trigger a reconciliation including the associated GitRepository resource")
}

type DeleteOptions struct {
	Namespace string
}

func (do *DeleteOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cli.NamespaceFlag(ctx, cmd, c, &do.Namespace)
}

type ListOptions struct {
	Namespace   string
	ServerUrl   string
	FromContext bool
}

func (lo *ListOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cli.NamespaceFlag(ctx, cmd, c, &lo.Namespace)
	cmd.Flags().StringVar(&lo.ServerUrl, "server-url", "", "Accelerator UI server URL to use for retriving accelerators")
	cmd.Flags().BoolVar(&lo.FromContext, "from-context", false, "Retrieve Accelerators from current context defined in kubeconfig")
}

type GetOptions struct {
	Namespace   string
	ServerUrl   string
	FromContext bool
}

func (gopts *GetOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cli.NamespaceFlag(ctx, cmd, c, &gopts.Namespace)
	cmd.Flags().StringVar(&gopts.ServerUrl, "server-url", "", "Accelerator UI server URL to use for retriving accelerators")
	cmd.Flags().BoolVar(&gopts.FromContext, "from-context", false, "Retrieve Accelerator from current context defined in kubeconfig")
}
