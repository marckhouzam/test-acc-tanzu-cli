package commands

import (
	"testing"
	"time"

	"github.com/fluxcd/pkg/apis/meta"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	clitesting "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime/testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestFragmentUpdateCmd(t *testing.T) {
	acceleratorName := "test-fragment"
	namespace := "accelerator-system"
	secretRef := "mysecret"
	interval := "2m"
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
			Name:         "Invalid fragment",
			Args:         []string{"non-existent"},
			ShouldError:  true,
			ExpectOutput: "accelerator fragment non-existent not found\n",
		},
		{
			Name: "Error updating fragment",
			Args: []string{acceleratorName},
			WithReactors: []clitesting.ReactionFunc{
				clitesting.InduceFailure("update", "Fragment"),
			},
			GivenObjects: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Fragment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta1.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				}),
			},
			ExpectUpdates: []client.Object{
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta1.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				},
			},
			ShouldError:  true,
			ExpectOutput: "there was an error updating accelerator fragment test-fragment\n",
		},
		{
			Name: "Updates fragment",
			Args: []string{acceleratorName, "--interval", interval, "--secret-ref", secretRef},
			GivenObjects: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Fragment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta1.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				}),
			},
			ExpectUpdates: []client.Object{
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta1.GitRepositoryRef{
								Branch: "main",
							},
							SecretRef: &meta.LocalObjectReference{
								Name: secretRef,
							},
							Interval: expectedInterval,
						},
					},
				},
			},
			ExpectOutput: "accelerator fragment test-fragment updated successfully\n",
		},
	}

	table.Run(t, scheme, FragmentUpdateCmd)
}
