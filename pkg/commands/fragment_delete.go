/*
Copyright 2021 VMware, Inc. All Rights Reserved.
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

func FragmentDeleteCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := DeleteOptions{}
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete an accelerator fragment",
		Long:  `Delete the accelerator fragment resource with the specified name.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must specify the name of the accelerator fragment")
			}
			return nil
		},
		ValidArgsFunction: SuggestFragmentNamesFromConfig(context.Background(), c),
		Example:           "tanzu accelerator fragment delete <fragment-name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			fragment := &acceleratorv1alpha1.Fragment{}
			err := c.Get(ctx, client.ObjectKey{Namespace: opts.Namespace, Name: args[0]}, fragment)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "accelerator fragment %s not found\n", args[0])
				return err
			}
			err = c.Delete(ctx, fragment)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "There was a problem trying to delete accelerator fragment %s\n", args[0])
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "deleted accelerator fragment %s in namespace %s\n", args[0], opts.Namespace)
			return nil
		},
	}
	opts.DefineFlags(ctx, deleteCmd, c)
	return deleteCmd
}
