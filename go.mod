module github.com/pivotal/acc-tanzu-cli

go 1.16

require (
	github.com/fluxcd/pkg/apis/meta v0.9.0
	github.com/fsnotify/fsnotify v1.5.0 // indirect
	github.com/google/go-containerregistry v0.7.0
	github.com/imdario/mergo v0.3.12
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/pivotal/acc-controller v0.5.0
	github.com/spf13/cobra v1.2.1
	github.com/vmware-tanzu/tanzu-cli-apps-plugins v0.3.0
	github.com/vmware-tanzu/tanzu-framework v0.10.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.22.4
	k8s.io/client-go v0.22.4
	sigs.k8s.io/controller-runtime v0.10.3
)

replace go.mongodb.org/mongo-driver v1.1.2 => go.mongodb.org/mongo-driver v1.5.1

replace github.com/containerd/containerd v1.5.7 => github.com/containerd/containerd v1.5.9
