## tanzu accelerator get

Get accelerator info

### Synopsis

Get accelerator info.

You can choose to get the accelerator from the Application Accelerator server using --server-url flag
or from a Kubernetes context using --from-context flag. The default is to get accelerators from the
Application Acceleratior server and you can set the ACC_SERVER_URL environment variable with the URL for
the Application Acceleratior server you want to access.


```
tanzu accelerator get [flags]
```

### Examples

```
tanzu accelerator get <accelerator-name> --from-context
```

### Options

```
      --from-context        retrieve resources from current context defined in kubeconfig
  -h, --help                help for get
  -n, --namespace name      kubernetes namespace (defaulted from kube config)
      --server-url string   the URL for the Application Accelerator server (default "http://localhost:8877")
```

### Options inherited from parent commands

```
      --context name      name of the kubeconfig context to use (default is current-context defined by kubeconfig)
      --kubeconfig file   kubeconfig file (default is $HOME/.kube/config)
```

### SEE ALSO

* [tanzu accelerator](tanzu_accelerator.md)	 - Manage accelerators in a Kubernetes cluster

