/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/imdario/mergo"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	fluxcdv1beta1 "github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu-private/tanzu-cli-apps-plugins/pkg/cli-runtime"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func UpdateCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := UpdateOptions{}
	requestedAtAnnotation := "reconcile.accelerator.apps.tanzu.vmware.com/requestedAt"
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update accelerator",
		Long:  `Update accelerator`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must pass the name of the accelerator")
			}
			return nil
		},
		Example: "tanzu accelerator update <accelerator-name> --description \"Lorem Ipsum\"",
		RunE: func(cmd *cobra.Command, args []string) error {
			accelerator := &acceleratorv1alpha1.Accelerator{}
			err := c.Get(context.Background(), client.ObjectKey{Namespace: opts.Namespace, Name: args[0]}, accelerator)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "accelerator %s not found\n", args[0])
				return err
			}
			updatedAccelerator := &acceleratorv1alpha1.Accelerator{
				TypeMeta: v1.TypeMeta{
					APIVersion: "accelerator.tanzu.vmware.com/v1alpha1",
					Kind:       "Accelerator",
				},
				ObjectMeta: v1.ObjectMeta{
					Namespace: opts.Namespace,
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
						},
					},
				},
			}
			if opts.Reconcile {
				if accelerator.ObjectMeta.Annotations == nil {
					accelerator.ObjectMeta.Annotations = make(map[string]string)
				}
				accelerator.ObjectMeta.Annotations[requestedAtAnnotation] = time.Now().UTC().Format(time.RFC3339)
			}
			mergo.Merge(updatedAccelerator, *accelerator)
			err = c.Update(ctx, updatedAccelerator)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "there was an error updating accelerator %s\n", args[0])
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "accelerator %s updated successfully\n", args[0])
			return nil
		},
	}
	opts.DefineFlags(ctx, updateCmd, c)
	return updateCmd
}
