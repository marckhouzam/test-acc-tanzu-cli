## tanzu accelerator delete

Delete an accelerator

### Synopsis

Delete the accelerator resource with the specified name.

```
tanzu accelerator delete [flags]
```

### Examples

```
tanzu accelerator delete <accelerator-name>
```

### Options

```
  -h, --help             help for delete
  -n, --namespace name   kubernetes namespace (defaulted from kube config)
```

### Options inherited from parent commands

```
      --context name      name of the kubeconfig context to use (default is current-context defined by kubeconfig)
      --kubeconfig file   kubeconfig file (default is $HOME/.kube/config)
```

### SEE ALSO

* [tanzu accelerator](tanzu_accelerator.md)	 - Manage accelerators in a Kubernetes cluster.

Accelerators contain complete and runnable application code and/or deployment configurations.
The accelerator also contains metadata for altering the code and deployment configurations
based on input values provided for specific options that are defined in the accelerator metadata.

Operators would typically use create, update and delete commands for managing accelerators in a
Kubernetes context. Developers would use the list, get and generate commands for using accelerators
available in an Application Accelerator server. When operators want to use get and list commands
they can specify the --from-context flag to access accelerators in a Kubernetes context.



