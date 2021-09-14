/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package main

import (
	"context"
	"fmt"
	"github/vmware-tanzu-private/tanzu-cli-app-accelerator/pkg/commands"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
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
		Name:           "accelerator",
		Version:        "v0.3.0-dev",
		Description:    "Manage accelerators in your kubernetes cluster",
		Group:          tanzucliv1alpha1.BuildCmdGroup,
		CompletionType: tanzucliv1alpha1.NativePluginCompletion,
		Aliases:        []string{"acc"},
	})
	if err != nil {
		panic(err)
	}

	c := cli.Initialize(fmt.Sprintf("tanzu %s", p.Cmd.Use), scheme)
	p.Cmd.CompletionOptions.DisableDefaultCmd = true
	defaultUiServerUrl := commands.EnvVar("ACC_SERVER_URL", "http://localhost:8877")

	if err != nil {
		panic(err)
	}
	p.AddCommands(
		commands.CreateCmd(ctx, c),
		commands.DeleteCmd(ctx, c),
		commands.ListCmd(ctx, c),
		commands.GetCmd(ctx, c),
		commands.UpdateCmd(ctx, c),
		commands.GenerateCmd(defaultUiServerUrl),
	)

	p.Cmd.PersistentFlags().StringVar(&c.KubeConfigFile, "kubeconfig", "", "kubeconfig `file` (default is $HOME/.kube/config)")
	p.Cmd.PersistentFlags().StringVar(&c.CurrentContext, "context", "", "`name` of the kubeconfig context to use (default is current-context defined by kubeconfig)")

	if err := p.Execute(); err != nil {
		panic(err)
	}

}
