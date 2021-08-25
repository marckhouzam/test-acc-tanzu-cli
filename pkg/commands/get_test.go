package commands

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"text/tabwriter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/acc-controller/api/clientset/fake"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	fluxcdv1beta1 "github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("command get", func() {
	acceleratorName := "test"
	namespace := "default"
	invalidAcceleratorName := "non-existent"
	acc := &acceleratorv1alpha1.Accelerator{
		TypeMeta: v1.TypeMeta{
			APIVersion: "accelerator.tanzu.vmware.com/v1alpha1",
			Kind:       "Accelerator",
		},
		ObjectMeta: v1.ObjectMeta{
			Namespace: namespace,
			Name:      acceleratorName,
		},
		Spec: acceleratorv1alpha1.AcceleratorSpec{
			Git: acceleratorv1alpha1.Git{
				URL: "http://www.test.com",
				Reference: &fluxcdv1beta1.GitRepositoryRef{
					Branch: "main",
				},
			},
		},
	}
	clientset := fake.NewSimpleClientset(acc)
	getCmd := GetCmd(clientset.FakeAcceleratorV1Alpha1())
	b := new(bytes.Buffer)
	getCmd.SetOut(b)
	getCmd.SetErr(b)
	Context("get()", func() {
		When("looks for existing accelerator", func() {
			It("Should return the accelerator", func() {
				tempBuf := bytes.NewBufferString("")
				w := new(tabwriter.Writer)
				w.Init(tempBuf, 0, 8, 3, ' ', 0)
				fmt.Fprintln(w, "NAME\tGIT REPOSITORY\tBRANCH\tTAG")
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", acc.Name, acc.Spec.Git.URL, acc.Spec.Git.Reference.Branch, acc.Spec.Reference.Tag)
				w.Flush()
				getCmd.SetArgs([]string{acceleratorName})
				getCmd.ExecuteContext(context.Background())
				out, err := ioutil.ReadAll(b)
				if err != nil {
					Fail("Error testing GET command")
				}
				expected, _ := ioutil.ReadAll(tempBuf)
				Expect(string(out)).Should(Equal(string(expected)))
			})

			It("Should throw error for non existent accelerator", func() {
				expectErrorMsg := fmt.Sprintf("Error getting accelerator %s", invalidAcceleratorName)
				getCmd.SetArgs([]string{invalidAcceleratorName})
				err := getCmd.ExecuteContext(context.Background())
				Expect(err).ShouldNot(BeNil())
				out, _ := ioutil.ReadAll(b)
				Expect(strings.HasPrefix(string(out), expectErrorMsg)).Should(BeTrue())
			})
		})

	})

})
