package commands

import (
	"testing"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	clitesting "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime/testing"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestApplyCommand(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)

	acceleratorName := "test-accelerator"
	gitRepoUrl := "https://www.test.com"
	gitBranch := "main"
	namespace := "accelerator-system"
	filename := "testdata/test.yml"

	table := clitesting.CommandTestSuite{
		{
			Name:        "Missing arg",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "Create Accelerator",
			Args: []string{acceleratorName, "--filename", filename},
			ExpectCreates: []client.Object{
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      acceleratorName,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta1.GitRepositoryRef{
								Branch: gitBranch,
							},
						},
					},
				},
			},
			ExpectOutput: "created accelerator test-accelerator in namespace accelerator-system\n",
		},
		{
			Name: "Update Accelerator",
			Args: []string{acceleratorName, "--filename", filename},
			GivenObjects: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Accelerator{
					ObjectMeta: v1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta1.GitRepositoryRef{
								Branch: "not-main",
							},
						},
					},
				}),
			},
			ExpectUpdates: []client.Object{
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: v1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta1.GitRepositoryRef{
								Branch: gitBranch,
							},
						},
					},
				},
			},
			ExpectOutput: "updated accelerator test-accelerator in namespace accelerator-system\n",
		},
	}
	table.Run(t, scheme, ApplyCmd)
}
