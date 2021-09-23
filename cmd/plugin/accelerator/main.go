/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package main

import (
	"context"
	"fmt"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-tanzu-cli/pkg/commands"
	tanzucliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/command/plugin"
	"k8s.io/apimachinery/pkg/runtime"

	// load credential helpers
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	_ "k8s.io/client-go/plugin/pkg/client/auth/openstack"

	cli "github.com/vmware-tanzu-private/tanzu-cli-apps-plugins/pkg/cli-runtime"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = acceleratorv1alpha1.AddToScheme(scheme)
}

func main() {
	ctx := context.Background()
	p, err := plugin.NewPlugin(&tanzucliv1alpha1.PluginDescriptor{
		Name:    "accelerator",
		Version: "v0.3.0-dev",
		Description: `Manage accelerators in a Kubernetes cluster.

Accelerators contain complete and runnable application code and/or deployment configurations.
The accelerator also contains metadata for altering the code and deployment configurations
based on input values provided for specific options that are defined in the accelerator metadata.

Operators would typically use create, update and delete commands for managing accelerators in a
Kubernetes context. Developers would use the list, get and generate commands for using accelerators
available in an Application Accelerator server. When operators want to use get and list commands
they can specify the --from-context flag to access accelerators in a Kubernetes context.

`,
		Group:          tanzucliv1alpha1.BuildCmdGroup,
		CompletionType: tanzucliv1alpha1.NativePluginCompletion,
		Aliases:        []string{"acc"},
	})
	if err != nil {
		panic(err)
	}

	c := cli.Initialize(fmt.Sprintf("tanzu %s", p.Cmd.Use), scheme)
	p.Cmd.CompletionOptions.DisableDefaultCmd = true

	if err != nil {
		panic(err)
	}
	p.AddCommands(
		commands.CreateCmd(ctx, c),
		commands.DeleteCmd(ctx, c),
		commands.ListCmd(ctx, c),
		commands.GetCmd(ctx, c),
		commands.UpdateCmd(ctx, c),
		commands.DocsCommand(ctx, c),
		commands.GenerateCmd(),
	)

	p.Cmd.PersistentFlags().StringVar(&c.KubeConfigFile, "kubeconfig", "", "kubeconfig `file` (default is $HOME/.kube/config)")
	p.Cmd.PersistentFlags().StringVar(&c.CurrentContext, "context", "", "`name` of the kubeconfig context to use (default is current-context defined by kubeconfig)")

	if err := p.Execute(); err != nil {
		panic(err)
	}

}
