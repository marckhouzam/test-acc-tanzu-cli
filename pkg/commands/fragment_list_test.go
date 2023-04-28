package commands

import (
	"testing"

	clitesting "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestFragmentListCommand(t *testing.T) {
	fragmentName := "test-fragment"
	namespace := "accelerator-system"

	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)

	table := clitesting.CommandTestSuite{
		{
			Name: "empty from context",
			Args: []string{},
			ExpectOutput: `
No accelerator fragments found.
`,
			ShouldError: false,
		},
		{
			Name: "Error listing accelerator fragments from context",
			Args: []string{},
			WithReactors: []clitesting.ReactionFunc{
				clitesting.InduceFailure("list", "FragmentList"),
			},
			GivenObjects: []client.Object{
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      fragmentName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta2.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				},
			},
			ShouldError:  true,
			ExpectOutput: "There was an error listing accelerator fragments\n",
		},
		{
			Name: "List accelerator fragments from context",
			Args: []string{},
			GivenObjects: []client.Object{
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      fragmentName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta2.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				},
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-fragment",
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta2.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				},
			},
			ExpectOutput: `
NAME               READY
another-fragment   unknown
test-fragment      unknown
`,
		},
		{
			Name: "List accelerator fragments from context with verbose flag",
			Args: []string{"--verbose"},
			GivenObjects: []client.Object{
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      fragmentName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta2.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				},
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-fragment",
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta2.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				},
			},
			ExpectOutput: `
NAME               READY     REPOSITORY
another-fragment   unknown   https://www.test.com:main
test-fragment      unknown   https://www.test.com:main
`,
		},
		{
			Name: "List accelerator fragments from context",
			Args: []string{},
			GivenObjects: []client.Object{
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      fragmentName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta2.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				},
			},
			ExpectOutput: `
NAME            READY
test-fragment   unknown
`,
		},
	}
	table.Run(t, scheme, FragmentListCmd)
}
