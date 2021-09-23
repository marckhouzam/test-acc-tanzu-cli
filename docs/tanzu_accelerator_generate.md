## tanzu accelerator generate

Generate project from accelerator

### Synopsis

Generate a project from an accelerator using provided options and download project artifacts as a ZIP file.

Generation options are provided as a JSON string and should match the metadata options that are specified for the
accelerator used for the generation. The options can include "projectName" which defaults to the name of the accelerator.
This "projectName" will be used as the name of the generated ZIP file.

You can see the available options by using the "tanzu accelerator list <accelerator-name>" command.

Here is an example of an options JSON string that specifies the "projectName" and an "includeKubernetes" boolean flag:

    --options '{"projectName":"test", "includeKubernetes": true}'

You can also provide a file that specifies the JSON string using the --options-file flag.

The generate command needs access to the Application Accelerator server. You can specify the --server-url flag or set
an ACC_SERVER_URL environment variable. If you specify the --server-url flag it will override the ACC_SERVER_URL
environmnet variable if it is set.


```
tanzu accelerator generate [flags]
```

### Examples

```
tanzu accelerator generate <accelerator-name> --options '{"projectName":"test"}'
```

### Options

```
  -h, --help                  help for generate
      --options string        options JSON string
      --options-file string   path to file containing options JSON string
      --output-dir string     directory that the zip file will be written to
      --server-url string     the URL for the Application Accelerator server (default "http://localhost:8877")
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



