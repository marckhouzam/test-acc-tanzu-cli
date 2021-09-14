/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"
	"errors"
	"fmt"
	"text/tabwriter"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu-private/tanzu-cli-apps-plugins/pkg/cli-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := GetOptions{}
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get accelerator",
		Long:  `Get accelerator`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must pass the name of the accelerator")
			}
			return nil
		},
		Example: "tanzu accelerator get <accelerator-name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			accelerator := &acceleratorv1alpha1.Accelerator{}
			err := c.Get(ctx, client.ObjectKey{Namespace: opts.Namespace, Name: args[0]}, accelerator)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error getting accelerator %s\n", args[0])
				return err
			}
			w := new(tabwriter.Writer)
			w.Init(cmd.OutOrStdout(), 0, 8, 3, ' ', 0)
			fmt.Fprintln(w, "NAME\tGIT REPOSITORY\tBRANCH\tTAG")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", accelerator.Name, accelerator.Spec.Git.URL, accelerator.Spec.Git.Reference.Branch, accelerator.Spec.Git.Reference.Tag)
			w.Flush()
			return nil
		},
	}
	opts.DefineFlags(ctx, getCmd, c)
	return getCmd
}
