/*
Copyright 2022 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"
	"fmt"
	"text/tabwriter"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/spf13/cobra"
	cli "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func FragmentListCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := FragmentListOptions{}
	var fragmentListCmd = &cobra.Command{
		Use:     "list",
		Short:   "List accelerator fragments",
		Long:    "List all accelerator fragments.",
		Example: "tanzu accelerator fragment list",
		RunE: func(cmd *cobra.Command, args []string) error {
			w := new(tabwriter.Writer)
			w.Init(cmd.OutOrStdout(), 0, 8, 3, ' ', 0)
			return printFragmentListFromClient(ctx, c, opts, cmd, w)
		},
	}
	opts.DefineFlags(ctx, fragmentListCmd, c)
	return fragmentListCmd
}

func printFragmentListFromClient(ctx context.Context, c *cli.Config, opts FragmentListOptions, cmd *cobra.Command, w *tabwriter.Writer) error {
	fragments := &acceleratorv1alpha1.FragmentList{}
	err := c.List(ctx, fragments, client.InNamespace(opts.Namespace))
	if err != nil {
		fmt.Fprintf(cmd.OutOrStderr(), "There was an error listing accelerator fragments\n")
		return err
	}

	fragList := [][]string{}

	for _, fragment := range fragments.Items {
		values := []string{fragment.Name}

		status := "unknown"
		for _, cond := range fragment.Status.Conditions {
			if cond.Type == "Ready" {
				if cond.Status == "True" {
					status = "true"
				} else {
					status = "false"
				}
				break
			}
		}

		repo := ""
		if fragment.Spec.Git != nil {
			repo = fragment.Spec.Git.URL
			if fragment.Spec.Git.Reference.Tag != "" {
				repo = repo + ":" + fragment.Spec.Git.Reference.Tag
			} else if fragment.Spec.Git.Reference.Branch != "" {
				repo = repo + ":" + fragment.Spec.Git.Reference.Branch
			}
			if fragment.Spec.Git.SubPath != nil {
				repo = repo + ":/" + *fragment.Spec.Git.SubPath
			}
			values = append(values, repo, status)
		} else if fragment.Spec.Source != nil {
			repo = "source-image: " + fragment.Spec.Source.Image
			values = append(values, repo, status)
		} else {
			values = append(values, "", "")
		}
		fragList = append(fragList, values)
	}

	printAcceleratorFragmentList(c, cmd, w, fragList)
	return nil
}

func printAcceleratorFragmentList(c *cli.Config, cmd *cobra.Command, w *tabwriter.Writer, fragments [][]string) {
	if len(fragments) == 0 {
		c.Infof("No accelerator fragments found.\n")
	} else {
		fmt.Fprintln(w, "NAME\tREADY\tREPOSITORY")
		for _, values := range fragments {
			fmt.Fprintf(w, "%s\t%s\t%s\n", values[0], values[2], values[1])
		}
		w.Flush()
	}
}
