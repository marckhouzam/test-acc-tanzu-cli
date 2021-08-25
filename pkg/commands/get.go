/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"
	"fmt"
	"text/tabwriter"

	acceleratorClientSet "github.com/pivotal/acc-controller/api/clientset"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetCmd(clientset acceleratorClientSet.AcceleratorV1Alpha1Interface) *cobra.Command {
	opts := GetOptions{}
	var getCmd = &cobra.Command{
		Use:     "get",
		Short:   "Get accelerator",
		Long:    `Get accelerator`,
		Example: "tanzu accelerator get <accelerator-name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			accelerator, err := clientset.Accelerators(opts.Namespace).Get(context.Background(), args[0], v1.GetOptions{})
			w := new(tabwriter.Writer)
			w.Init(cmd.OutOrStdout(), 0, 8, 3, ' ', 0)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error getting accelerator %s", args[0])
				return err
			}
			fmt.Fprintln(w, "NAME\tGIT REPOSITORY\tBRANCH\tTAG")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", accelerator.Name, accelerator.Spec.Git.URL, accelerator.Spec.Git.Reference.Branch, accelerator.Spec.Git.Reference.Tag)
			w.Flush()
			return nil
		},
	}
	opts.DefineFlags(getCmd)
	return getCmd
}
