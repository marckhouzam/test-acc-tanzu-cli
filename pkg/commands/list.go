/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/spf13/cobra"
	cli "github.com/vmware-tanzu-private/tanzu-cli-apps-plugins/pkg/cli-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ListCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	var accServerUrl string
	opts := ListOptions{}
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List accelerators",
		Long: `List all accelerators.

You can choose to list the accelerators from the Application Accelerator server using --server-url flag
or from a Kubernetes context using --from-context flag. The default is to list accelerators from the
Application Acceleratior server and you can set the ACC_SERVER_URL environment variable with the URL for
the Application Acceleratior server you want to access.
`,
		Example: "tanzu accelerator list",
		RunE: func(cmd *cobra.Command, args []string) error {
			var context, kubeconfig bool
			if cmd.Parent() != nil {
				context = cmd.Parent().PersistentFlags().Changed("context")
				kubeconfig = cmd.Parent().PersistentFlags().Changed("kubeconfig")
			}
			serverUrl := accServerUrl
			if opts.ServerUrl != "" {
				serverUrl = opts.ServerUrl
			}
			w := new(tabwriter.Writer)
			w.Init(cmd.OutOrStdout(), 0, 8, 3, ' ', 0)
			if !opts.FromContext && !context && !kubeconfig {
				return printListFromUiServer(serverUrl, w, cmd)
			} else {
				return printListFromClient(ctx, c, opts, cmd, w)
			}
		},
	}
	accServerUrl = EnvVar("ACC_SERVER_URL", "http://localhost:8877")
	opts.DefineFlags(ctx, listCmd, c)
	return listCmd
}

func printListFromUiServer(url string, w *tabwriter.Writer, cmd *cobra.Command) error {
	Accelerators, err := GetAcceleratorsFromUiServer(url, cmd)
	if err != nil {
		return err
	}
	sort.Slice(Accelerators, func(i, j int) bool {
		return strings.Compare(Accelerators[i].Name, Accelerators[j].Name) < 0
	})
	fmt.Fprintln(w, "NAME\tGIT REPOSITORY\tBRANCH\tTAG")
	for _, accelerator := range Accelerators {
		gitRepoUrl := accelerator.SpecGitRepositoryUrl
		if gitRepoUrl == "" {
			gitRepoUrl = accelerator.SourceUrl
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", accelerator.Name, gitRepoUrl, accelerator.SourceBranch, accelerator.SourceTag)
	}
	w.Flush()

	return nil
}

func printListFromClient(ctx context.Context, c *cli.Config, opts ListOptions, cmd *cobra.Command, w *tabwriter.Writer) error {
	accelerators := &acceleratorv1alpha1.AcceleratorList{}
	err := c.List(ctx, accelerators, client.InNamespace(opts.Namespace))
	if err != nil {
		fmt.Fprintf(cmd.OutOrStderr(), "There was an error listing accelerators\n")
		return err
	}
	if len(accelerators.Items) == 0 {
		errorMsg := "no accelerators found"
		fmt.Fprintf(cmd.OutOrStderr(), errorMsg+".\n")
		return errors.New(errorMsg)
	}
	fmt.Fprintln(w, "NAME\tGIT REPOSITORY\tBRANCH\tTAG\tIMAGE")
	for _, accelerator := range accelerators.Items {
		values := []string{accelerator.Name}

		if accelerator.Spec.Git != nil {
			values = append(values, accelerator.Spec.Git.URL, accelerator.Spec.Git.Reference.Branch, accelerator.Spec.Reference.Tag)
		} else {
			values = append(values, "", "", "")
		}

		if accelerator.Spec.Source != nil {
			values = append(values, accelerator.Spec.Source.Image)
		} else {
			values = append(values, "")
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", values[0], values[1], values[2], values[3], values[4])
	}
	w.Flush()
	return nil
}
