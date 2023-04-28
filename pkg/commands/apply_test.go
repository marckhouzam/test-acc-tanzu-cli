package commands

import (
	"testing"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta2"
	clitesting "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime/testing"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestApplyCommand(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)

	acceleratorName := "test-accelerator"
	fragmentName := "test-fragment"
	gitRepoUrl := "https://www.test.com"
	gitBranch := "main"
	namespace := "accelerator-system"
	acceleratorFilename := "testdata/test-accelerator.yml"
	fragmentFilename := "testdata/test-fragment.yml"

	table := clitesting.CommandTestSuite{
		{
			Name:        "Missing arg",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name:         "Apply missing file",
			Args:         []string{acceleratorName, "--filename", "testdata/test.yml"},
			ShouldError:  true,
			ExpectOutput: "Error loading file testdata/test.yml\n",
		},
		{
			Name: "Create Accelerator",
			Args: []string{acceleratorName, "--filename", acceleratorFilename},
			ExpectCreates: []client.Object{
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
			ExpectOutput: "created accelerator test-accelerator in namespace accelerator-system\n",
		},
		{
			Name: "Update Accelerator",
			Args: []string{acceleratorName, "--filename", acceleratorFilename},
			GivenObjects: []client.Object{
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: v1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta2.GitRepositoryRef{
								Branch: "not-main",
							},
						},
					},
				},
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
							Reference: &v1beta2.GitRepositoryRef{
								Branch: gitBranch,
							},
						},
					},
				},
			},
			ExpectOutput: "updated accelerator test-accelerator in namespace accelerator-system\n",
		},
		{
			Name: "Create Fragment",
			Args: []string{fragmentName, "--filename", fragmentFilename},
			ExpectCreates: []client.Object{
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
			ExpectOutput: "created accelerator fragment test-fragment in namespace accelerator-system\n",
		},
		{
			Name: "Update Fragment",
			Args: []string{fragmentName, "--filename", fragmentFilename},
			GivenObjects: []client.Object{
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: v1.ObjectMeta{
						Name:      fragmentName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta2.GitRepositoryRef{
								Branch: "not-main",
							},
						},
					},
				},
			},
			ExpectUpdates: []client.Object{
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: v1.ObjectMeta{
						Name:      fragmentName,
						Namespace: namespace,
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
			ExpectOutput: "updated accelerator fragment test-fragment in namespace accelerator-system\n",
		},
	}
	table.Run(t, scheme, ApplyCmd)
}
