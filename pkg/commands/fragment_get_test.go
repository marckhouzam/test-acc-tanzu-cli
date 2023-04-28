package commands

import (
	"testing"
	"time"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta2"
	clitesting "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime/testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestFragmentGetCommand(t *testing.T) {

	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)
	fragmentName := "test-fragment"
	namespace := "accelerator-system"
	ignore := ".ignore"
	duration, _ := time.ParseDuration("2m")

	testFragment := acceleratorv1alpha1.Fragment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fragmentName,
			Namespace: namespace,
		},
		Spec: acceleratorv1alpha1.FragmentSpec{
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
		Status: acceleratorv1alpha1.FragmentStatus{
			DisplayName: "Test Fragment",
			ArtifactInfo: acceleratorv1alpha1.ArtifactInfo{
				Ready:   true,
				Message: "test",
				URL:     "http://www.test.com",
				Imports: map[string]string{"java-version": "http://www.example.com"},
			},
			Options: `[{"defaultValue": "","name":"test","label":"test"}]`,
		},
	}

	testAcceleratorEmptyValues := acceleratorv1alpha1.Fragment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fragmentName,
			Namespace: namespace,
		},
		Spec: acceleratorv1alpha1.FragmentSpec{
			Git: &acceleratorv1alpha1.Git{
				Ignore: &ignore,
				URL:    "http://www.test.com",
				Reference: &v1beta2.GitRepositoryRef{
					Branch: "main",
					Tag:    "v1.0.0",
				},
			},
		},
		Status: acceleratorv1alpha1.FragmentStatus{
			DisplayName: "Test Fragment",
			ArtifactInfo: acceleratorv1alpha1.ArtifactInfo{
				Ready:   true,
				Message: "test",
				URL:     "http://www.test.com",
			},
		},
	}

	table := clitesting.CommandTestSuite{
		{
			Name:        "Missing args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name:        "Error getting fragment",
			Args:        []string{"non-existent"},
			ShouldError: true,
		},
		{
			Name: "Error getting fragment from context",
			Args: []string{fragmentName},
			WithReactors: []clitesting.ReactionFunc{
				clitesting.InduceFailure("get", "Fragment"),
			},
			GivenObjects: []client.Object{
				&testFragment,
			},
			ShouldError:  true,
			ExpectOutput: "Error getting accelerator fragment test-fragment\n",
		},
		{
			Name: "Get a fragment from context",
			Args: []string{fragmentName},
			GivenObjects: []client.Object{
				&testFragment,
			},
			ExpectOutput: `
name: test-fragment
namespace: accelerator-system
displayName: Test Fragment
git:
  interval: 2m0s
  ignore: .ignore
  url: http://www.test.com
  ref:
    branch: main
    tag: v1.0.0
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
importedBy:
  None
`,
		},
		{
			Name: "Get an fragment with empty values from context",
			Args: []string{fragmentName},
			GivenObjects: []client.Object{
				&testAcceleratorEmptyValues,
			},
			ExpectOutput: `
name: test-fragment
namespace: accelerator-system
displayName: Test Fragment
git:
  ignore: .ignore
  url: http://www.test.com
  ref:
    branch: main
    tag: v1.0.0
ready: true
options: []
artifact:
  message: test
  ready: true
  url: http://www.test.com
imports:
  None
importedBy:
  None
`,
		},
	}
	table.Run(t, scheme, FragmentGetCmd)
}
