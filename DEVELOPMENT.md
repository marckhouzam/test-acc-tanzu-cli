# App Accelerator tanzu CLI plugin - Development Guide

Before installing the plugin in your tanzu instance, you need to follow the steps described for installing [Tanzu CLI in the TAP Installation docs](https://docs.vmware.com/en/VMware-Tanzu-Application-Platform/0.1/tap-0-1/GUID-install.html#install-the-tanzu-cli-and-package-plugin-4).

Then install the Tanzu builder and test plugins:

```sh
tanzu plugin repo add -b tanzu-cli-admin-plugins -n admin -p artifacts-admin
```

```sh
tanzu plugin install builder
```

```sh
tanzu plugin install test
```

# Build the plugin

to build the plugin you need to:

 - Run the command `go mod vendor`
 - Run the `make build` task
 - The build task will create the `artifacts/` directory
 - To install the plugin to tanzu, you need to run `tanzu plugin install accelerator --local ./artifacts`

# Troubleshooting
If you get the error `fatal: could not read Username for 'https://github.com': terminal prompts disabled` while downloading the go dependencies: 

- run the following command:

    ```
    git config --global url."git@github.com:".insteadOf https://github.com/
    ```

- also make sure `GOPRIVATE` is set for 

    - github.com/pivotal
    - github.com/vmware-tanzu/*
    - github.com/vmware-tanzu-private/*`

    check with:

    ```
    go env GOPRIVATE
    ```

    if it is not set add it with something like the following (adjust if you aready have something in there):

    ```
    go env -w GOPRIVATE='github.com/pivotal,github.com/vmware-tanzu/*,github.com/vmware-tanzu-private/*'
    ```
