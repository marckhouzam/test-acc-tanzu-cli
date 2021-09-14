package commands

import (
	"testing"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	clitesting "github.com/vmware-tanzu-private/tanzu-cli-apps-plugins/pkg/cli-runtime/testing"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestCreateCommand(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)

	acceleratorName := "test-accelerator"
	gitRepoUrl := "https://www.test.com"
	gitBranch := "main"
	namespace := "default"

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
						Git: acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta1.GitRepositoryRef{
								Branch: gitBranch,
							},
						},
					},
				}),
			},
			ExpectOutput: "Error creating accelerator test-accelerator\n",
			ShouldError:  true,
		},
		{
			Name: "Create Accelerator",
			Args: []string{acceleratorName, "--git-repository", gitRepoUrl},
			ExpectCreates: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Accelerator{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      acceleratorName,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta1.GitRepositoryRef{
								Branch: gitBranch,
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
