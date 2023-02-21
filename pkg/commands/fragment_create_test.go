package commands

import (
	"testing"
	"time"

	"github.com/fluxcd/pkg/apis/meta"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	clitesting "github.com/vmware-tanzu/apps-cli-plugin/pkg/cli-runtime/testing"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestFragmentCreateCommand(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = acceleratorv1alpha1.AddToScheme(scheme)

	fragmentName := "test-fragment"
	gitRepoUrl := "https://www.test.com"
	gitRepoSubPath := "test-path"
	noGitBranch := "main"
	noGitTag := ""
	gitBranch := "main"
	gitTag := "v0.0.1"
	namespace := "accelerator-system"
	interval := "2m"
	secretRef := "mysecret"
	expectedDuration, _ := time.ParseDuration(interval)

	table := clitesting.CommandTestSuite{
		{
			Name:        "Missing args",
			Args:        []string{"--git-repository", gitRepoUrl},
			ShouldError: true,
		},
		{
			Name: "Error creating accelerator fragment",
			Args: []string{fragmentName, "--git-repository", gitRepoUrl},
			WithReactors: []clitesting.ReactionFunc{
				clitesting.InduceFailure("create", "Fragment"),
			},
			ExpectCreates: []client.Object{
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      fragmentName,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta1.GitRepositoryRef{
								Branch: noGitBranch,
								Tag:    noGitTag,
							},
						},
					},
				},
			},
			ExpectOutput: "Error creating accelerator fragment test-fragment\n",
			ShouldError:  true,
		},
		{
			Name: "Create Fragment with GitRepository",
			Args: []string{fragmentName, "--git-repository", gitRepoUrl},
			ExpectCreates: []client.Object{
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      fragmentName,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta1.GitRepositoryRef{
								Branch: noGitBranch,
								Tag:    noGitTag,
							},
						},
					},
				},
			},
			ExpectOutput: "created accelerator fragment test-fragment in namespace accelerator-system\n",
		},
		{
			Name: "Create Fragment with Branch and Tag and Secret ref",
			Args: []string{fragmentName,
				"--git-repository", gitRepoUrl,
				"--git-branch", gitBranch,
				"--git-tag", gitTag,
				"--git-sub-path", gitRepoSubPath,
				"--interval", interval,
				"--secret-ref", secretRef,
			},
			ExpectCreates: []client.Object{
				&acceleratorv1alpha1.Fragment{
					ObjectMeta: v1.ObjectMeta{
						Namespace: namespace,
						Name:      fragmentName,
					},
					Spec: acceleratorv1alpha1.FragmentSpec{
						Git: &acceleratorv1alpha1.Git{
							URL: gitRepoUrl,
							Reference: &v1beta1.GitRepositoryRef{
								Branch: gitBranch,
								Tag:    gitTag,
							},
							SubPath: &gitRepoSubPath,
							SecretRef: &meta.LocalObjectReference{
								Name: secretRef,
							},
							Interval: &v1.Duration{
								Duration: expectedDuration,
							},
						},
					},
				},
			},
			ExpectOutput: "created accelerator fragment test-fragment in namespace accelerator-system\n",
		},
	}
	table.Run(t, scheme, FragmentCreateCmd)
}
