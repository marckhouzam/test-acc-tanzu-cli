package commands

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	cli "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
	clitesting "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta2"
	"github.com/pivotal/acc-controller/sourcecontroller/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestAcceleratorListCommand(t *testing.T) {
	acceleratorName := "test-accelerator"
	namespace := "accelerator-system"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockAccelerator := UiAcceleratorsApiResponse{
			Emdedded: Embedded{
				Accelerators: []Accelerator{
					{
						Name:           "mock",
						IconUrl:        "http://icon-url.png",
						SourceUrl:      "http://www.test.com",
						SourceBranch:   "main",
						SourceTag:      "v1.0.0",
						Tags:           []string{"first", "second"},
						Description:    "Lorem Ipsum",
						DisplayName:    "Mock",
						Ready:          true,
						ArchiveUrl:     "http://archive.tar.gz",
						ArchiveReady:   true,
						ArchiveMessage: "Lorem Ipsum archive",
					},
				},
			},
		}
		mockResponse, _ := json.Marshal(mockAccelerator)
		w.Write(mockResponse)
	}))

	os.Setenv("ACC_SERVER_URL", ts.URL)
	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)

	table := clitesting.CommandTestSuite{
		{
			Name: "Wrong acc server URL",
			Prepare: func(t *testing.T, ctx context.Context, config *cli.Config, tc *clitesting.CommandTestCase) (context.Context, error) {
				os.Setenv("ACC_SERVER_URL", "http://not-found")
				return ctx, nil
			},
			Args: []string{},
			ExpectOutput: "Error getting accelerators from http://not-found," +
				" check that --server-url or the ACC_SERVER_URL env variable is set with the correct value," +
				" or use the --from-context flag to get the accelerators from your current context\n",
			CleanUp: func(t *testing.T, ctx context.Context, config *cli.Config, tc *clitesting.CommandTestCase) error {
				os.Setenv("ACC_SERVER_URL", ts.URL)
				return nil
			},
			ShouldError: true,
		},
		{
			Name: "empty from context",
			Args: []string{"--from-context"},
			ExpectOutput: `
No accelerators found.
`,
			ShouldError: false,
		},
		{
			Name: "Error listing accelerators from context",
			Args: []string{"--from-context"},
			WithReactors: []clitesting.ReactionFunc{
				clitesting.InduceFailure("list", "AcceleratorList"),
			},
			GivenObjects: []client.Object{
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
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
			ExpectOutput: "There was an error listing accelerators\n",
		},
		{
			Name: "List accelerators from context",
			Args: []string{"--from-context"},
			GivenObjects: []client.Object{
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta2.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				},
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-accelerator",
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta2.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				},
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "image-accelerator",
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Source: &v1alpha1.ImageRepositorySpec{
							Image: "test-image",
						},
					},
				},
			},
			ExpectOutput: `
NAME                  TAGS   READY
another-accelerator   []     unknown
image-accelerator     []     unknown
test-accelerator      []     unknown
`,
		},
		{
			Name: "List accelerators from context with verbose flag",
			Args: []string{"--from-context", "--verbose"},
			GivenObjects: []client.Object{
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta2.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				},
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-accelerator",
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: "https://www.test.com",
							Reference: &v1beta2.GitRepositoryRef{
								Branch: "main",
							},
						},
					},
				},
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "image-accelerator",
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
						Source: &v1alpha1.ImageRepositorySpec{
							Image: "test-image",
						},
					},
				},
			},
			ExpectOutput: `
NAME                  TAGS   READY     REPOSITORY
another-accelerator   []     unknown   https://www.test.com:main
image-accelerator     []     unknown   source-image: test-image
test-accelerator      []     unknown   https://www.test.com:main
`,
		},
		{
			Name: "List accelerators server-url",
			Args: []string{"--server-url", ts.URL},
			ExpectOutput: `
NAME   TAGS             READY
mock   [first second]   true
`,
		},
		{
			Name: "List accelerators server-url with verbose flag",
			Args: []string{"--server-url", ts.URL, "--verbose"},
			ExpectOutput: `
NAME   TAGS             READY   REPOSITORY
mock   [first second]   true    
`,
		},
		{
			Name: "List accelerators from context",
			Args: []string{"--from-context"},
			GivenObjects: []client.Object{
				&acceleratorv1alpha1.Accelerator{
					ObjectMeta: metav1.ObjectMeta{
						Name:      acceleratorName,
						Namespace: namespace,
					},
					Spec: acceleratorv1alpha1.AcceleratorSpec{
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
NAME               TAGS   READY
test-accelerator   []     unknown
`,
		},
	}
	table.Run(t, scheme, ListCmd)
}
