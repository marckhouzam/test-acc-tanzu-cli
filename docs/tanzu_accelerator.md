## tanzu accelerator

Manage accelerators in a Kubernetes cluster.

Accelerators contain complete and runnable application code and/or deployment configurations.
The accelerator also contains metadata for altering the code and deployment configurations
based on input values provided for specific options that are defined in the accelerator metadata.

Operators would typically use create, update and delete commands for managing accelerators in a
Kubernetes context. Developers would use the list, get and generate commands for using accelerators
available in an Application Accelerator server. When operators want to use get and list commands
they can specify the --from-context flag to access accelerators in a Kubernetes context.



### Options

```
      --context name      name of the kubeconfig context to use (default is current-context defined by kubeconfig)
  -h, --help              help for accelerator
      --kubeconfig file   kubeconfig file (default is $HOME/.kube/config)
```

### SEE ALSO

* [tanzu accelerator create](tanzu_accelerator_create.md)	 - Create a new accelerator
* [tanzu accelerator delete](tanzu_accelerator_delete.md)	 - Delete an accelerator
* [tanzu accelerator generate](tanzu_accelerator_generate.md)	 - Generate project from accelerator
* [tanzu accelerator get](tanzu_accelerator_get.md)	 - Get accelerator info
* [tanzu accelerator list](tanzu_accelerator_list.md)	 - List accelerators
* [tanzu accelerator update](tanzu_accelerator_update.md)	 - Update an accelerator

