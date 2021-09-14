# App Accelerator tanzu CLI plugin - Development Guide

Before installing the plugin in your tanzu instance, you need to follow the steps described here https://github.com/vmware-tanzu/tanzu-framework/blob/main/docs/cli/getting-started.md
for running plugins locally

# Build the plugin

to build the plugin you need to:

 - Run the command `go mod vendor`
 - Run the `make build` task
 - The build task will create the `artifacts/` directory
 - To install the plugin to tanzu, you need to run `tanzu plugin install accelerator --local ./artifacts`

# Troubleshooting
If you get the error `fatal: could not read Username for 'https://github.com': terminal prompts disabled` while downloading the go dependencies
run the following command
```
git config --global url."git@github.com:".insteadOf https://github.com/
```