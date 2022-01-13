## tanzu accelerator apply

Apply accelerator

### Synopsis

Create or update accelerator resource using specified accelerator manifest file.

```
tanzu accelerator apply [flags]
```

### Examples

```
tanzu accelerator apply --filename <path-to-accelerator-manifest>
```

### Options

```
  -f, --filename string   path of manifest file for the accelerator
  -h, --help              help for apply
```

### Options inherited from parent commands

```
      --context name      name of the kubeconfig context to use (default is current-context defined by kubeconfig)
      --kubeconfig file   kubeconfig file (default is $HOME/.kube/config)
```

### SEE ALSO

* [tanzu accelerator](tanzu_accelerator.md)	 - Manage accelerators in a Kubernetes cluster

