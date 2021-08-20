/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"
	"fmt"

	"github.com/imdario/mergo"
	acceleratorClientSet "github.com/pivotal/acc-controller/api/clientset"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	fluxcdv1beta1 "github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func UpdateCmd(clientset *acceleratorClientSet.AcceleratorV1Alpha1Client) *cobra.Command {
	opts := UpdateOptions{}
	var updateCmd = &cobra.Command{
		Use:     "update",
		Short:   "Update accelerator",
		Long:    `Update accelerator`,
		Example: "tanzu accelerator update <accelerator-name> --description \"Lorem Ipsum\"",
		Run: func(cmd *cobra.Command, args []string) {
			accelerator, err := clientset.Accelerators(opts.Namespace).Get(context.Background(), args[0], v1.GetOptions{})
			if err != nil {
				panic(err.Error())
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
			updatedAcceleratorStruct := *updatedAccelerator
			err = mergo.Merge(&updatedAcceleratorStruct, *accelerator)
			if err != nil {
				panic(err.Error())
			}
			clientset.Accelerators(opts.Namespace).Update(context.Background(), &updatedAcceleratorStruct, v1.UpdateOptions{})
			fmt.Printf("updated accelerator %s in namespace %s\n", args[0], opts.Namespace)
		},
	}
	opts.DefineFlags(updateCmd)
	return updateCmd
}
