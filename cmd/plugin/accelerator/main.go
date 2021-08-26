/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package main

import (
	"github/vmware-tanzu-private/tanzu-cli-app-accelerator/pkg/commands"

	acceleratorClientSet "github.com/pivotal/acc-controller/api/clientset"
	tanzucliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/command/plugin"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	// load credential helpers
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	_ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
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

	p.Cmd.CompletionOptions.DisableDefaultCmd = true

	kubeconfig := homedir.HomeDir() + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	defaultUiServerUrl := commands.EnvVar("ACC_UI_SERVER_URL", "http://acc-ui-server.accelerator-system")

	clientset, err := acceleratorClientSet.NewForConfig(config)

	if err != nil {
		panic(err)
	}
	p.AddCommands(
		commands.CreateCmd(clientset),
		commands.DeleteCmd(clientset),
		commands.ListCmd(clientset),
		commands.GetCmd(clientset),
		commands.UpdateCmd(clientset),
		commands.RunCmd(defaultUiServerUrl),
	)

	if err := p.Execute(); err != nil {
		panic(err)
	}

}
