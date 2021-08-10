# App Accelerator tanzu CLI plugin

To check the plugin documentation go to `cmd/plugin/accelerator/README.md`

# Build the plugin1

to build the plugin you need to:

 - Run the command `go mod vendor`
 - Run the `make build` task
 - The build task will create the `artifacts/` directory
 - To install the plugin to tanzu, you need to run `tanzu plugin install accelerator --local ./artifacts`