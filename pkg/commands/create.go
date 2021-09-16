/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"
	"errors"
	"fmt"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	fluxcdv1beta1 "github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	"github.com/spf13/cobra"
	cli "github.com/vmware-tanzu-private/tanzu-cli-apps-plugins/pkg/cli-runtime"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := CreateOptions{}
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new accelerator",
		Long:    `Create a new accelerator resource using the provided options`,
		Example: "tanzu accelerator create <accelerator-name> -git-repository <git-repo-URL>",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must specify the name of the accelerator")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			acc := &acceleratorv1alpha1.Accelerator{
				TypeMeta: v1.TypeMeta{
					APIVersion: "accelerator.tanzu.vmware.com/v1alpha1",
					Kind:       "Accelerator",
				},
				ObjectMeta: v1.ObjectMeta{
					Namespace: opts.Namespace,
					Name:      args[0],
				},
				Spec: acceleratorv1alpha1.AcceleratorSpec{
					DisplayName: opts.DisplayName,
					Description: opts.Description,
					IconUrl:     opts.IconUrl,
					Tags:        opts.Tags,
					Git: acceleratorv1alpha1.Git{
						URL: opts.GitRepoUrl,
						Reference: &fluxcdv1beta1.GitRepositoryRef{
							Branch: opts.GitBranch,
							Tag:    opts.GitTag,
						},
					},
				},
			}
			err := c.Create(ctx, acc)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error creating accelerator %s\n", args[0])
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "created accelerator %s in namespace %s\n", acc.Name, acc.Namespace)
			return nil

		},
	}
	opts.DefineFlags(ctx, cmd, c)
	return cmd
}
