package commands

import (
	"context"
	"net/url"
	"testing"

	ggcrregistry "github.com/google/go-containerregistry/pkg/registry"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	cli "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
	clitesting "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime/testing"
	"github.com/vmware-tanzu/apps-cli-plugin/pkg/source"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func TestPushCmd(t *testing.T) {
	registry, err := ggcrregistry.TLS("localhost")
	utilruntime.Must(err)
	defer registry.Close()
	u, err := url.Parse(registry.URL)
	utilruntime.Must(err)
	registryHost := u.Host
	localPath := "testdata/test-acc"
	imageName := registryHost + "/test-image:test"
	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)

	c := cli.NewDefaultConfig("test", scheme)
	c.Client = clitesting.NewFakeCliClient(clitesting.NewFakeClient(scheme))

	table := clitesting.CommandTestSuite{
		{
			Name:        "Missing options",
			Args:        []string{},
			ShouldError: true,
		},

		{
			Name:        "Missing --local-path",
			Args:        []string{"--source-image", imageName},
			ShouldError: true,
		},

		{
			Name:        "Missing --source-image",
			Args:        []string{"--local-path", localPath},
			ShouldError: true,
		},

		{
			Config: c,
			Prepare: func(t *testing.T, ctx context.Context, config *cli.Config, tc *clitesting.CommandTestCase) (context.Context, error) {
				return source.StashGgcrRemoteOptions(ctx, remote.WithTransport(registry.Client().Transport)), nil
			},
			Name:         "Push to source-image",
			Args:         []string{"--local-path", localPath, "--source-image", imageName},
			ShouldError:  false,
			ExpectOutput: "publishing accelerator source in \"testdata/test-acc\" to \"" + imageName + "\"...\npublished accelerator\n",
		},
	}

	table.Run(t, scheme, PushCmd)
}
