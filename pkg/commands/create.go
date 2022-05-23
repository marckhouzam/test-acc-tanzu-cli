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
	"github.com/pivotal/acc-controller/sourcecontroller/api/v1alpha1"
	"github.com/spf13/cobra"
	cli "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := CreateOptions{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new accelerator",
		Long: `Create a new accelerator resource with specified configuration.

Accelerator configuration options include:
- Git repository URL and branch/tag where accelerator code and metadata is defined
- Metadata like description, display-name, tags and icon-url

The Git repository option is required. Metadata options are optional and will override any values for
the same options specified in the accelerator metadata retrieved from the Git repository.
`,
		Example: "tanzu accelerator create <accelerator-name> --git-repository <URL> --git-branch <branch>",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must specify the name of the accelerator")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			if opts.GitRepoUrl == "" && opts.SourceImage == "" && opts.LocalPath == "" {
				return errors.New("you must provide --git-repository or --source-image")
			}

			if opts.GitRepoUrl != "" && opts.SourceImage != "" {
				return errors.New("you may only provide one of --git-repository or --source-image")
			}

			if opts.LocalPath != "" && opts.SourceImage == "" {
				return errors.New("you must provide --source-image when using --local-path")
			}

			acc := &acceleratorv1alpha1.Accelerator{
				TypeMeta: v1.TypeMeta{
					APIVersion: "accelerator.tanzu.vmware.com/v1alpha1",
					Kind:       "Accelerator",
				},
				ObjectMeta: v1.ObjectMeta{
					Namespace: opts.Namespace,
					Name:      args[0],
				},
				Spec: acceleratorv1alpha1.AcceleratorSpec{
					DisplayName: opts.DisplayName,
					Description: opts.Description,
					IconUrl:     opts.IconUrl,
					Tags:        opts.Tags,
				},
			}

			if opts.GitRepoUrl != "" {
				acc.Spec.Git = &acceleratorv1alpha1.Git{
					URL: opts.GitRepoUrl,
					Reference: &fluxcdv1beta1.GitRepositoryRef{
						Branch: opts.GitBranch,
						Tag:    opts.GitTag,
					},
				}
			}

			if opts.LocalPath != "" {
				if err := opts.PublishLocalSource(ctx, c); err != nil {
					return err
				}
			}

			if opts.SourceImage != "" {
				acc.Spec.Source = &v1alpha1.ImageRepositorySpec{
					Image: opts.SourceImage,
				}
			}

			if opts.Interval != "" {
				duration, _ := time.ParseDuration(opts.Interval)
				interval := v1.Duration{
					Duration: duration,
				}

				if acc.Spec.Source != nil {
					acc.Spec.Source.Interval = &interval
				}

				if acc.Spec.Git != nil {
					acc.Spec.Git.Interval = &interval
				}
			}

			if opts.SecretRef != "" {
				if opts.SourceImage != "" {
					ref := corev1.LocalObjectReference{
						Name: opts.SecretRef,
					}
					acc.Spec.Source.ImagePullSecrets = []corev1.LocalObjectReference{ref}
				} else {
					ref := meta.LocalObjectReference{
						Name: opts.SecretRef,
					}
					acc.Spec.Git.SecretRef = &ref
				}
			}

			err := c.Create(ctx, acc)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error creating accelerator %s\n", args[0])
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "created accelerator %s in namespace %s\n", acc.Name, acc.Namespace)
			return nil

		},
	}
	opts.DefineFlags(ctx, cmd, c)
	return cmd
}
