/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"github.com/spf13/cobra"
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

func (co *CreateOptions) DefineFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&co.Namespace, "namespace", "n", "default", "Kubernetes namespace")
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
}

func (uo *UpdateOptions) DefineFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&uo.Namespace, "namespace", "n", "default", "Kubernetes namespace")
	cmd.Flags().StringVar(&uo.Description, "description", "", "Accelerator description")
	cmd.Flags().StringVar(&uo.DisplayName, "display-name", "", "Accelerator display name")
	cmd.Flags().StringVar(&uo.IconUrl, "icon-url", "", "Accelerator icon location")
	cmd.Flags().StringSliceVar(&uo.Tags, "tags", []string{}, "Accelerator Tags")
	cmd.Flags().StringVar(&uo.GitRepoUrl, "git-repository", "", "Accelerator repo URL")
	cmd.Flags().StringVar(&uo.GitBranch, "git-branch", "main", "Accelerator repo branch")
}

type DeleteOptions struct {
	Namespace string
}

func (do *DeleteOptions) DefineFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&do.Namespace, "namespace", "n", "default", "Kubernetes namespace")
}

type ListOptions struct {
	Namespace string
}

func (lo *ListOptions) DefineFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&lo.Namespace, "namespace", "n", "default", "Kubernetes namespace")
}

type GetOptions struct {
	Namespace string
}

func (gopts *GetOptions) DefineFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&gopts.Namespace, "namespace", "n", "default", "Kubernetes namespace")
}
