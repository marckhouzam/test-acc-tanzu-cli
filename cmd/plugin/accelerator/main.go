/*
Copyright 2021-2023 VMware, Inc. All Rights Reserved.
*/
package main

import (
	"context"
	"fmt"
	"os"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-tanzu-cli/pkg/commands"
	"github.com/vmware-tanzu/tanzu-plugin-runtime/config/types"
	"github.com/vmware-tanzu/tanzu-plugin-runtime/plugin"
	"github.com/vmware-tanzu/tanzu-plugin-runtime/plugin/buildinfo"
	"k8s.io/apimachinery/pkg/runtime"

	// load credential helpers
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	cli "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = acceleratorv1alpha1.AddToScheme(scheme)
}

func main() {
	ctx := context.Background()
	p, err := plugin.NewPlugin(&plugin.PluginDescriptor{
		Name:           "accelerator",
		Version:        buildinfo.Version,
		Description:    "Manage accelerators in a Kubernetes cluster",
		Target:         types.TargetK8s,
		BuildSHA:       buildinfo.SHA,
		Group:          plugin.RunCmdGroup,
		CompletionType: plugin.NativePluginCompletion,
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
		commands.PushCmd(ctx, c),
		commands.ApplyCmd(ctx, c),
		commands.FragmentCmd(ctx, c),
		commands.LocalGenerateCmd(),
	)

	p.Cmd.PersistentFlags().StringVar(&c.KubeConfigFile, "kubeconfig", "", "kubeconfig `file` (default is $HOME/.kube/config)")
	p.Cmd.PersistentFlags().StringVar(&c.CurrentContext, "context", "", "`name` of the kubeconfig context to use (default is current-context defined by kubeconfig)")

	if err := p.Execute(); err != nil {
		println(err.Error())
		os.Exit(1)
	}

}
