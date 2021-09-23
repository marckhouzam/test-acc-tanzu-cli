## tanzu accelerator update

Update an accelerator

### Synopsis

Udate an accelerator resource with the specified name using the specified configuration.

Accelerator configuration options include:
- Git repository URL and branch/tag where accelerator code and metadata is defined
- Metadata like description, display-name, tags and icon-url

The update command also provides a --reoncile flag that will force the accelerator to be refreshed
with any changes made to the associated Git repository.


```
tanzu accelerator update [flags]
```

### Examples

```
tanzu accelerator update <accelerator-name> --description "Lorem Ipsum"
```

### Options

```
      --description string      description of this accelerator
      --display-name string     display name for the accelerator
      --git-branch string       Git repository branch to be used (default "main")
      --git-repository string   Git repository URL for the accelerator
      --git-tag string          Git repository tag to be used
  -h, --help                    help for update
      --icon-url string         URL for icon to use with the accelerator
  -n, --namespace name          kubernetes namespace (defaulted from kube config)
      --reconcile               trigger a reconciliation including the associated GitRepository resource
      --tags strings            tags that can be used to search for accelerators
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



