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
	"github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
	"gopkg.in/yaml.v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func FragmentGetCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := FragmentGetOptions{}
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get accelerator fragment info",
		Long:  "Get accelerator fragment info.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must specify the name of the fragment")
			}
			return nil
		},
		Example:           "tanzu accelerator get <fragment-name>",
		ValidArgsFunction: getSuggestion(ctx, c),
		RunE: func(cmd *cobra.Command, args []string) error {
			w := new(tabwriter.Writer)
			w.Init(cmd.OutOrStdout(), 0, 8, 3, ' ', 0)
			return printFragmentFromClient(ctx, opts, cmd, args[0], w, c)
		},
	}
	opts.DefineFlags(ctx, getCmd, c)
	return getCmd
}

func printFragmentFromClient(ctx context.Context, opts FragmentGetOptions, cmd *cobra.Command, name string, w *tabwriter.Writer, c *cli.Config) error {
	fragment := &acceleratorv1alpha1.Fragment{}
	err := c.Get(ctx, client.ObjectKey{Namespace: opts.Namespace, Name: name}, fragment)
	if err != nil {
		fmt.Fprintf(cmd.OutOrStderr(), "Error getting accelerator fragment %s\n", name)
		return err
	}
	var options []interface{}
	yaml.Unmarshal([]byte(fragment.Status.Options), &options)
	optionsYaml, _ := yaml.Marshal(options)
	fmt.Fprintf(cmd.OutOrStdout(), "name: %s\n", fragment.Name)
	fmt.Fprintf(cmd.OutOrStdout(), "namespace: %s\n", fragment.Namespace)
	fmt.Fprintf(cmd.OutOrStdout(), "displayName: %s\n", fragment.Status.DisplayName)
	if fragment.Spec.Source != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "source:\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  image: %s\n", fragment.Spec.Source.Image)
		if fragment.Spec.Source.ImagePullSecrets != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "  secret-ref: %s\n", fragment.Spec.Source.ImagePullSecrets)
		}
	} else {
		if fragment.Spec.Git != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "git:")
			if fragment.Spec.Git.Interval != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "  interval: %s\n", fragment.Spec.Git.Interval.Duration)
			}
			if fragment.Spec.Git.Ignore != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "  ignore: %s\n", *fragment.Spec.Ignore)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  url: %s\n", fragment.Spec.Git.URL)
			fmt.Fprintf(cmd.OutOrStdout(), "  ref:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "    branch: %s\n", fragment.Spec.Git.Reference.Branch)
			if fragment.Spec.Git.Reference.Tag != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "    tag: %s\n", fragment.Spec.Git.Reference.Tag)
			}
			if fragment.Spec.Git.SubPath != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "  subPath: %s\n", *fragment.Spec.Git.SubPath)
			}
			if fragment.Spec.Git.SecretRef != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "  secret-ref: %s\n", fragment.Spec.Git.SecretRef.Name)
			}
		}
	}
	if len(fragment.Status.Conditions) > 0 {
		ready := false
		reason := ""
		message := ""
		for i := 0; i < len(fragment.Status.Conditions); i++ {
			if fragment.Status.Conditions[i].Type == "Ready" {
				reason = fragment.Status.Conditions[i].Reason
				message = fragment.Status.Conditions[i].Message
				if fragment.Status.Conditions[i].Status == "True" {
					ready = true
				}
				break
			}
		}
		fmt.Fprintf(cmd.OutOrStdout(), "ready: %t\n", ready)
		if !ready {
			fmt.Fprintf(cmd.OutOrStdout(), "reason: %s\n", reason)
			fmt.Fprintf(cmd.OutOrStdout(), "message: %s\n", message)
		}
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "ready: %s\n", "true")
	}
	if string(optionsYaml) != "[]\n" {
		fmt.Fprintln(cmd.OutOrStdout(), "options:")
		fmt.Fprint(cmd.OutOrStdout(), string(optionsYaml))
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "options: %s", string(optionsYaml))
	}
	fmt.Fprintln(cmd.OutOrStdout(), "artifact:")
	fmt.Fprintf(cmd.OutOrStdout(), "  message: %s\n", fragment.Status.ArtifactInfo.Message)
	fmt.Fprintf(cmd.OutOrStdout(), "  ready: %t\n", fragment.Status.ArtifactInfo.Ready)
	fmt.Fprintf(cmd.OutOrStdout(), "  url: %s\n", fragment.Status.ArtifactInfo.URL)

	fmt.Fprintln(cmd.OutOrStdout(), "imports:")
	imports := fragment.Status.ArtifactInfo.Imports
	if len(imports) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "  None\n")
	} else {
		for key, _ := range imports {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", key)
		}
	}

	fmt.Fprintln(cmd.OutOrStdout(), "importedBy:")
	accelerators := &acceleratorv1alpha1.AcceleratorList{}
	err = c.List(ctx, accelerators, client.InNamespace(opts.Namespace), client.HasLabels{"imports.accelerator.apps.tanzu.vmware.com/" + fragment.Name})
	if err != nil {
		fmt.Fprintf(cmd.OutOrStderr(), "  Unable to find any importing accelerators\n")
	} else {
		if len(accelerators.Items) > 0 {
			for _, accelerator := range accelerators.Items {
				fmt.Fprintf(cmd.OutOrStderr(), "  accelerator/%s\n", accelerator.Name)
			}
		}
	}
	fragments := &acceleratorv1alpha1.FragmentList{}
	err = c.List(ctx, fragments, client.InNamespace(opts.Namespace), client.HasLabels{"imports.accelerator.apps.tanzu.vmware.com/" + fragment.Name})
	if err != nil {
		fmt.Fprintf(cmd.OutOrStderr(), "  Unable to find any importing fragments\n")
	} else {
		if len(fragments.Items) > 0 {
			for _, fragments := range fragments.Items {
				fmt.Fprintf(cmd.OutOrStderr(), "  fragment/%s\n", fragments.Name)
			}
		}
	}
	if len(accelerators.Items) == 0 && len(fragments.Items) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "  None\n")
	}

	w.Flush()
	return nil
}
