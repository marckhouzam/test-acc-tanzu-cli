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
			GivenObjects: []client.Object{
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      fragmentName,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
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
					Kind:      "Fragment",
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
			GivenObjects: []client.Object{
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      fragmentName,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta2.GitRepositoryRef{
								Branch: gitBranch,
							},
						},
					},
				},
			},
			ExpectOutput: "accelerator fragment non-existent not found\n",
			ShouldError:  true,
		},
		{
			Name: "Delete Fragment",
			Args: []string{fragmentName},
			GivenObjects: []client.Object{
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      fragmentName,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
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
					Kind:      "Fragment",
					Namespace: namespace,
					Name:      fragmentName,
				},
			},
		},
	}
	table.Run(t, scheme, FragmentDeleteCmd)

}
