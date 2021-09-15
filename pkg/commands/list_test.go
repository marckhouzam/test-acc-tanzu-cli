package commands

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	cli "github.com/vmware-tanzu-private/tanzu-cli-apps-plugins/pkg/cli-runtime"
	clitesting "github.com/vmware-tanzu-private/tanzu-cli-apps-plugins/pkg/cli-runtime/testing"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestAcceleratorListCommand(t *testing.T) {
	acceleratorName := "test-accelerator"
	namespace := "default"

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
			Name: "empty",
			Prepare: func(t *testing.T, ctx context.Context, config *cli.Config, tc *clitesting.CommandTestCase) (context.Context, error) {
				os.Setenv("ACC_SERVER_URL", "http://not-found")
				return ctx, nil
			},
			Args: []string{},
			ExpectOutput: `
Error getting accelerators from http://not-found
`,
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
no accelerators found.
`,
			ShouldError: true,
		},
		{
			Name: "Error listing accelerators from context",
			Args: []string{"--from-context"},
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
			Name: "List accelerators from context",
			Args: []string{"--from-context"},
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
		{
			Name: "List accelerators server-url",
			Args: []string{"--server-url", ts.URL},
			ExpectOutput: `
NAME   GIT REPOSITORY        BRANCH   TAG
mock   http://www.test.com   main     v1.0.0
`,
		},
		{
			Name: "List accelerators from context",
			Args: []string{"--from-context"},
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
