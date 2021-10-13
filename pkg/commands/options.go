/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/tanzu-cli-apps-plugins/pkg/cli-runtime"
)

type CreateOptions struct {
	Namespace   string
	DisplayName string
	Description string
	IconUrl     string
	GitRepoUrl  string
	GitBranch   string
	GitTag      string
	GitInterval string
	LocalPath   string
	SourceImage string
	Tags        []string
}

func (co *CreateOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cli.NamespaceFlag(ctx, cmd, c, &co.Namespace)
	cmd.Flags().StringVar(&co.Description, "description", "", "description of this accelerator")
	cmd.Flags().StringVar(&co.DisplayName, "display-name", "", "display name for the accelerator")
	cmd.Flags().StringVar(&co.IconUrl, "icon-url", "", "URL for icon to use with the accelerator")
	cmd.Flags().StringSliceVar(&co.Tags, "tags", []string{}, "tags that can be used to search for accelerators")
	cmd.Flags().StringVar(&co.GitRepoUrl, "git-repository", "", "Git repository URL for the accelerator")
	cmd.Flags().StringVar(&co.GitBranch, "git-branch", "", "Git repository branch to be used")
	cmd.Flags().StringVar(&co.GitTag, "git-tag", "", "Git repository tag to be used")
	cmd.Flags().StringVar(&co.GitInterval, "git-interval", "", "interval at which to check for Git repository updates")
	cmd.Flags().StringVar(&co.SourceImage, "source-image", "", "name of the source image for the accelerator")
	cmd.Flags().StringVar(&co.LocalPath, "local-path", "", "path to the directory containing the source for the accelerator")
}

type UpdateOptions struct {
	Namespace   string
	DisplayName string
	Description string
	IconUrl     string
	GitRepoUrl  string
	GitBranch   string
	GitTag      string
	GitInterval string
	Tags        []string
	SourceImage string
	Reconcile   bool
}

func (uo *UpdateOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cli.NamespaceFlag(ctx, cmd, c, &uo.Namespace)
	cmd.Flags().StringVar(&uo.Description, "description", "", "description of this accelerator")
	cmd.Flags().StringVar(&uo.DisplayName, "display-name", "", "display name for the accelerator")
	cmd.Flags().StringVar(&uo.IconUrl, "icon-url", "", "URL for icon to use with the accelerator")
	cmd.Flags().StringSliceVar(&uo.Tags, "tags", []string{}, "tags that can be used to search for accelerators")
	cmd.Flags().StringVar(&uo.GitRepoUrl, "git-repository", "", "Git repository URL for the accelerator")
	cmd.Flags().StringVar(&uo.GitBranch, "git-branch", "main", "Git repository branch to be used")
	cmd.Flags().StringVar(&uo.GitTag, "git-tag", "", "Git repository tag to be used")
	cmd.Flags().BoolVar(&uo.Reconcile, "reconcile", false, "trigger a reconciliation including the associated GitRepository resource")
	cmd.Flags().StringVar(&uo.GitInterval, "git-interval", "", "interval at which to check for Git repository updates")
	cmd.Flags().StringVar(&uo.SourceImage, "source-image", "", "name of the source image for the accelerator")
}

type PushOptions struct {
	LocalPath   string
	SourceImage string
}

func (po *PushOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cmd.Flags().StringVar(&po.SourceImage, "source-image", "", "name of the source image for the accelerator")
	cmd.MarkFlagRequired("source-image")
	cmd.Flags().StringVar(&po.LocalPath, "local-path", "", "path to the directory containing the source for the accelerator")
	cmd.MarkFlagRequired("local-path")
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
	defaultUiServerUrl := EnvVar("ACC_SERVER_URL", "http://localhost:8877")
	cli.NamespaceFlag(ctx, cmd, c, &lo.Namespace)
	cmd.Flags().StringVar(&lo.ServerUrl, "server-url", defaultUiServerUrl, "the URL for the Application Accelerator server")
	cmd.Flags().BoolVar(&lo.FromContext, "from-context", false, "retrieve resources from current context defined in kubeconfig")
}

type GetOptions struct {
	Namespace   string
	ServerUrl   string
	FromContext bool
}

func (gopts *GetOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	defaultUiServerUrl := EnvVar("ACC_SERVER_URL", "http://localhost:8877")
	cli.NamespaceFlag(ctx, cmd, c, &gopts.Namespace)
	cmd.Flags().StringVar(&gopts.ServerUrl, "server-url", defaultUiServerUrl, "the URL for the Application Accelerator server")
	cmd.Flags().BoolVar(&gopts.FromContext, "from-context", false, "retrieve resources from current context defined in kubeconfig")
}
