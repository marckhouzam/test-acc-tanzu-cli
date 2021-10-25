package commands

import (
	"testing"
	"time"

	"github.com/fluxcd/pkg/apis/meta"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	"github.com/pivotal/acc-controller/sourcecontroller/api/v1alpha1"
	clitesting "github.com/vmware-tanzu/tanzu-cli-apps-plugins/pkg/cli-runtime/testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestUpdateCmd(t *testing.T) {
	acceleratorName := "test-accelerator"
	testDescription := "another description"
	namespace := "default"
	repositoryUrl := "http://www.test.com"
	imageName := "test-image"
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
			Args: []string{acceleratorName, "--description", testDescription, "--interval", interval},
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
							Interval: expectedInterval,
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
		{
			Name: "Updates repo url from accelerator",
			Args: []string{acceleratorName, "--git-repository", repositoryUrl, "--secret-ref", secretRef},
			GivenObjects: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Description: "first description",
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
						Description: "first description",
						Git: &acceleratorv1alpha1.Git{
							URL: "http://www.test.com",
							Reference: &v1beta1.GitRepositoryRef{
								Branch: "main",
							},
							SecretRef: &meta.LocalObjectReference{
								Name: secretRef,
							},
						},
					},
				}),
			},
			ExpectOutput: "accelerator test-accelerator updated successfully\n",
		},
		{
			Name: "Updates repo image name from accelerator",
			Args: []string{acceleratorName, "--source-image", imageName, "--secret-ref", secretRef, "--interval", interval},
			GivenObjects: []clitesting.Factory{
				clitesting.Wrapper(&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Description: testDescription,
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
						Source: &v1alpha1.ImageRepositorySpec{
							Image: imageName,
							ImagePullSecrets: []meta.LocalObjectReference{
								{
									Name: secretRef,
								},
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
