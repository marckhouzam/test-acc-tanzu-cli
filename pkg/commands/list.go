/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"
	"fmt"
	"text/tabwriter"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/spf13/cobra"
	cli "github.com/vmware-tanzu-private/tanzu-cli-apps-plugins/pkg/cli-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ListCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := ListOptions{}
	var listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List accelerators",
		Long:    `List the accelerators, you can choose with namespace to use passing the flag -namespace`,
		Example: "tanzu accelerator list",
		RunE: func(cmd *cobra.Command, args []string) error {
			accelerators := &acceleratorv1alpha1.AcceleratorList{}
			err := c.List(ctx, accelerators, client.InNamespace(opts.Namespace))
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "There was an error listing accelerators\n")
				return err
			}
			if len(accelerators.Items) == 0 {
				fmt.Fprintf(cmd.OutOrStderr(), "No accelerators found.\n")
				return nil
			}
			w := new(tabwriter.Writer)
			w.Init(cmd.OutOrStdout(), 0, 8, 3, ' ', 0)
			fmt.Fprintln(w, "NAME\tGIT REPOSITORY\tBRANCH\tTAG")
			for _, accelerator := range accelerators.Items {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", accelerator.Name, accelerator.Spec.Git.URL, accelerator.Spec.Git.Reference.Branch, accelerator.Spec.Git.Reference.Tag)
			}
			w.Flush()
			return nil
		},
	}
	opts.DefineFlags(ctx, listCmd, c)
	return listCmd
}
