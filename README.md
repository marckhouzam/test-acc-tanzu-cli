# App Accelerator tanzu CLI plugin

![ci](https://github.com/pivotal/acc-tanzu-cli/actions/workflows/ci.yaml/badge.svg?branch=main)

Early prototype for Accelerator CLI commands.

To check the plugin documentation go to

## Commands

- [tanzu accelerator](./cmd/plugin/accelerator/README.md)

## Install

The [Tanzu CLI](https://docs.vmware.com/en/VMware-Tanzu-Application-Platform/1.3/tap/GUID-install-tanzu-cli.html#install-or-update-the-tanzu-cli-and-plugins-3) is required to use the Accelerator CLI plugin.

### From a pre-built distribution

Download `tanzu-accelerator-plugin-<version>.tar.gz` from the most recent release listed on the [App Accelerator tanzu CLI plugin releases page](https://github.com/pivotal/acc-tanzu-cli/releases).

Extract the archive to a local directory:

```sh
mkdir plugin_bundle
tar -zxvf tanzu-accelerator-plugin-*.tar.gz -C plugin_bundle
```

Install the accelerator plugin:

- For MacOS:
    ```sh
    tanzu plugin install accelerator -t kubernetes --local ./plugin_bundle/darwin/amd64
    ```
- For Linux:
    ```sh
    tanzu plugin install accelerator -t kubernetes --local ./plugin_bundle/linux/amd64
    ```
- For Windows:
    ```sh
    tanzu plugin install accelerator -t kubernetes --local ./plugin_bundle/windows/amd64
    ```

### Build from source

See the [Development Guide](./DEVELOPMENT.md).
