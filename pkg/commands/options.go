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
	cmd.Flags().StringVar(&co.Namespace, "namespace", "default", "Kubernetes namespace")
	cmd.Flags().StringVar(&co.Description, "description", "", "Accelerator description")
	cmd.Flags().StringVar(&co.DisplayName, "displayName", "", "Accelerator display name")
	cmd.Flags().StringVar(&co.IconUrl, "iconUrl", "", "Accelerator icon location")
	cmd.Flags().StringSliceVar(&co.Tags, "tags", []string{}, "Accelerator Tags")
	cmd.Flags().StringVar(&co.GitRepoUrl, "gitRepoUrl", "", "Accelerator repo URL")
	cmd.Flags().StringVar(&co.GitBranch, "gitBranch", "", "Accelerator repo branch")

	cmd.MarkFlagRequired("gitRepoUrl")
	cmd.MarkFlagRequired("gitBranch")
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
	cmd.Flags().StringVar(&uo.Namespace, "namespace", "default", "Kubernetes namespace")
	cmd.Flags().StringVar(&uo.Description, "description", "", "Accelerator description")
	cmd.Flags().StringVar(&uo.DisplayName, "displayName", "", "Accelerator display name")
	cmd.Flags().StringVar(&uo.IconUrl, "iconUrl", "", "Accelerator icon location")
	cmd.Flags().StringSliceVar(&uo.Tags, "tags", []string{}, "Accelerator Tags")
	cmd.Flags().StringVar(&uo.GitRepoUrl, "gitRepoUrl", "", "Accelerator repo URL")
	cmd.Flags().StringVar(&uo.GitBranch, "gitBranch", "", "Accelerator repo branch")
}

type DeleteOptions struct {
	Namespace string
}

func (do *DeleteOptions) DefineFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&do.Namespace, "namespace", "default", "Kubernetes namespace")
}

type ListOptions struct {
	Namespace string
}

func (lo *ListOptions) DefineFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&lo.Namespace, "namespace", "default", "Kubernetes namespace")
}

type GetOptions struct {
	Namespace string
}

func (gopts *GetOptions) DefineFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&gopts.Namespace, "namespace", "default", "Kubernetes namespace")
}
