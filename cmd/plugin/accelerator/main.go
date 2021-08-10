/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package main

import (
	"log"

	"github/vmware-tanzu-private/tanzu-cli-app-accelerator/pkg/commands"

	acceleratorClientSet "github.com/pivotal/acc-controller/api/clientset"
	tanzucliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/command/plugin"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	p, err := plugin.NewPlugin(&tanzucliv1alpha1.PluginDescriptor{
		Name:           "accelerator",
		Version:        "v0.2.0",
		Description:    "Manage accelerators in your kubernetes cluster",
		Group:          tanzucliv1alpha1.BuildCmdGroup,
		CompletionType: tanzucliv1alpha1.NativePluginCompletion,
		Aliases:        []string{"acc"},
	})
	if err != nil {
		log.Fatal(err)
	}

	p.Cmd.CompletionOptions.DisableDefaultCmd = true

	kubeconfig := homedir.HomeDir() + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := acceleratorClientSet.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	p.AddCommands(
		commands.CreateCmd(clientset),
		commands.DeleteCmd(clientset),
		commands.ListCmd(clientset),
		commands.GetCmd(clientset),
		commands.UpdateCmd(clientset),
	)

	if err := p.Execute(); err != nil {
		log.Fatal(err)
	}

}
