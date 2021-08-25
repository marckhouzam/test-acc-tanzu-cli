package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/tabwriter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/acc-controller/api/clientset/fake"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	fluxcdv1beta1 "github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("command list", func() {
	var acc []runtime.Object
	for i := 0; i < 3; i++ {
		acc = append(acc, &acceleratorv1alpha1.Accelerator{
			TypeMeta: v1.TypeMeta{
				APIVersion: "accelerator.tanzu.vmware.com/v1alpha1",
				Kind:       "Accelerator",
			},
			ObjectMeta: v1.ObjectMeta{
				Namespace: "default",
				Name:      fmt.Sprintf("test-%d", i),
			},
			Spec: acceleratorv1alpha1.AcceleratorSpec{
				Git: acceleratorv1alpha1.Git{
					URL: "http://www.test.com",
					Reference: &fluxcdv1beta1.GitRepositoryRef{
						Branch: "main",
					},
				},
			},
		})
	}
	clientset := fake.NewSimpleClientset(acc...)
	listCmd := ListCmd(clientset.FakeAcceleratorV1Alpha1())
	b := bytes.NewBufferString("")
	listCmd.SetOut(b)
	listCmd.SetErr(b)
	Context("list()", func() {
		When("looks for existing accelerators", func() {
			It("Should return the list of accelerators", func() {
				tempBuf := bytes.NewBufferString("")
				w := new(tabwriter.Writer)
				w.Init(tempBuf, 0, 8, 3, ' ', 0)
				fmt.Fprintln(w, "NAME\tGIT REPOSITORY\tBRANCH\tTAG")
				for _, obj := range acc {
					accelerator := obj.(*acceleratorv1alpha1.Accelerator)
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", accelerator.Name, accelerator.Spec.Git.URL, accelerator.Spec.Git.Reference.Branch, accelerator.Spec.Git.Reference.Tag)
				}
				w.Flush()
				listCmd.Execute()
				out, _ := ioutil.ReadAll(b)
				expected, _ := ioutil.ReadAll(tempBuf)
				Expect(string(out)).Should(Equal(string(expected)))
			})
		})

	})

})
