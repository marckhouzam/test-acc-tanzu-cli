package commands

import (
	"context"
	"net/url"
	"testing"
	"time"

	ggcrregistry "github.com/google/go-containerregistry/pkg/registry"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	"github.com/pivotal/acc-controller/sourcecontroller/api/v1alpha1"
	cli "github.com/vmware-tanzu/tanzu-cli-apps-plugins/pkg/cli-runtime"
	clitesting "github.com/vmware-tanzu/tanzu-cli-apps-plugins/pkg/cli-runtime/testing"
	"github.com/vmware-tanzu/tanzu-cli-apps-plugins/pkg/source"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func TestCreateCommand(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)

	acceleratorName := "test-accelerator"
	gitRepoUrl := "https://www.test.com"
	imageName := "test-image"
	noGitBranch := ""
	noGitTag := ""
	gitBranch := "main"
	gitTag := "v0.0.1"
	namespace := "default"
	gitInterval := "2m"
	expectedDuration, _ := time.ParseDuration(gitInterval)

	table := clitesting.CommandTestSuite{
		{
			Name:        "Missing args",
			Args:        []string{"--git-repository", gitRepoUrl},
			ShouldError: true,
		},
		{
			Name: "Error creating accelerator",
			Args: []string{acceleratorName, "--git-repository", gitRepoUrl},
			WithReactors: []clitesting.ReactionFunc{
				clitesting.InduceFailure("create", "Accelerator"),
			},
			ExpectCreates: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Accelerator{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      acceleratorName,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta1.GitRepositoryRef{
								Branch: noGitBranch,
								Tag:    noGitTag,
							},
						},
					},
				}),
			},
			ExpectOutput: "Error creating accelerator test-accelerator\n",
			ShouldError:  true,
		},
		{
			Name: "Create Accelerator just GitRepository",
			Args: []string{acceleratorName, "--git-repository", gitRepoUrl},
			ExpectCreates: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Accelerator{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      acceleratorName,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta1.GitRepositoryRef{
								Branch: noGitBranch,
								Tag:    noGitTag,
							},
						},
					},
				}),
			},
			ExpectOutput: "created accelerator test-accelerator in namespace default\n",
		},
		{
			Name: "Create Accelerator Image",
			Args: []string{acceleratorName, "--source-image", imageName},
			ExpectCreates: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Accelerator{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      acceleratorName,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Source: &v1alpha1.ImageRepositorySpec{
							Image: imageName,
						},
					},
				}),
			},
			ExpectOutput: "created accelerator test-accelerator in namespace default\n",
		},
		{
			Name: "Create Accelerator with Branch and Tag",
			Args: []string{acceleratorName,
				"--git-repository", gitRepoUrl,
				"--git-branch", gitBranch,
				"--git-tag", gitTag,
				"--git-interval", gitInterval,
			},
			ExpectCreates: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Accelerator{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      acceleratorName,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta1.GitRepositoryRef{
								Branch: gitBranch,
								Tag:    gitTag,
							},
							Interval: &v1.Duration{
								Duration: expectedDuration,
							},
						},
					},
				}),
			},
			ExpectOutput: "created accelerator test-accelerator in namespace default\n",
		},
	}
	table.Run(t, scheme, CreateCmd)
}

func TestCreateCommandLocalPath(t *testing.T) {
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

	acceleratorName := "test-accelerator"
	namespace := "default"

	table := clitesting.CommandTestSuite{
		{
			Name:        "Missing args",
			Args:        []string{"--source-image", imageName, "--local-path", localPath},
			ShouldError: true,
		},
		{
			Config: c,
			Prepare: func(t *testing.T, ctx context.Context, config *cli.Config, tc *clitesting.CommandTestCase) (context.Context, error) {
				return source.StashGgcrRemoteOptions(ctx, remote.WithTransport(registry.Client().Transport)), nil
			},
			Name: "Create Accelerator Image from Local Path",
			Args: []string{acceleratorName, "--source-image", imageName, "--local-path", localPath},
			ExpectCreates: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Accelerator{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      acceleratorName,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Source: &v1alpha1.ImageRepositorySpec{
							Image: imageName,
						},
					},
				}),
			},
			ExpectOutput: "publishing accelerator source in \"testdata/test-acc\" to \"" + imageName + "\"...\npublished accelerator\ncreated accelerator test-accelerator in namespace default\n",
		},
	}
	table.Run(t, scheme, CreateCmd)
}
