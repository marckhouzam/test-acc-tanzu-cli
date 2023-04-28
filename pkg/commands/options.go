/*
Copyright 2021-2023 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
)

type CreateOptions struct {
	Namespace   string
	DisplayName string
	Description string
	IconUrl     string
	GitRepoUrl  string
	GitBranch   string
	GitTag      string
	GitSubPath  string
	Interval    string
	// Deprecated: LocalPath is deprecated
	LocalPath string
	// Deprecated: SourceImage is deprecated
	SourceImage string
	SecretRef   string
	Tags        []string
}

type FragmentCreateOptions struct {
	Namespace   string
	DisplayName string
	GitRepoUrl  string
	GitBranch   string
	GitTag      string
	GitSubPath  string
	Interval    string
	// Deprecated: LocalPath is deprecated
	LocalPath string
	// Deprecated: SourceImage is deprecated
	SourceImage string
	SecretRef   string
}

func normalizeGitRepoRun(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "git-repository":
		name = "git-repo"
	}
	return pflag.NormalizedName(name)
}

func (co *CreateOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cmd.Flags().StringVarP(&co.Namespace, "namespace", "n", "accelerator-system", "namespace for accelerator system")
	cmd.Flags().StringVar(&co.Description, "description", "", "description of this accelerator")
	cmd.Flags().StringVar(&co.DisplayName, "display-name", "", "display name for the accelerator")
	cmd.Flags().StringVar(&co.IconUrl, "icon-url", "", "URL for icon to use with the accelerator")
	cmd.Flags().StringSliceVar(&co.Tags, "tags", []string{}, "tags that can be used to search for accelerators")
	cmd.Flags().StringVar(&co.GitRepoUrl, "git-repository", "", "Git repository URL for the accelerator")
	cmd.Flags().StringVar(&co.GitBranch, "git-branch", "main", "Git repository branch to be used")
	cmd.Flags().StringVar(&co.GitTag, "git-tag", "", "Git repository tag to be used")
	cmd.Flags().StringVar(&co.GitSubPath, "git-sub-path", "", "Git repository subPath to be used")
	cmd.Flags().StringVar(&co.Interval, "interval", "", "interval for checking for updates to Git or image repository")
	cmd.Flags().StringVar(&co.SourceImage, "source-image", "", "(DEPRECATED) name of the source image for the accelerator")
	cmd.Flags().StringVar(&co.SecretRef, "secret-ref", "", "name of secret containing credentials for private Git or image repository")
	cmd.Flags().StringVar(&co.LocalPath, "local-path", "", "(DEPRECATED) path to the directory containing the source for the accelerator")
	cmd.Flags().SetNormalizeFunc(normalizeGitRepoRun)
}

func (co *FragmentCreateOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cmd.Flags().StringVarP(&co.Namespace, "namespace", "n", "accelerator-system", "namespace for accelerator system")
	cmd.Flags().StringVar(&co.DisplayName, "display-name", "", "display name for the accelerator fragment")
	cmd.Flags().StringVar(&co.GitRepoUrl, "git-repository", "", "Git repository URL for the accelerator fragment")
	cmd.Flags().StringVar(&co.GitBranch, "git-branch", "main", "Git repository branch to be used")
	cmd.Flags().StringVar(&co.GitTag, "git-tag", "", "Git repository tag to be used")
	cmd.Flags().StringVar(&co.GitSubPath, "git-sub-path", "", "Git repository subPath to be used")
	cmd.Flags().StringVar(&co.Interval, "interval", "", "interval for checking for updates to Git or image repository")
	cmd.Flags().StringVar(&co.SourceImage, "source-image", "", "(DEPRECATED) name of the source image for the accelerator")
	cmd.Flags().StringVar(&co.SecretRef, "secret-ref", "", "name of secret containing credentials for private Git or image repository")
	cmd.Flags().StringVar(&co.LocalPath, "local-path", "", "(DEPRECATED) path to the directory containing the source for the accelerator fragment")
	cmd.Flags().SetNormalizeFunc(normalizeGitRepoRun)
}

type UpdateOptions struct {
	Namespace   string
	DisplayName string
	Description string
	IconUrl     string
	GitRepoUrl  string
	GitBranch   string
	GitTag      string
	GitSubPath  string
	Interval    string
	// Deprecated: SourceImage is deprecated
	SourceImage string
	SecretRef   string
	Tags        []string
	Reconcile   bool
}

func (uo *UpdateOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cmd.Flags().StringVarP(&uo.Namespace, "namespace", "n", "accelerator-system", "namespace for accelerator system")
	cmd.Flags().StringVar(&uo.Description, "description", "", "description of this accelerator")
	cmd.Flags().StringVar(&uo.DisplayName, "display-name", "", "display name for the accelerator")
	cmd.Flags().StringVar(&uo.IconUrl, "icon-url", "", "URL for icon to use with the accelerator")
	cmd.Flags().StringSliceVar(&uo.Tags, "tags", []string{}, "tags that can be used to search for accelerators")
	cmd.Flags().StringVar(&uo.GitRepoUrl, "git-repository", "", "Git repository URL for the accelerator")
	cmd.Flags().StringVar(&uo.GitBranch, "git-branch", "", "Git repository branch to be used")
	cmd.Flags().StringVar(&uo.GitTag, "git-tag", "", "Git repository tag to be used")
	cmd.Flags().StringVar(&uo.GitSubPath, "git-sub-path", "", "Git repository subPath to be used")
	cmd.Flags().BoolVar(&uo.Reconcile, "reconcile", false, "trigger a reconciliation including the associated GitRepository resource")
	cmd.Flags().StringVar(&uo.Interval, "interval", "", "interval for checking for updates to Git or image repository")
	cmd.Flags().StringVar(&uo.SourceImage, "source-image", "", "(DEPRECATED) name of the source image for the accelerator")
	cmd.Flags().StringVar(&uo.SecretRef, "secret-ref", "", "name of secret containing credentials for private Git or image repository")
	cmd.Flags().SetNormalizeFunc(normalizeGitRepoRun)
}

type FragmentUpdateOptions struct {
	Namespace   string
	DisplayName string
	GitRepoUrl  string
	GitBranch   string
	GitTag      string
	GitSubPath  string
	Interval    string
	// Deprecated: SourceImage is deprecated
	SourceImage string
	SecretRef   string
	Reconcile   bool
}

func (uo *FragmentUpdateOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cmd.Flags().StringVarP(&uo.Namespace, "namespace", "n", "accelerator-system", "namespace for accelerator fragments")
	cmd.Flags().StringVar(&uo.DisplayName, "display-name", "", "display name for the accelerator fragment")
	cmd.Flags().StringVar(&uo.GitRepoUrl, "git-repository", "", "Git repository URL for the accelerator fragment")
	cmd.Flags().StringVar(&uo.GitBranch, "git-branch", "", "Git repository branch to be used")
	cmd.Flags().StringVar(&uo.GitTag, "git-tag", "", "Git repository tag to be used")
	cmd.Flags().StringVar(&uo.GitSubPath, "git-sub-path", "", "Git repository subPath to be used")
	cmd.Flags().BoolVar(&uo.Reconcile, "reconcile", false, "trigger a reconciliation including the associated GitRepository resource")
	cmd.Flags().StringVar(&uo.Interval, "interval", "", "interval for checking for updates to Git repository")
	cmd.Flags().StringVar(&uo.SourceImage, "source-image", "", "(DEPRECATED) name of the source image for the accelerator fragment")
	cmd.Flags().StringVar(&uo.SecretRef, "secret-ref", "", "name of secret containing credentials for private Git repository")
	cmd.Flags().SetNormalizeFunc(normalizeGitRepoRun)
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
	cmd.Flags().StringVarP(&do.Namespace, "namespace", "n", "accelerator-system", "namespace for accelerator system")

}

type ListOptions struct {
	Tags        []string
	Namespace   string
	ServerUrl   string
	FromContext bool
	Verbose     bool
}

func (lo *ListOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cmd.Flags().StringSliceVarP(&lo.Tags, "tags", "t", []string{}, "accelerator tags to match against")
	cmd.Flags().StringVarP(&lo.Namespace, "namespace", "n", "accelerator-system", "namespace for accelerator system")
	cmd.Flags().StringVar(&lo.ServerUrl, "server-url", "", "the URL for the Application Accelerator server")
	cmd.Flags().BoolVar(&lo.FromContext, "from-context", false, "retrieve resources from current context defined in kubeconfig")
	cmd.Flags().BoolVarP(&lo.Verbose, "verbose", "v", false, "include repository and show long URLs or image digests in the output")
}

type FragmentListOptions struct {
	Namespace string
	Verbose   bool
}

func (lo *FragmentListOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cmd.Flags().StringVarP(&lo.Namespace, "namespace", "n", "accelerator-system", "namespace for accelerator system")
	cmd.Flags().BoolVarP(&lo.Verbose, "verbose", "v", false, "include repository and show long URLs or image digests in the output")
}

type GetOptions struct {
	Namespace   string
	ServerUrl   string
	FromContext bool
	Verbose     bool
}

type FragmentGetOptions struct {
	Namespace string
}

func (gopts *GetOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cmd.Flags().StringVarP(&gopts.Namespace, "namespace", "n", "accelerator-system", "namespace for accelerator system")
	cmd.Flags().StringVar(&gopts.ServerUrl, "server-url", "", "the URL for the Application Accelerator server")
	cmd.Flags().BoolVar(&gopts.FromContext, "from-context", false, "retrieve resources from current context defined in kubeconfig")
	cmd.Flags().BoolVarP(&gopts.Verbose, "verbose", "v", false, "include all fields and show long URLs in the output")
}

func (gopts *FragmentGetOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cmd.Flags().StringVarP(&gopts.Namespace, "namespace", "n", "accelerator-system", "namespace for accelerator system")
}

type ApplyOptions struct {
	Namespace string
	FileName  string
}

func (appopts *ApplyOptions) DefineFlags(ctx context.Context, cmd *cobra.Command, c *cli.Config) {
	cmd.Flags().StringVarP(&appopts.Namespace, "namespace", "n", "accelerator-system", "namespace for the resource")
	cmd.Flags().StringVarP(&appopts.FileName, "filename", "f", "", "path of manifest file for the resource")
	cmd.MarkFlagRequired("filename")
}
