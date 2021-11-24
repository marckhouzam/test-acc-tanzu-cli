# App Accelerator tanzu CLI plugin

![ci](https://github.com/pivotal/acc-tanzu-cli/actions/workflows/ci.yaml/badge.svg?branch=main)

Early prototype for Accelerator CLI commands.

To check the plugin documentation go to 

## Commands

- [tanzu accelerator](./cmd/plugin/accelerator/README.md)

## Install

The [Tanzu CLI](https://docs.vmware.com/en/VMware-Tanzu-Application-Platform/0.1/tap-0-1/GUID-install.html#install-the-tanzu-cli-and-package-plugin-4) is required to use the Accelerator CLI plugin.

### From a pre-built distribution

Download `tanzu-accelerator-plugin-<version>.tar.gz` from the most recent release listed on the [App Accelerator tanzu CLI plugin releases page](https://github.com/pivotal/acc-tanzu-cli/releases).

Extract the archive to a local directory:

```sh
tar -zxvf tanzu-accelerator-plugin-*.tar.gz
```

Install the accelerator plugin:

```sh
tanzu plugin install accelerator --local ./artifacts --version v0.5.0-dev
```

### Build from source

See the [Development Guide](./DEVELOPMENT.md).
