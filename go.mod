module github.com/pivotal/acc-tanzu-cli

go 1.16

require (
	github.com/fluxcd/pkg/apis/meta v0.9.0
	github.com/fsnotify/fsnotify v1.5.0 // indirect
	github.com/google/go-containerregistry v0.6.0
	github.com/imdario/mergo v0.3.12
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/pivotal/acc-controller v0.3.1-0.20211008011320-aad710cbc6c1
	github.com/spf13/cobra v1.2.1
	github.com/vmware-tanzu/tanzu-cli-apps-plugins v0.2.1-0.20211007192245-181c97eeb1d0
	github.com/vmware-tanzu/tanzu-framework v0.6.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	sigs.k8s.io/controller-runtime v0.10.2
)
