## tanzu accelerator list

List accelerators

### Synopsis

List all accelerators.

You can choose to list the accelerators from the Application Accelerator server using --server-url flag
or from a Kubernetes context using --from-context flag. The default is to list accelerators from the
Application Acceleratior server and you can set the ACC_SERVER_URL environment variable with the URL for
the Application Acceleratior server you want to access.


```
tanzu accelerator list [flags]
```

### Examples

```
tanzu accelerator list
```

### Options

```
      --from-context        retrieve resources from current context defined in kubeconfig
  -h, --help                help for list
  -n, --namespace name      kubernetes namespace (defaulted from kube config)
      --server-url string   the URL for the Application Accelerator server (default "http://localhost:8877")
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



