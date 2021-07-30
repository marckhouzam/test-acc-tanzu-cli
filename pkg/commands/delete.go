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

func DeleteCmd(clientset *acceleratorClientSet.AcceleratorV1Alpha1Client) *cobra.Command {
	opts := DeleteOptions{}
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete accelerator",
		Long:  `Delete will delete an accelerator from the given name`,
		Run: func(cmd *cobra.Command, args []string) {
			err := clientset.Accelerators(opts.Namespace).Delete(context.Background(), args[0], v1.DeleteOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("deleted accelerator %s in namespace %s\n", args[0], opts.Namespace)
		},
	}
	opts.DefineFlags(deleteCmd)
	return deleteCmd
}
