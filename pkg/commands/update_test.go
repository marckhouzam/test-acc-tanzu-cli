package commands

import (
	"testing"
	"time"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	clitesting "github.com/vmware-tanzu-private/tanzu-cli-apps-plugins/pkg/cli-runtime/testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestUpdateCmd(t *testing.T) {
	acceleratorName := "test-accelerator"
	testDescription := "another description"
	namespace := "default"
	expectedDuration, _ := time.ParseDuration("2m")
	expectedInterval := &metav1.Duration{
		Duration: expectedDuration,
	}
	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)

	table := clitesting.CommandTestSuite{
		{
			Name:        "Missing args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name:         "Invalid accelerator",
			Args:         []string{"non-existent"},
			ShouldError:  true,
			ExpectOutput: "accelerator non-existent not found\n",
		},
		{
			Name: "Error updating accelerator",
			Args: []string{acceleratorName, "--description", testDescription},
			WithReactors: []clitesting.ReactionFunc{
				clitesting.InduceFailure("update", "Accelerator"),
			},
			GivenObjects: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Description: "first description",
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta1.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				}),
			},
			ExpectUpdates: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Description: testDescription,
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta1.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				}),
			},
			ShouldError:  true,
			ExpectOutput: "there was an error updating accelerator test-accelerator\n",
		},
		{
			Name: "Updates accelerator",
			Args: []string{acceleratorName, "--description", testDescription, "--git-interval", "2m"},
			GivenObjects: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Description: "first description",
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta1.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				}),
			},
			ExpectUpdates: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Description: testDescription,
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta1.GitRepositoryRef{
								Branch: "main",
							},
							Interval: expectedInterval,
						},
					},
				}),
			},
			ExpectOutput: "accelerator test-accelerator updated successfully\n",
		},
	}

	table.Run(t, scheme, UpdateCmd)
}
