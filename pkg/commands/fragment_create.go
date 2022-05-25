/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fluxcd/pkg/apis/meta"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	fluxcdv1beta1 "github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	"github.com/spf13/cobra"
	cli "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func FragmentCreateCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := FragmentCreateOptions{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new accelerator fragment",
		Long: `Create a new accelerator fragment resource with specified configuration.

Accelerator configuration options include:
- Git repository URL and branch/tag where accelerator code and metadata is defined
- Metadata like description, display-name, tags and icon-url

The Git repository option is required. Metadata options are optional and will override any values for
the same options specified in the accelerator metadata retrieved from the Git repository.
`,
		Example: "tanzu acceleratorent fragm create <fragment-name> --git-repository <URL> --git-branch <branch> --git-sub-path <sub-path>",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must specify the name of the accelerator fragment")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			frag := &acceleratorv1alpha1.Fragment{
				TypeMeta: v1.TypeMeta{
					APIVersion: "accelerator.tanzu.vmware.com/v1alpha1",
					Kind:       "Fragment",
				},
				ObjectMeta: v1.ObjectMeta{
					Namespace: opts.Namespace,
					Name:      args[0],
				},
				Spec: acceleratorv1alpha1.FragmentSpec{
					DisplayName: opts.DisplayName,
				},
			}

			if opts.GitRepoUrl != "" {
				frag.Spec.Git = &acceleratorv1alpha1.Git{
					URL: opts.GitRepoUrl,
					Reference: &fluxcdv1beta1.GitRepositoryRef{
						Branch: opts.GitBranch,
						Tag:    opts.GitTag,
					},
				}
				if opts.GitSubPath != "" {
					frag.Spec.Git.SubPath = &opts.GitSubPath
				}
			}

			if opts.Interval != "" {
				duration, _ := time.ParseDuration(opts.Interval)
				interval := v1.Duration{
					Duration: duration,
				}

				if frag.Spec.Git != nil {
					frag.Spec.Git.Interval = &interval
				}
			}

			if opts.SecretRef != "" {
				ref := meta.LocalObjectReference{
					Name: opts.SecretRef,
				}
				frag.Spec.Git.SecretRef = &ref
			}

			err := c.Create(ctx, frag)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error creating accelerator fragment %s\n", args[0])
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "created accelerator fragment %s in namespace %s\n", frag.Name, frag.Namespace)
			return nil

		},
	}
	opts.DefineFlags(ctx, cmd, c)
	return cmd
}
