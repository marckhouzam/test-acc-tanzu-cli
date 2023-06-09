package commands

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/fluxcd/pkg/apis/meta"
	ggcrregistry "github.com/google/go-containerregistry/pkg/registry"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta2"
	"github.com/pivotal/acc-controller/sourcecontroller/api/v1alpha1"
	cli "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
	clitesting "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime/testing"
	"github.com/vmware-tanzu/apps-cli-plugin/pkg/source"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestCreateCommand(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)

	acceleratorName := "test-accelerator"
	gitRepoUrl := "https://www.test.com"
	imageName := "test-image"
	noGitBranch := "main"
	noGitTag := ""
	gitBranch := "main"
	gitTag := "v0.0.1"
	namespace := "accelerator-system"
	interval := "2m"
	secretRef := "mysecret"
	expectedDuration, _ := time.ParseDuration(interval)

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
			ExpectCreates: []client.Object{
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      acceleratorName,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta2.GitRepositoryRef{
								Branch: noGitBranch,
								Tag:    noGitTag,
							},
						},
					},
				},
			},
			ExpectOutput: "Error creating accelerator test-accelerator\n",
			ShouldError:  true,
		},
		{
			Name: "Create Accelerator just GitRepository",
			Args: []string{acceleratorName, "--git-repository", gitRepoUrl},
			ExpectCreates: []client.Object{
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      acceleratorName,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta2.GitRepositoryRef{
								Branch: noGitBranch,
								Tag:    noGitTag,
							},
						},
					},
				},
			},
			ExpectOutput: "created accelerator test-accelerator in namespace accelerator-system\n",
		},
		{
			Name: "Create Accelerator Image with Secret ref",
			Args: []string{acceleratorName, "--source-image", imageName, "--secret-ref", secretRef, "--interval", interval},
			ExpectCreates: []client.Object{
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      acceleratorName,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Source: &v1alpha1.ImageRepositorySpec{
							Image: imageName,
							ImagePullSecrets: []corev1.LocalObjectReference{
								{
									Name: secretRef,
								},
							},
							Interval: &v1.Duration{
								Duration: expectedDuration,
							},
						},
					},
				},
			},
			ExpectOutput: "created accelerator test-accelerator in namespace accelerator-system\n",
		},
		{
			Name: "Create Accelerator with Branch and Tag and Secret ref",
			Args: []string{acceleratorName,
				"--git-repository", gitRepoUrl,
				"--git-branch", gitBranch,
				"--git-tag", gitTag,
				"--interval", interval,
				"--secret-ref", secretRef,
			},
			ExpectCreates: []client.Object{
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      acceleratorName,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta2.GitRepositoryRef{
								Branch: gitBranch,
								Tag:    gitTag,
							},
							SecretRef: &meta.LocalObjectReference{
								Name: secretRef,
							},
							Interval: &v1.Duration{
								Duration: expectedDuration,
							},
						},
					},
				},
			},
			ExpectOutput: "created accelerator test-accelerator in namespace accelerator-system\n",
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
	namespace := "accelerator-system"

	table := clitesting.CommandTestSuite{
		{
			Name:        "Missing args",
			Args:        []string{"--source-image", imageName, "--local-path", localPath},
			ShouldError: true,
		},
		{
			Config: c,
			Prepare: func(t *testing.T, ctx context.Context, config *cli.Config, tc *clitesting.CommandTestCase) (context.Context, error) {
				return source.StashContainerRemoteTransport(ctx, registry.Client().Transport), nil
			},
			Name: "Create Accelerator Image from Local Path",
			Args: []string{acceleratorName, "--source-image", imageName, "--local-path", localPath},
			ExpectCreates: []client.Object{
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      acceleratorName,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Source: &v1alpha1.ImageRepositorySpec{
							Image: imageName,
						},
					},
				},
			},
			ExpectOutput: "publishing accelerator source in \"testdata/test-acc\" to \"" + imageName + "\"...\npublished accelerator\ncreated accelerator test-accelerator in namespace accelerator-system\n",
		},
	}
	table.Run(t, scheme, CreateCmd)
}
