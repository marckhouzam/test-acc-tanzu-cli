/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"
	"fmt"

	acceleratorClientSet "github.com/pivotal/acc-controller/api/clientset"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	fluxcdv1beta1 "github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateCmd(clientset *acceleratorClientSet.AcceleratorV1Alpha1Client) *cobra.Command {
	opts := CreateOptions{}
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create accelerator",
		Example: "tanzu accelerator create <accelerator-name> -git-repository <git-repo-URL>",
		Long:    `Create will add the accelerator resource using the given parameters`,
		Run: func(cmd *cobra.Command, args []string) {
			if args[0] == "" {
				panic("you need to pass the name of the accelerator")
			}
			if opts.DisplayName == "" {
				opts.DisplayName = args[0]
			}

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
						},
					},
				},
			}
			result, err := clientset.Accelerators(opts.Namespace).Create(context.Background(), acc, v1.CreateOptions{})
			if err != nil {
				panic(err.Error())
			}

			fmt.Printf("created accelerator %s in namespace %s\n", result.Name, result.Namespace)

		},
	}
	opts.DefineFlags(cmd)
	return cmd
}
