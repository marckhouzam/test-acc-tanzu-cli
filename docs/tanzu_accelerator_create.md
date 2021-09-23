## tanzu accelerator create

Create a new accelerator

### Synopsis

Create a new accelerator resource with specified configuration.

Accelerator configuration options include:
- Git repository URL and branch/tag where accelerator code and metadata is defined
- Metadata like description, display-name, tags and icon-url

The Git repository option is required. Metadata options are optional and will override any values for
the same options specified in the accelerator metadata retrieved from the Git repository.


```
tanzu accelerator create [flags]
```

### Examples

```
tanzu accelerator create <accelerator-name> --git-repository <URL> --git-branch <branch>
```

### Options

```
      --description string      description of this accelerator
      --display-name string     display name for the accelerator
      --git-branch string       Git repository branch to be used
      --git-repository string   Git repository URL for the accelerator
      --git-tag string          Git repository tag to be used
  -h, --help                    help for create
      --icon-url string         URL for icon to use with the accelerator
  -n, --namespace name          kubernetes namespace (defaulted from kube config)
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



