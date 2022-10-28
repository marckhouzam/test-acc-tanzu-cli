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
	"github.com/imdario/mergo"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	fluxcdv1beta1 "github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	"github.com/pivotal/acc-controller/sourcecontroller/api/v1alpha1"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func FragmentUpdateCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := FragmentUpdateOptions{}
	requestedAtAnnotation := "reconcile.accelerator.apps.tanzu.vmware.com/requestedAt"
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an accelerator fragment",
		Long: `Update an accelerator fragment resource with the specified name using the specified configuration.

Accelerator configuration options include:
- Git repository URL and branch/tag where accelerator code and metadata is defined
- Metadata like display-name

The update command also provides a --reoncile flag that will force the accelerator fragment to be refreshed
with any changes made to the associated Git repository.
`,
		ValidArgsFunction: SuggestFragmentNamesFromConfig(context.Background(), c),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must specify the name of the accelerator fragment")
			}
			return nil
		},
		Example: "tanzu accelerator update <accelerator-name> --description \"Lorem Ipsum\"",
		RunE: func(cmd *cobra.Command, args []string) error {

			fragment := &acceleratorv1alpha1.Fragment{}
			err := c.Get(context.Background(), client.ObjectKey{Namespace: opts.Namespace, Name: args[0]}, fragment)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "accelerator fragment %s not found\n", args[0])
				return err
			}
			updatedFragment := &acceleratorv1alpha1.Fragment{
				TypeMeta: v1.TypeMeta{
					APIVersion: "accelerator.tanzu.vmware.com/v1alpha1",
					Kind:       "Fragment",
				},
				ObjectMeta: v1.ObjectMeta{
					Namespace: opts.Namespace,
				},
				Spec: acceleratorv1alpha1.FragmentSpec{
					DisplayName: opts.DisplayName,
				},
			}

			if opts.GitRepoUrl != "" {
				updatedFragment.Spec.Git = &acceleratorv1alpha1.Git{
					URL: opts.GitRepoUrl,
					Reference: &fluxcdv1beta1.GitRepositoryRef{
						Branch: opts.GitBranch,
						Tag:    opts.GitTag,
					},
				}
				if opts.GitSubPath != "" {
					updatedFragment.Spec.Git.SubPath = &opts.GitSubPath
				}
			}

			if opts.SourceImage != "" {
				updatedFragment.Spec.Source = &v1alpha1.ImageRepositorySpec{
					Image: opts.SourceImage,
				}
				updatedFragment.Spec.Git = nil
			}

			if opts.Reconcile {
				if fragment.ObjectMeta.Annotations == nil {
					fragment.ObjectMeta.Annotations = make(map[string]string)
				}
				fragment.ObjectMeta.Annotations[requestedAtAnnotation] = time.Now().UTC().Format(time.RFC3339)
			}

			mergo.Merge(updatedFragment, *fragment)

			if opts.GitRepoUrl == "" && opts.GitBranch != "" {
				updatedFragment.Spec.Git.Reference.Branch = opts.GitBranch
			}

			if opts.GitRepoUrl == "" && opts.GitTag != "" {
				updatedFragment.Spec.Git.Reference.Tag = opts.GitTag
			}

			if opts.GitRepoUrl == "" && opts.GitSubPath != "" {
				updatedFragment.Spec.Git.SubPath = &opts.GitSubPath
			}

			if opts.Interval != "" {
				duration, _ := time.ParseDuration(opts.Interval)
				interval := v1.Duration{
					Duration: duration,
				}

				if updatedFragment.Spec.Git != nil {
					updatedFragment.Spec.Git.Interval = &interval
				}
			}

			if opts.SecretRef != "" {
				if updatedFragment.Spec.Git != nil {
					ref := meta.LocalObjectReference{
						Name: opts.SecretRef,
					}
					updatedFragment.Spec.Git.SecretRef = &ref
				}
			}

			err = c.Update(ctx, updatedFragment)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "there was an error updating accelerator fragment %s\n", args[0])
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "accelerator fragment %s updated successfully\n", args[0])
			return nil
		},
	}
	opts.DefineFlags(ctx, updateCmd, c)
	return updateCmd
}
