/*
Copyright 2021-2023 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"
	"errors"
	"fmt"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/spf13/cobra"
	cli "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func DeleteCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := DeleteOptions{}
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete an accelerator",
		Long:  `Delete the accelerator resource with the specified name.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must specify the name of the accelerator")
			}
			return nil
		},
		ValidArgsFunction: SuggestAcceleratorNamesFromConfig(context.Background(), c),
		Example:           "tanzu accelerator delete <accelerator-name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			accelerator := &acceleratorv1alpha1.Accelerator{}
			err := c.Get(ctx, client.ObjectKey{Namespace: opts.Namespace, Name: args[0]}, accelerator)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "accelerator %s not found\n", args[0])
				return err
			}
			err = c.Delete(ctx, accelerator)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "There was a problem trying to delete accelerator %s\n", args[0])
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "deleted accelerator %s in namespace %s\n", args[0], opts.Namespace)
			return nil
		},
	}
	opts.DefineFlags(ctx, deleteCmd, c)
	return deleteCmd
}
