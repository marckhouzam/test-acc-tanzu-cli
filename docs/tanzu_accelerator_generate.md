## tanzu accelerator generate

Generate project from accelerator

### Synopsis

Generate a project from an accelerator and download project artifacts as a ZIP file

```
tanzu accelerator generate [flags]
```

### Options

```
  -h, --help                  help for generate
      --options string        Enter options string
      --options-file string   Enter file path with json body
      --output-dir string     Directory where the zip file should be written
      --server-url string     The App Accelerator server URL, this will override ACC_SERVER_URL env variable (default "http://localhost:8877")
```

### Options inherited from parent commands

```
      --context name      name of the kubeconfig context to use (default is current-context defined by kubeconfig)
      --kubeconfig file   kubeconfig file (default is $HOME/.kube/config)
```

### SEE ALSO

* [tanzu accelerator](tanzu_accelerator.md)	 - Manage accelerators in your kubernetes cluster

