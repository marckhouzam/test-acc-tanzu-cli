/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"
	"fmt"

	acceleratorClientSet "github.com/pivotal/acc-controller/api/clientset"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListCmd(clientset *acceleratorClientSet.AcceleratorV1Alpha1Client) *cobra.Command {
	opts := ListOptions{}
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List accelerators",
		Long:  `List the accelerators, you can choose with namespace to use passing the flag -namespace`,
		Run: func(cmd *cobra.Command, args []string) {
			accelerators, err := clientset.Accelerators(opts.Namespace).List(context.Background(), v1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}
			for _, accelerator := range accelerators.Items {
				fmt.Println(accelerator.Name)
			}
		},
	}
	opts.DefineFlags(listCmd)
	return listCmd
}
