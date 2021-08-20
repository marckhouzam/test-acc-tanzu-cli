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

func ListCmd(clientset *acceleratorClientSet.AcceleratorV1Alpha1Client, w *tabwriter.Writer) *cobra.Command {
	opts := ListOptions{}
	var listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List accelerators",
		Long:    `List the accelerators, you can choose with namespace to use passing the flag -namespace`,
		Example: "tanzu accelerator list",
		Run: func(cmd *cobra.Command, args []string) {
			accelerators, err := clientset.Accelerators(opts.Namespace).List(context.Background(), v1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Fprintln(w, "NAME\tGIT REPOSITORY\tBRANCH\tTAG")
			for _, accelerator := range accelerators.Items {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", accelerator.Name, accelerator.Spec.Git.URL, accelerator.Spec.Git.Reference.Branch, accelerator.Spec.Git.Reference.Tag)
			}
			w.Flush()
		},
	}
	opts.DefineFlags(listCmd)
	return listCmd
}
