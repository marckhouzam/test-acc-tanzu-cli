package commands

import (
	"testing"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	clitesting "github.com/vmware-tanzu-private/tanzu-cli-apps-plugins/pkg/cli-runtime/testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestGetCommand(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)
	acceleratorName := "test-accelerator"
	namespace := "default"

	table := clitesting.CommandTestSuite{
		{
			Name:        "Missing args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "Error getting accelerator",
			Args: []string{acceleratorName},
			WithReactors: []clitesting.ReactionFunc{
				clitesting.InduceFailure("get", "Accelerator"),
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
			ExpectOutput: "Error getting accelerator test-accelerator\n",
		},
		{
			Name: "Get an accelerator",
			Args: []string{acceleratorName},
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

	table.Run(t, scheme, GetCmd)
}
