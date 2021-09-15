package commands

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	clitesting "github.com/vmware-tanzu-private/tanzu-cli-apps-plugins/pkg/cli-runtime/testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestGetCommand(t *testing.T) {
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
	acceleratorName := "test-accelerator"
	namespace := "default"

	table := clitesting.CommandTestSuite{
		{
			Name:        "Missing args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name:         "Error getting accelerator",
			Args:         []string{"non-existent"},
			ShouldError:  true,
			ExpectOutput: "accelertor non-existent not found.\n",
		},
		{
			Name: "Error getting accelerator from context",
			Args: []string{acceleratorName, "--from-context"},
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
			Name: "Get an accelerator from context",
			Args: []string{acceleratorName, "--from-context"},
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
			Name: "Get accelerators from server-url",
			Args: []string{"mock", "--server-url", ts.URL},
			ExpectOutput: `
NAME   GIT REPOSITORY        BRANCH   TAG
mock   http://www.test.com   main     v1.0.0
`,
		},
	}

	table.Run(t, scheme, GetCmd)
}
