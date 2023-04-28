package commands

import (
	"testing"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta2"
	rtesting "github.com/vmware-labs/reconciler-runtime/testing"
	clitesting "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime/testing"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
			GivenObjects: []client.Object{
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
							},
						},
					},
				},
			},
			ExpectDeletes: []rtesting.DeleteRef{
				{
					Group:     "accelerator.apps.tanzu.vmware.com",
					Kind:      "Accelerator",
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
			GivenObjects: []client.Object{
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
							},
						},
					},
				},
			},
			ExpectOutput: "accelerator non-existent not found\n",
			ShouldError:  true,
		},
		{
			Name: "Delete Accelerator",
			Args: []string{acceleratorName},
			GivenObjects: []client.Object{
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
							},
						},
					},
				},
			},
			ExpectDeletes: []rtesting.DeleteRef{
				{
					Group:     "accelerator.apps.tanzu.vmware.com",
					Kind:      "Accelerator",
					Namespace: namespace,
					Name:      acceleratorName,
				},
			},
		},
	}
	table.Run(t, scheme, DeleteCmd)

}
