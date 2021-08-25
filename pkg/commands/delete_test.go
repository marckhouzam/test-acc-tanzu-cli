package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/acc-controller/api/clientset/fake"
	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	fluxcdv1beta1 "github.com/pivotal/acc-controller/fluxcd/api/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("command delete", func() {
	Context("delete()", func() {
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
		DeleteCmd := DeleteCmd(clientset.FakeAcceleratorV1Alpha1())
		b := bytes.NewBufferString("")
		DeleteCmd.SetOut(b)
		DeleteCmd.SetErr(b)
		When("deletes exisitng accelerator", func() {
			It("Should delete accelerator without errors", func() {
				DeleteCmd.SetArgs([]string{acceleratorName})
				DeleteCmd.Execute()
				out, _ := ioutil.ReadAll(b)
				Expect(string(out)).Should(Equal(fmt.Sprintf("deleted accelerator %s in namespace %s\n", acceleratorName, namespace)))
			})
			It("Should throw error for non existent accelerator", func() {
				DeleteCmd.SetArgs([]string{invalidAcceleratorName})
				err := DeleteCmd.Execute()
				Expect(err).ShouldNot(BeNil())
				out, _ := ioutil.ReadAll(b)
				Expect(strings.HasPrefix(string(out), fmt.Sprintf("There was a problem trying to delete accelerator %s", invalidAcceleratorName))).Should(BeTrue())
			})
		})

	})

})
