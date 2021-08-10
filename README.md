# App Accelerator tanzu CLI plugin

To check the plugin documentation go to `cmd/plugin/accelerator/README.md`
# Setup

Before installing the plugin in your tanzu instance, you need to follow the steps described here https://github.com/vmware-tanzu/tanzu-framework/blob/main/docs/cli/getting-started.md
for running plugins locally

# Build the plugin

to build the plugin you need to:

 - Run the command `go mod vendor`
 - Run the `make build` task
 - The build task will create the `artifacts/` directory
 - To install the plugin to tanzu, you need to run `tanzu plugin install accelerator --local ./artifacts`