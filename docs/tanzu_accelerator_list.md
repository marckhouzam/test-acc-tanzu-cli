## tanzu accelerator list

List accelerators

### Synopsis

List all accelerators.

You can choose to list accelerators from a server using --server-url flag 
or from a Kubernetes context using --from-context flag.

```
tanzu accelerator list [flags]
```

### Examples

```
tanzu accelerator list
```

### Options

```
      --from-context        Retrieve resources from current context defined in kubeconfig
  -h, --help                help for list
  -n, --namespace name      kubernetes namespace (defaulted from kube config)
      --server-url string   Accelerator server URL to use for retrieving resources
```

### Options inherited from parent commands

```
      --context name      name of the kubeconfig context to use (default is current-context defined by kubeconfig)
      --kubeconfig file   kubeconfig file (default is $HOME/.kube/config)
```

### SEE ALSO

* [tanzu accelerator](tanzu_accelerator.md)	 - Manage accelerators in your kubernetes cluster

