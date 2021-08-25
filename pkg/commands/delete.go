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

func DeleteCmd(clientset acceleratorClientSet.AcceleratorV1Alpha1Interface) *cobra.Command {
	opts := DeleteOptions{}
	var deleteCmd = &cobra.Command{
		Use:     "delete",
		Short:   "Delete accelerator",
		Example: "tanzu accelerator delete <accelerator-name>",
		Long:    `Delete will delete an accelerator from the given name`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := clientset.Accelerators(opts.Namespace).Delete(context.Background(), args[0], v1.DeleteOptions{})
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "There was a problem trying to delete accelerator %s", args[0])
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "deleted accelerator %s in namespace %s\n", args[0], opts.Namespace)
			return nil
		},
	}
	opts.DefineFlags(deleteCmd)
	return deleteCmd
}
