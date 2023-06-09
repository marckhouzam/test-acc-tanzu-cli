/*
Copyright 2021-2023 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fluxcd/pkg/apis/meta"
	"github.com/imdario/mergo"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	fluxcdv1beta1 "github.com/pivotal/acc-controller/fluxcd/api/v1beta2"
	"github.com/pivotal/acc-controller/sourcecontroller/api/v1alpha1"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func UpdateCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := UpdateOptions{}
	requestedAtAnnotation := "reconcile.accelerator.apps.tanzu.vmware.com/requestedAt"
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an accelerator",
		Long: `Update an accelerator resource with the specified name using the specified configuration.

Accelerator configuration options include:
- Git repository URL and branch/tag where accelerator code and metadata is defined
- Metadata like description, display-name, tags and icon-url

The update command also provides a --reoncile flag that will force the accelerator to be refreshed
with any changes made to the associated Git repository.
`,
		ValidArgsFunction: SuggestAcceleratorNamesFromConfig(context.Background(), c),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must specify the name of the accelerator")
			}
			return nil
		},
		Example: "tanzu accelerator update <accelerator-name> --description \"Lorem Ipsum\"",
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.GitRepoUrl != "" && opts.SourceImage != "" {
				return errors.New("you may only provide one of --git-repository or --source-image")
			}

			accelerator := &acceleratorv1alpha1.Accelerator{}
			err := c.Get(context.Background(), client.ObjectKey{Namespace: opts.Namespace, Name: args[0]}, accelerator)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "accelerator %s not found\n", args[0])
				return err
			}
			updatedAccelerator := &acceleratorv1alpha1.Accelerator{
				TypeMeta: v1.TypeMeta{
					APIVersion: "accelerator.tanzu.vmware.com/v1alpha1",
					Kind:       "Accelerator",
				},
				ObjectMeta: v1.ObjectMeta{
					Namespace: opts.Namespace,
				},
				Spec: acceleratorv1alpha1.AcceleratorSpec{
					DisplayName: opts.DisplayName,
					Description: opts.Description,
					IconUrl:     opts.IconUrl,
					Tags:        opts.Tags,
				},
			}

			if opts.GitRepoUrl != "" {
				updatedAccelerator.Spec.Git = &acceleratorv1alpha1.Git{
					URL: opts.GitRepoUrl,
					Reference: &fluxcdv1beta1.GitRepositoryRef{
						Branch: opts.GitBranch,
						Tag:    opts.GitTag,
					},
				}
				if opts.GitSubPath != "" {
					updatedAccelerator.Spec.Git.SubPath = &opts.GitSubPath
				}
				accelerator.Spec.Source = nil
			}

			if opts.SourceImage != "" {
				updatedAccelerator.Spec.Source = &v1alpha1.ImageRepositorySpec{
					Image: opts.SourceImage,
				}
				accelerator.Spec.Git = nil
			}

			if opts.Reconcile {
				if accelerator.ObjectMeta.Annotations == nil {
					accelerator.ObjectMeta.Annotations = make(map[string]string)
				}
				accelerator.ObjectMeta.Annotations[requestedAtAnnotation] = time.Now().UTC().Format(time.RFC3339)
			}

			mergo.Merge(updatedAccelerator, *accelerator)

			if opts.GitRepoUrl == "" && opts.GitBranch != "" {
				updatedAccelerator.Spec.Git.Reference.Branch = opts.GitBranch
			}

			if opts.GitRepoUrl == "" && opts.GitTag != "" {
				updatedAccelerator.Spec.Git.Reference.Tag = opts.GitTag
			}

			if opts.GitRepoUrl == "" && opts.GitSubPath != "" {
				updatedAccelerator.Spec.Git.SubPath = &opts.GitSubPath
			}

			if opts.Interval != "" {
				duration, _ := time.ParseDuration(opts.Interval)
				interval := v1.Duration{
					Duration: duration,
				}

				if updatedAccelerator.Spec.Source != nil {
					updatedAccelerator.Spec.Source.Interval = &interval
				}

				if updatedAccelerator.Spec.Git != nil {
					updatedAccelerator.Spec.Git.Interval = &interval
				}
			}

			if opts.SecretRef != "" {
				if updatedAccelerator.Spec.Source != nil {
					ref := corev1.LocalObjectReference{
						Name: opts.SecretRef,
					}
					updatedAccelerator.Spec.Source.ImagePullSecrets = []corev1.LocalObjectReference{ref}
				}

				if updatedAccelerator.Spec.Git != nil {
					ref := meta.LocalObjectReference{
						Name: opts.SecretRef,
					}
					updatedAccelerator.Spec.Git.SecretRef = &ref
				}
			}

			err = c.Update(ctx, updatedAccelerator)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "there was an error updating accelerator %s\n", args[0])
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "accelerator %s updated successfully\n", args[0])
			return nil
		},
	}
	opts.DefineFlags(ctx, updateCmd, c)
	return updateCmd
}
