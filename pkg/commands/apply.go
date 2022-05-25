/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/imdario/mergo"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/spf13/cobra"
	cli "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ApplyCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := ApplyOptions{}
	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Apply accelerator resource",
		Long:    "Create or update accelerator resource using specified manifest file.",
		Example: "tanzu accelerator apply --filename <path-to-resource-manifest>",
		RunE: func(cmd *cobra.Command, args []string) error {
			var fileObj runtime.RawExtension
			err := loadResourceFromFile(opts.FileName, &fileObj)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error loading file %s\n", opts.FileName)
				return err
			}
			obj, _, err := unstructured.UnstructuredJSONScheme.Decode(fileObj.Raw, nil, nil)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error decoding file %s\n", opts.FileName)
				return err
			}

			if obj.GetObjectKind().GroupVersionKind().Kind == "Accelerator" {
				providedResource := acceleratorv1alpha1.Accelerator{}
				_, _, err := unstructured.UnstructuredJSONScheme.Decode(fileObj.Raw, nil, &providedResource)
				if err != nil {
					fmt.Fprintf(cmd.OutOrStderr(), "Error decoding Accelerator resource from file %s\n", opts.FileName)
					return err
				}
				err = saveAcceleratorResource(ctx, c, providedResource, opts, cmd)
				if err != nil {
					return err
				}
			} else if obj.GetObjectKind().GroupVersionKind().Kind == "Fragment" {
				providedResource := acceleratorv1alpha1.Fragment{}
				_, _, err := unstructured.UnstructuredJSONScheme.Decode(fileObj.Raw, nil, &providedResource)
				if err != nil {
					fmt.Fprintf(cmd.OutOrStderr(), "Error decoding Fragment resource from file %s\n", opts.FileName)
					return err
				}
				err = saveFragmentResource(ctx, c, providedResource, opts, cmd)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("the resource kind \"%s\" in the provided file \"%s\" does not match \"Accelerator\" or \"Fragment\"", obj.GetObjectKind().GroupVersionKind().Kind, opts.FileName)
			}
			return nil
		},
	}
	opts.DefineFlags(ctx, cmd, c)
	return cmd
}

func loadResourceFromFile(file string, a *runtime.RawExtension) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	d := yamlutil.NewYAMLOrJSONDecoder(f, 4096)
	documents := 0
	for {
		accelerator := &runtime.RawExtension{}
		if err := d.Decode(&accelerator); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if accelerator == nil {
			continue
		}
		if documents > 0 {
			return fmt.Errorf("files containing multiple resources are not supported")
		}
		accelerator.DeepCopyInto(a)
		documents++
	}
	return nil
}

func saveAcceleratorResource(ctx context.Context, c *cli.Config, providedResource acceleratorv1alpha1.Accelerator, opts ApplyOptions, cmd *cobra.Command) error {
	acceleratorName := providedResource.ObjectMeta.Name
	if providedResource.ObjectMeta.Namespace > "" && providedResource.ObjectMeta.Namespace != opts.Namespace {
		return fmt.Errorf("the namespace specified in the provided file \"%s\" does not match the namespace \"%s\". You must pass '--namespace=%s' to perform this operation.", providedResource.ObjectMeta.Namespace, opts.Namespace, providedResource.ObjectMeta.Namespace)
	}
	acceleratorNamespace := opts.Namespace
	if providedResource.ObjectMeta.Namespace == "" {
		providedResource.ObjectMeta.Namespace = opts.Namespace
	}
	currentAcc := &acceleratorv1alpha1.Accelerator{}
	err := c.Get(ctx, client.ObjectKey{Namespace: acceleratorNamespace, Name: acceleratorName}, currentAcc)
	if err != nil && !errors.IsNotFound(err) {
		fmt.Fprintf(cmd.OutOrStderr(), "Error getting accelerator %s\n", providedResource.Name)
		return err
	} else if err != nil && errors.IsNotFound(err) {
		err = c.Create(ctx, &providedResource)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error creating accelerator %s\n", providedResource.Name)
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "created accelerator %s in namespace %s\n", providedResource.Name, providedResource.Namespace)
		return nil
	} else {
		mergo.Merge(currentAcc, providedResource, mergo.WithOverride)
		c.Update(ctx, currentAcc)
		fmt.Fprintf(cmd.OutOrStdout(), "updated accelerator %s in namespace %s\n", currentAcc.Name, currentAcc.Namespace)
		return nil
	}
}

func saveFragmentResource(ctx context.Context, c *cli.Config, providedResource acceleratorv1alpha1.Fragment, opts ApplyOptions, cmd *cobra.Command) error {
	fragmentName := providedResource.ObjectMeta.Name
	if providedResource.ObjectMeta.Namespace > "" && providedResource.ObjectMeta.Namespace != opts.Namespace {
		return fmt.Errorf("the namespace specified in the provided file \"%s\" does not match the namespace \"%s\". You must pass '--namespace=%s' to perform this operation.", providedResource.ObjectMeta.Namespace, opts.Namespace, providedResource.ObjectMeta.Namespace)
	}
	fragmentNamespace := opts.Namespace
	if providedResource.ObjectMeta.Namespace == "" {
		providedResource.ObjectMeta.Namespace = opts.Namespace
	}
	currentAcc := &acceleratorv1alpha1.Fragment{}
	err := c.Get(ctx, client.ObjectKey{Namespace: fragmentNamespace, Name: fragmentName}, currentAcc)
	if err != nil && !errors.IsNotFound(err) {
		fmt.Fprintf(cmd.OutOrStderr(), "Error getting accelerator fragment %s\n", providedResource.Name)
		return err
	} else if err != nil && errors.IsNotFound(err) {
		err = c.Create(ctx, &providedResource)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error creating accelerator fragment %s\n", providedResource.Name)
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "created accelerator fragment %s in namespace %s\n", providedResource.Name, providedResource.Namespace)
		return nil
	} else {
		mergo.Merge(currentAcc, providedResource, mergo.WithOverride)
		c.Update(ctx, currentAcc)
		fmt.Fprintf(cmd.OutOrStdout(), "updated accelerator fragment %s in namespace %s\n", currentAcc.Name, currentAcc.Namespace)
		return nil
	}
}
