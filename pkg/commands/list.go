/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/spf13/cobra"
	cli "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
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
Kubernetes context. To override this, you can set the ACC_SERVER_URL environment variable with the URL for
the Application Accelerator server you want to access.
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
			if serverUrl != "" && !opts.FromContext && !context && !kubeconfig {
				return printListFromUiServer(c, serverUrl, opts, cmd, w)
			} else {
				return printListFromClient(ctx, c, opts, cmd, w)
			}
		},
	}
	accServerUrl = EnvVar("ACC_SERVER_URL", "")
	opts.DefineFlags(ctx, listCmd, c)
	return listCmd
}

func printListFromUiServer(c *cli.Config, url string, opts ListOptions, cmd *cobra.Command, w *tabwriter.Writer) error {
	accelerators, err := GetAcceleratorsFromApiServer(url, cmd)
	if err != nil {
		return err
	}
	sort.Slice(accelerators, func(i, j int) bool {
		return strings.Compare(accelerators[i].Name, accelerators[j].Name) < 0
	})

	accList := [][]string{}

	for _, accelerator := range accelerators {
		if len(opts.Tags) == 0 || contains(accelerator.Tags, opts.Tags) {
			repo := ""
			tags := fmt.Sprintf("%v", accelerator.Tags)
			if accelerator.SpecGitRepositoryUrl != "" {
				repo = accelerator.SpecGitRepositoryUrl
				if accelerator.SourceTag != "" {
					repo = repo + ":" + accelerator.SourceTag
				} else if accelerator.SourceBranch != "" {
					repo = repo + ":" + accelerator.SourceBranch
				}
			} else if accelerator.SpecImageRepository != "" {
				repo = "source-image: " + accelerator.SpecImageRepository
			}
			status := "unknown"
			if accelerator.Ready {
				status = "true"
			} else {
				status = "false"
			}

			accList = append(accList, []string{accelerator.Name, tags, repo, status})
		}
	}
	w.Flush()

	printAcceleratorList(c, cmd, w, accList)
	return nil
}

func printListFromClient(ctx context.Context, c *cli.Config, opts ListOptions, cmd *cobra.Command, w *tabwriter.Writer) error {
	accelerators := &acceleratorv1alpha1.AcceleratorList{}
	err := c.List(ctx, accelerators, client.InNamespace(opts.Namespace))
	if err != nil {
		fmt.Fprintf(cmd.OutOrStderr(), "There was an error listing accelerators\n")
		return err
	}

	accList := [][]string{}

	for _, accelerator := range accelerators.Items {

		if len(opts.Tags) == 0 || contains(accelerator.Status.Tags, opts.Tags) {

			values := []string{accelerator.Name}

			status := "unknown"
			for _, cond := range accelerator.Status.Conditions {
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
			tags := fmt.Sprintf("%v", accelerator.Status.Tags)
			if accelerator.Spec.Git != nil {
				repo = accelerator.Spec.Git.URL
				if accelerator.Spec.Git.Reference.Tag != "" {
					repo = repo + ":" + accelerator.Spec.Git.Reference.Tag
				} else if accelerator.Spec.Git.Reference.Branch != "" {
					repo = repo + ":" + accelerator.Spec.Git.Reference.Branch
				}
				if accelerator.Spec.Git.SubPath != nil {
					repo = repo + ":/" + *accelerator.Spec.Git.SubPath
				}
				values = append(values, tags, repo, status)
			} else if accelerator.Spec.Source != nil {
				repo = "source-image: " + accelerator.Spec.Source.Image
				values = append(values, tags, repo, status)
			} else {
				values = append(values, "", "", "")
			}
			accList = append(accList, values)
		}
	}

	printAcceleratorList(c, cmd, w, accList)
	return nil
}

func contains(accTags []string, input []string) bool {

	for _, inputTag := range input {
		result := false
		for _, accTag := range accTags {
			if strings.TrimSpace(inputTag) == accTag {
				result = true
			}
		}
		if result == false {
			return false
		}
	}

	return true
}

func printAcceleratorList(c *cli.Config, cmd *cobra.Command, w *tabwriter.Writer, accelerators [][]string) {
	if len(accelerators) == 0 {
		c.Infof("No accelerators found.\n")
	} else {
		fmt.Fprintln(w, "NAME\tTAGS\tREADY\tREPOSITORY")
		for _, values := range accelerators {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", values[0], values[1], values[3], values[2])
		}
		w.Flush()
	}
}
