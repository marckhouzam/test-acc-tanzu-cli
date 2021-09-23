## tanzu accelerator update

Update an accelerator

### Synopsis

Update an accelerator resource using the provided options

```
tanzu accelerator update [flags]
```

### Examples

```
tanzu accelerator update <accelerator-name> --description "Lorem Ipsum"
```

### Options

```
      --description string      Accelerator description
      --display-name string     Accelerator display name
      --git-branch string       Accelerator repo branch (default "main")
      --git-repository string   Accelerator repo URL
  -h, --help                    help for update
      --icon-url string         Accelerator icon location
  -n, --namespace name          kubernetes namespace (defaulted from kube config)
      --reconcile               Trigger a reconciliation including the associated GitRepository resource
      --tags strings            Accelerator Tags
```

### Options inherited from parent commands

```
      --context name      name of the kubeconfig context to use (default is current-context defined by kubeconfig)
      --kubeconfig file   kubeconfig file (default is $HOME/.kube/config)
```

### SEE ALSO

* [tanzu accelerator](tanzu_accelerator.md)	 - Manage accelerators in your kubernetes cluster

