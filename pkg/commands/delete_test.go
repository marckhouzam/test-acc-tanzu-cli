package commands

import (
	"testing"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	clitesting "github.com/vmware-tanzu/tanzu-cli-apps-plugins/pkg/cli-runtime/testing"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestDeleteCommand(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)

	acceleratorName := "test-accelerator"
	acceleratorNotFound := "non-existent"
	gitRepoUrl := "https://www.test.com"
	gitBranch := "main"
	namespace := "accelerator-system"

	table := clitesting.CommandTestSuite{
		{
			Name:        "Missing args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "Error Deleting Accelerator",
			Args: []string{acceleratorName},
			WithReactors: []clitesting.ReactionFunc{
				clitesting.InduceFailure("delete", "Accelerator"),
			},
			GivenObjects: []clitesting.Factory{
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
							},
						},
					},
				}),
			},
			ExpectDeletes: []clitesting.DeleteRef{
				{
					Group:     "accelerator.apps.tanzu.vmware.com",
					Resource:  "Accelerator",
					Namespace: namespace,
					Name:      acceleratorName,
				},
			},
			ExpectOutput: "There was a problem trying to delete accelerator test-accelerator\n",
			ShouldError:  true,
		},
		{
			Name: "Error accelerator not found",
			Args: []string{acceleratorNotFound},
			GivenObjects: []clitesting.Factory{
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
							},
						},
					},
				}),
			},
			ExpectOutput: "accelerator non-existent not found\n",
			ShouldError:  true,
		},
		{
			Name: "Delete Accelerator",
			Args: []string{acceleratorName},
			GivenObjects: []clitesting.Factory{
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
							},
						},
					},
				}),
			},
			ExpectDeletes: []clitesting.DeleteRef{
				{
					Group:     "accelerator.apps.tanzu.vmware.com",
					Resource:  "Accelerator",
					Namespace: namespace,
					Name:      acceleratorName,
				},
			},
		},
	}
	table.Run(t, scheme, DeleteCmd)

}
