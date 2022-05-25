package commands

import (
	"testing"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	clitesting "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime/testing"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestFragmentDeleteCommand(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)

	fragmentName := "test-fragment"
	fragmentNotFound := "non-existent"
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
			Name: "Error Deleting Fragment",
			Args: []string{fragmentName},
			WithReactors: []clitesting.ReactionFunc{
				clitesting.InduceFailure("delete", "Fragment"),
			},
			GivenObjects: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Fragment{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      fragmentName,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
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
					Resource:  "Fragment",
					Namespace: namespace,
					Name:      fragmentName,
				},
			},
			ExpectOutput: "There was a problem trying to delete accelerator fragment test-fragment\n",
			ShouldError:  true,
		},
		{
			Name: "Error accelerator fragment not found",
			Args: []string{fragmentNotFound},
			GivenObjects: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Fragment{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      fragmentName,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta1.GitRepositoryRef{
								Branch: gitBranch,
							},
						},
					},
				}),
			},
			ExpectOutput: "accelerator fragment non-existent not found\n",
			ShouldError:  true,
		},
		{
			Name: "Delete Fragment",
			Args: []string{fragmentName},
			GivenObjects: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Fragment{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      fragmentName,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
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
					Resource:  "Fragment",
					Namespace: namespace,
					Name:      fragmentName,
				},
			},
		},
	}
	table.Run(t, scheme, FragmentDeleteCmd)

}
