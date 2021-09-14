package commands

import (
	"testing"

	clitesting "github.com/vmware-tanzu-private/tanzu-cli-apps-plugins/pkg/cli-runtime/testing"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestAcceleratorListCommand(t *testing.T) {
	acceleratorName := "test-accelerator"
	namespace := "default"

	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)

	table := clitesting.CommandTestSuite{
		{
			Name: "empty",
			Args: []string{},
			ExpectOutput: `
No accelerators found.
`,
		},
		{
			Name: "Error listing accelerators",
			Args: []string{},
			WithReactors: []clitesting.ReactionFunc{
				clitesting.InduceFailure("list", "AcceleratorList"),
			},
			GivenObjects: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta1.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				}),
			},
			ShouldError:  true,
			ExpectOutput: "There was an error listing accelerators\n",
		},
		{
			Name: "List accelerators",
			Args: []string{},
			GivenObjects: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta1.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				}),
			},
			ExpectOutput: `
NAME               GIT REPOSITORY         BRANCH   TAG
test-accelerator   https://www.test.com   main     
`,
		},
	}

	table.Run(t, scheme, ListCmd)
}
