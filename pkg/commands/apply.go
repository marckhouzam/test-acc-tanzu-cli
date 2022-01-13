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
	cli "github.com/vmware-tanzu/tanzu-cli-apps-plugins/pkg/cli-runtime"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func loadAcceleratorFromFile(file string, a *acceleratorv1alpha1.Accelerator) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	d := yaml.NewYAMLOrJSONDecoder(f, 4096)
	documents := 0
	for {
		accelerator := &acceleratorv1alpha1.Accelerator{}
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
			return fmt.Errorf("files containing multiple accelerators are not supported")
		}
		accelerator.DeepCopyInto(a)
		documents++
	}
	return nil
}

func ApplyCmd(ctx context.Context, c *cli.Config) *cobra.Command {
	opts := ApplyOptions{}
	cmd := &cobra.Command{
		Use:     "apply",
		Short:   "Apply accelerator",
		Long:    "Create or update accelerator resource using specified accelerator manifest file.",
		Example: "tanzu accelerator apply --filename <path-to-accelerator-manifest>",
		RunE: func(cmd *cobra.Command, args []string) error {
			fileAcc := acceleratorv1alpha1.Accelerator{}
			err := loadAcceleratorFromFile(opts.FileName, &fileAcc)
			acceleratorName := fileAcc.ObjectMeta.Name
			acceleratorNamespace := fileAcc.ObjectMeta.Namespace
			if acceleratorNamespace == "" {
				acceleratorNamespace = "accelerator-system"
				fileAcc.ObjectMeta.Namespace = "accelerator-system"
			}
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error loading file %s\n", opts.FileName)
				return err
			}
			currentAcc := &acceleratorv1alpha1.Accelerator{}
			err = c.Get(ctx, client.ObjectKey{Namespace: acceleratorNamespace, Name: acceleratorName}, currentAcc)
			if err != nil && !errors.IsNotFound(err) {
				fmt.Fprintf(cmd.OutOrStderr(), "Error getting accelerator %s\n", args[0])
				return err
			} else if err != nil && errors.IsNotFound(err) {
				err = c.Create(ctx, &fileAcc)
				if err != nil {
					fmt.Fprintf(cmd.OutOrStderr(), "Error creating accelerator %s\n", args[0])
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "created accelerator %s in namespace %s\n", fileAcc.Name, fileAcc.Namespace)
				return nil
			} else {
				mergo.Merge(currentAcc, fileAcc, mergo.WithOverride)
				c.Update(ctx, currentAcc)
				fmt.Fprintf(cmd.OutOrStdout(), "updated accelerator %s in namespace %s\n", currentAcc.Name, currentAcc.Namespace)
				return nil
			}
		},
	}
	opts.DefineFlags(ctx, cmd, c)
	return cmd
}
