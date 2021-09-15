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
	var accServerUrl string
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
				return printAcceleratorFromUiServer(serverUrl, args[0], w, cmd)
			} else {
				return printAcceleratorFromClient(ctx, opts, cmd, args[0], w, c)
			}
		},
	}
	accServerUrl = EnvVar("ACC_SERVER_URL", "http://localhost:8877")
	opts.DefineFlags(ctx, getCmd, c)
	return getCmd
}

func printAcceleratorFromUiServer(url string, name string, w *tabwriter.Writer, cmd *cobra.Command) error {
	errorMsg := "accelertor %s not found"
	Accelerators, err := GetAcceleratorsFromUiServer(url, cmd)
	if err != nil {
		return err
	}
	for _, accelerator := range Accelerators {
		if accelerator.Name == name {
			gitRepoUrl := accelerator.SpecGitRepositoryUrl
			if gitRepoUrl == "" {
				gitRepoUrl = accelerator.SourceUrl
			}
			fmt.Fprintln(w, "NAME\tGIT REPOSITORY\tBRANCH\tTAG")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", accelerator.Name, gitRepoUrl, accelerator.SourceBranch, accelerator.SourceTag)
			w.Flush()
			return nil
		}
	}

	fmt.Fprintf(cmd.OutOrStderr(), errorMsg+".\n", name)
	return fmt.Errorf(errorMsg, name)
}

func printAcceleratorFromClient(ctx context.Context, opts GetOptions, cmd *cobra.Command, name string, w *tabwriter.Writer, c *cli.Config) error {
	accelerator := &acceleratorv1alpha1.Accelerator{}
	err := c.Get(ctx, client.ObjectKey{Namespace: opts.Namespace, Name: name}, accelerator)
	if err != nil {
		fmt.Fprintf(cmd.OutOrStderr(), "Error getting accelerator %s\n", name)
		return err
	}
	fmt.Fprintln(w, "NAME\tGIT REPOSITORY\tBRANCH\tTAG")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", accelerator.Name, accelerator.Spec.Git.URL, accelerator.Spec.Git.Reference.Branch, accelerator.Spec.Git.Reference.Tag)
	w.Flush()
	return nil
}
