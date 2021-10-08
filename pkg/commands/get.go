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
	"gopkg.in/yaml.v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type GitData struct {
	URL    string
	Branch string
	Tag    string
}
type GetOutput struct {
	Name        string
	Namespace   string
	Description string
	DisplayName string
	Options     string
	Tags        []string
	Ready       bool
	Git         GitData
}

func CastAcceleratorSpecToGetOutput(acc acceleratorv1alpha1.Accelerator) GetOutput {
	return GetOutput{
		Name:        acc.Name,
		Namespace:   acc.Namespace,
		Description: acc.Spec.Description,
		DisplayName: acc.Spec.DisplayName,
		Options:     acc.Status.Options,
		Tags:        acc.Spec.Tags,
		Ready:       acc.Status.ArtifactInfo.Ready,
		Git: GitData{
			URL:    acc.Spec.Git.URL,
			Branch: acc.Spec.Git.Reference.Branch,
			Tag:    acc.Spec.Git.Reference.Tag,
		},
	}
}

func getSuggestion(ctx context.Context, c *cli.Config) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if cmd.Flags().Changed("from-context") ||
			cmd.Parent().PersistentFlags().Changed("context") ||
			cmd.Parent().PersistentFlags().Changed("kubeconfig") {
			return SuggestAcceleratorNamesFromConfig(ctx, c)(cmd, args, toComplete)
		} else {
			return SuggestAcceleratorNamesFromUiServer(ctx)(cmd, args, toComplete)
		}
	}
}

func GetCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	var accServerUrl string
	opts := GetOptions{}
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get accelerator info",
		Long: `Get accelerator info.

You can choose to get the accelerator from the Application Accelerator server using --server-url flag
or from a Kubernetes context using --from-context flag. The default is to get accelerators from the
Application Acceleratior server and you can set the ACC_SERVER_URL environment variable with the URL for
the Application Acceleratior server you want to access.
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must specify the name of the accelerator")
			}
			return nil
		},
		Example:           "tanzu accelerator get <accelerator-name> --from-context",
		ValidArgsFunction: getSuggestion(ctx, c),
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
	errorMsg := "accelerator %s not found"
	Accelerators, err := GetAcceleratorsFromUiServer(url, cmd)
	if err != nil {
		return err
	}
	for _, accelerator := range Accelerators {
		if accelerator.Name == name {
			options, err := GetAcceleratorOptionsFromUiServer(url, accelerator.Name, cmd)
			if err != nil {
				return err
			}
			tagsYaml, _ := yaml.Marshal(accelerator.Tags)
			optionsYaml, _ := yaml.Marshal(options)
			fmt.Fprintf(cmd.OutOrStdout(), "name: %s\n", accelerator.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "description: %s\n", accelerator.Description)
			fmt.Fprintf(cmd.OutOrStdout(), "displayName: %s\n", accelerator.DisplayName)
			fmt.Fprintf(cmd.OutOrStdout(), "iconUrl: %s\n", accelerator.IconUrl)
			if accelerator.SpecImageRepository != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "source:\n")
				fmt.Fprintf(cmd.OutOrStdout(), "  image: %s\n", accelerator.SpecImageRepository)
			} else {
				if accelerator.SpecGitRepositoryUrl != "" {
					fmt.Fprintf(cmd.OutOrStdout(), "git:\n")
					fmt.Fprintf(cmd.OutOrStdout(), "  url: %s\n", accelerator.SpecGitRepositoryUrl)
					fmt.Fprintf(cmd.OutOrStdout(), "  ref:\n")
					fmt.Fprintf(cmd.OutOrStdout(), "    branch: %s\n", accelerator.SourceBranch)
					fmt.Fprintf(cmd.OutOrStdout(), "    tag: %s\n", accelerator.SourceTag)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "sourceUrl: %s\n", accelerator.SourceUrl)
				}
			}
			if string(tagsYaml) != "[]\n" {
				fmt.Fprintln(cmd.OutOrStdout(), "tags:")
				fmt.Fprint(cmd.OutOrStdout(), string(tagsYaml))
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "tags: %s", string(tagsYaml))
			}
			fmt.Fprintf(cmd.OutOrStdout(), "ready: %t\n", accelerator.Ready)
			if string(optionsYaml) != "[]\n" {
				fmt.Fprintln(cmd.OutOrStdout(), "options:")
				fmt.Fprint(cmd.OutOrStdout(), string(optionsYaml))
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "options: %s", string(optionsYaml))
			}
			fmt.Fprintln(cmd.OutOrStdout(), "artifact:")
			fmt.Fprintf(cmd.OutOrStdout(), "  message: %s\n", accelerator.ArchiveMessage)
			fmt.Fprintf(cmd.OutOrStdout(), "  ready: %t\n", accelerator.ArchiveReady)
			fmt.Fprintf(cmd.OutOrStdout(), "  url: %s\n", accelerator.ArchiveUrl)
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
	ignore := ""
	if accelerator.Spec.Git != nil && accelerator.Spec.Ignore != nil {
		ignore = *accelerator.Spec.Ignore
	}
	var options []interface{}
	yaml.Unmarshal([]byte(accelerator.Status.Options), &options)
	tagsYaml, _ := yaml.Marshal(accelerator.Status.Tags)
	optionsYaml, _ := yaml.Marshal(options)
	fmt.Fprintf(cmd.OutOrStdout(), "name: %s\n", accelerator.Name)
	fmt.Fprintf(cmd.OutOrStdout(), "namespace: %s\n", accelerator.Namespace)
	fmt.Fprintf(cmd.OutOrStdout(), "description: %s\n", accelerator.Status.Description)
	fmt.Fprintf(cmd.OutOrStdout(), "displayName: %s\n", accelerator.Status.DisplayName)
	fmt.Fprintf(cmd.OutOrStdout(), "iconUrl: %s\n", accelerator.Status.IconUrl)
	if accelerator.Spec.Git != nil {
		fmt.Fprintln(cmd.OutOrStdout(), "git:")
		if accelerator.Spec.Interval != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "  interval: %s\n", accelerator.Spec.Interval.Duration)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "  interval: \n")
		}
		fmt.Fprintf(cmd.OutOrStdout(), "  ignore: %s\n", ignore)
		fmt.Fprintf(cmd.OutOrStdout(), "  ref:\n")
		fmt.Fprintf(cmd.OutOrStdout(), "    branch: %s\n", accelerator.Spec.Git.Reference.Branch)
		fmt.Fprintf(cmd.OutOrStdout(), "    tag: %s\n", accelerator.Spec.Git.Reference.Tag)
		fmt.Fprintf(cmd.OutOrStdout(), "  url: %s\n", accelerator.Spec.Git.URL)
	}
	if accelerator.Spec.Source != nil {
		fmt.Fprintln(cmd.OutOrStdout(), "source:")
		fmt.Fprintf(cmd.OutOrStdout(), "  image: %s\n", accelerator.Spec.Source.Image)
	}
	if string(tagsYaml) != "[]\n" {
		fmt.Fprintln(cmd.OutOrStdout(), "tags:")
		fmt.Fprint(cmd.OutOrStdout(), string(tagsYaml))
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "tags: %s", string(tagsYaml))
	}
	fmt.Fprintf(cmd.OutOrStdout(), "ready: %t\n", accelerator.Status.ArtifactInfo.Ready)
	if string(optionsYaml) != "[]\n" {
		fmt.Fprintln(cmd.OutOrStdout(), "options:")
		fmt.Fprint(cmd.OutOrStdout(), string(optionsYaml))
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "options: %s", string(optionsYaml))
	}
	fmt.Fprintln(cmd.OutOrStdout(), "artifact:")
	fmt.Fprintf(cmd.OutOrStdout(), "  message: %s\n", accelerator.Status.ArtifactInfo.Message)
	fmt.Fprintf(cmd.OutOrStdout(), "  ready: %t\n", accelerator.Status.ArtifactInfo.Ready)
	fmt.Fprintf(cmd.OutOrStdout(), "  url: %s\n", accelerator.Status.ArtifactInfo.URL)
	w.Flush()
	return nil
}
