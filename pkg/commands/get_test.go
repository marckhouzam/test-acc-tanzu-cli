package commands

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta2"
	"github.com/pivotal/acc-controller/sourcecontroller/api/v1alpha1"
	cli "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime"
	clitesting "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime/testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
					{
						Name:           "mock-empty-tags",
						IconUrl:        "http://icon-url.png",
						SourceUrl:      "http://www.test.com",
						SourceBranch:   "main",
						SourceTag:      "v1.0.0",
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
		mockOptions := OptionsResponse{
			Options: []Option{
				{
					Name:         "test-option",
					DefaultValue: "test",
					Display:      true,
					DataType:     "choices",
					Choices: []Choice{
						{
							Text:  "first",
							Value: "first",
						},
					},
				},
				{
					Name:         "test-option-bool",
					DefaultValue: true,
					Display:      true,
					DataType:     "boolean",
				},
			},
		}
		emptyOptions := OptionsResponse{}
		var mockResponse []byte
		if strings.Contains(r.URL.Path, "options") && strings.Contains(r.URL.RawQuery, "empty") {
			mockResponse, _ = json.Marshal(emptyOptions)
		} else if strings.Contains(r.URL.Path, "options") {
			mockResponse, _ = json.Marshal(mockOptions)
		} else {
			mockResponse, _ = json.Marshal(mockAccelerator)
		}

		w.Write(mockResponse)
	}))
	os.Setenv("ACC_SERVER_URL", ts.URL)
	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)
	acceleratorName := "test-accelerator"
	namespace := "accelerator-system"
	ignore := ".ignore"
	duration, _ := time.ParseDuration("2m")

	testAccelerator := acceleratorv1alpha1.Accelerator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      acceleratorName,
			Namespace: namespace,
		},
		Spec: acceleratorv1alpha1.AcceleratorSpec{
			Git: &acceleratorv1alpha1.Git{
				Ignore: &ignore,
				URL:    "http://www.test.com",
				Reference: &v1beta2.GitRepositoryRef{
					Branch: "main",
					Tag:    "v1.0.0",
				},
				Interval: &metav1.Duration{
					Duration: duration,
				},
			},
		},
		Status: acceleratorv1alpha1.AcceleratorStatus{
			Description: "Lorem Ipsum",
			DisplayName: "Test Accelerator",
			IconUrl:     "http://icon.png",
			Tags:        []string{"first", "second"},
			ArtifactInfo: acceleratorv1alpha1.ArtifactInfo{
				Ready:   true,
				Message: "test",
				URL:     "http://www.test.com",
				Imports: map[string]string{"java-version": "http://www.example.com"},
			},
			Options: `[{"defaultValue": "","name":"test","label":"test"}]`,
		},
	}

	testAcceleratorEmptyValues := acceleratorv1alpha1.Accelerator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      acceleratorName,
			Namespace: namespace,
		},
		Spec: acceleratorv1alpha1.AcceleratorSpec{
			Git: &acceleratorv1alpha1.Git{
				Ignore: &ignore,
				URL:    "http://www.test.com",
				Reference: &v1beta2.GitRepositoryRef{
					Branch: "main",
					Tag:    "v1.0.0",
				},
			},
		},
		Status: acceleratorv1alpha1.AcceleratorStatus{
			Description: "Lorem Ipsum",
			DisplayName: "Test Accelerator",
			IconUrl:     "http://icon.png",
			ArtifactInfo: acceleratorv1alpha1.ArtifactInfo{
				Ready:   true,
				Message: "test",
				URL:     "http://www.test.com",
			},
		},
	}

	testAcceleratorImage := acceleratorv1alpha1.Accelerator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      acceleratorName,
			Namespace: namespace,
		},
		Spec: acceleratorv1alpha1.AcceleratorSpec{
			Source: &v1alpha1.ImageRepositorySpec{
				Image: "test-image",
			},
		},
		Status: acceleratorv1alpha1.AcceleratorStatus{
			Description: "Lorem Ipsum",
			DisplayName: "Test Accelerator",
			IconUrl:     "http://icon.png",
			Tags:        []string{"first", "second"},
			ArtifactInfo: acceleratorv1alpha1.ArtifactInfo{
				Ready:   true,
				Message: "test",
				URL:     "http://www.test.com",
			},
			Options: `[{"defaultValue": "","name":"test","label":"test"}]`,
		},
	}

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
			ExpectOutput: "accelerator non-existent not found.\n",
		},
		{
			Name: "Wrong acc server URL",
			Prepare: func(t *testing.T, ctx context.Context, config *cli.Config, tc *clitesting.CommandTestCase) (context.Context, error) {
				os.Setenv("ACC_SERVER_URL", "http://not-found")
				return ctx, nil
			},
			Args: []string{"error"},
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
			Name: "Error getting accelerator from context",
			Args: []string{acceleratorName, "--from-context"},
			WithReactors: []clitesting.ReactionFunc{
				clitesting.InduceFailure("get", "Accelerator"),
			},
			GivenObjects: []client.Object{
				&testAccelerator,
			},
			ShouldError:  true,
			ExpectOutput: "Error getting accelerator test-accelerator\n",
		},
		{
			Name: "Get an accelerator from context",
			Args: []string{acceleratorName, "--from-context"},
			GivenObjects: []client.Object{
				&testAccelerator,
			},
			ExpectOutput: `
name: test-accelerator
namespace: accelerator-system
description: Lorem Ipsum
displayName: Test Accelerator
git:
  interval: 2m0s
  ignore: .ignore
  url: http://www.test.com
  ref:
    branch: main
    tag: v1.0.0
tags:
- first
- second
ready: true
options:
- defaultValue: ""
  label: test
  name: test
artifact:
  message: test
  ready: true
imports:
  java-version
`,
		},
		{
			Name: "Get an accelerator from context with verbose flag",
			Args: []string{acceleratorName, "--from-context", "--verbose"},
			GivenObjects: []client.Object{
				&testAccelerator,
			},
			ExpectOutput: `
name: test-accelerator
namespace: accelerator-system
description: Lorem Ipsum
displayName: Test Accelerator
iconUrl: http://icon.png
git:
  interval: 2m0s
  ignore: .ignore
  url: http://www.test.com
  ref:
    branch: main
    tag: v1.0.0
tags:
- first
- second
ready: true
options:
- defaultValue: ""
  label: test
  name: test
artifact:
  message: test
  ready: true
  url: http://www.test.com
imports:
  java-version
`,
		},
		{
			Name: "Get an accelerator with empty values from context",
			Args: []string{acceleratorName, "--from-context"},
			GivenObjects: []client.Object{
				&testAcceleratorEmptyValues,
			},
			ExpectOutput: `
name: test-accelerator
namespace: accelerator-system
description: Lorem Ipsum
displayName: Test Accelerator
git:
  ignore: .ignore
  url: http://www.test.com
  ref:
    branch: main
    tag: v1.0.0
tags: []
ready: true
options: []
artifact:
  message: test
  ready: true
imports:
  None
`,
		},
		{
			Name: "Get an accelerator with image from context",
			Args: []string{acceleratorName, "--from-context"},
			GivenObjects: []client.Object{
				&testAcceleratorImage,
			},
			ExpectOutput: `
name: test-accelerator
namespace: accelerator-system
description: Lorem Ipsum
displayName: Test Accelerator
source:
  image: test-image
tags:
- first
- second
ready: true
options:
- defaultValue: ""
  label: test
  name: test
artifact:
  message: test
  ready: true
imports:
  None
`,
		},
		{
			Name: "Get accelerators from server-url",
			Args: []string{"mock", "--server-url", ts.URL},
			ExpectOutput: `
name: mock
description: Lorem Ipsum
displayName: Mock
sourceUrl: http://www.test.com
tags:
- first
- second
ready: true
options:
- name: test-option
  defaultValue: test
  display: true
  dataType: choices
  choices:
  - text: first
    value: first
- name: test-option-bool
  defaultValue: true
  display: true
  dataType: boolean
  choices: []
artifact:
  message: Lorem Ipsum archive
  ready: true
`,
		},
		{
			Name: "Get accelerators from server-url with verbose flag",
			Args: []string{"mock", "--server-url", ts.URL, "--verbose"},
			ExpectOutput: `
name: mock
description: Lorem Ipsum
displayName: Mock
iconUrl: http://icon-url.png
sourceUrl: http://www.test.com
tags:
- first
- second
ready: true
options:
- name: test-option
  defaultValue: test
  display: true
  dataType: choices
  choices:
  - text: first
    value: first
- name: test-option-bool
  defaultValue: true
  display: true
  dataType: boolean
  choices: []
artifact:
  message: Lorem Ipsum archive
  ready: true
  url: http://archive.tar.gz
`,
		},
		{
			Name: "Get empty tags accelerators from server-url",
			Args: []string{"mock-empty-tags", "--server-url", ts.URL},
			ExpectOutput: `
name: mock-empty-tags
description: Lorem Ipsum
displayName: Mock
sourceUrl: http://www.test.com
tags: []
ready: true
options: []
artifact:
  message: Lorem Ipsum archive
  ready: true
`,
		},
	}

	table.Run(t, scheme, GetCmd)
}
