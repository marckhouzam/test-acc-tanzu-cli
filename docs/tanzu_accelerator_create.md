## tanzu accelerator create

Create a new accelerator

### Synopsis

Create a new accelerator resource using the provided options

```
tanzu accelerator create [flags]
```

### Examples

```
tanzu accelerator create <accelerator-name> -git-repository <git-repo-URL>
```

### Options

```
      --description string      Accelerator description
      --display-name string     Accelerator display name
      --git-branch string       Accelerator repo branch
      --git-repository string   Accelerator repo URL
      --git-tag string          Accelerator repo tag
  -h, --help                    help for create
      --icon-url string         Accelerator icon location
  -n, --namespace name          kubernetes namespace (defaulted from kube config)
      --tags strings            Accelerator Tags
```

### Options inherited from parent commands

```
      --context name      name of the kubeconfig context to use (default is current-context defined by kubeconfig)
      --kubeconfig file   kubeconfig file (default is $HOME/.kube/config)
```

### SEE ALSO

* [tanzu accelerator](tanzu_accelerator.md)	 - Manage accelerators in your kubernetes cluster

