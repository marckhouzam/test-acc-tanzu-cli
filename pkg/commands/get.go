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

func GetCmd(clientset *acceleratorClientSet.AcceleratorV1Alpha1Client, w *tabwriter.Writer) *cobra.Command {
	opts := GetOptions{}
	var getCmd = &cobra.Command{
		Use:     "get",
		Short:   "Get accelerator",
		Long:    `Get accelerator`,
		Example: "tanzu accelerator get <accelerator-name>",
		Run: func(cmd *cobra.Command, args []string) {
			accelerator, err := clientset.Accelerators(opts.Namespace).Get(context.Background(), args[0], v1.GetOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Fprintln(w, "NAME\tGIT REPOSITORY\tBRANCH")
			fmt.Fprintf(w, "%s\t%s\t%s\n", accelerator.Name, accelerator.Spec.Git.URL, accelerator.Spec.Git.Reference.Branch)
			w.Flush()
		},
	}
	opts.DefineFlags(getCmd)
	return getCmd
}
