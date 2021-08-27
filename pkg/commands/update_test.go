package commands

import (
	"bytes"
	"context"
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

var _ = Describe("command update", func() {
	acceleratorName := "test"
	invalidAcceleratorName := "non-existent"
	updatedBranch := "another"
	acc := &acceleratorv1alpha1.Accelerator{
		TypeMeta: v1.TypeMeta{
			APIVersion: "accelerator.tanzu.vmware.com/v1alpha1",
			Kind:       "Accelerator",
		},
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
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
	b := bytes.NewBufferString("")
	updateCmd := UpdateCmd(clientset.FakeAcceleratorV1Alpha1())
	updateCmd.SetOut(b)
	updateCmd.SetErr(b)
	Context("update()", func() {
		When("updates existing accelerator", func() {
			It("Should update the accelerator", func() {
				updateCmd.SetArgs([]string{acceleratorName, "--git-branch", updatedBranch})
				updateCmd.Execute()
				out, _ := ioutil.ReadAll(b)
				Expect(string(out)).Should(Equal(fmt.Sprintf("accelerator %s updated successfully", acceleratorName)))
			})
		})

		When("updates non-existing accelerator", func() {
			It("Should throw error", func() {
				updateCmd.SetArgs([]string{invalidAcceleratorName, "--git-branch", updatedBranch})
				updateCmd.Execute()
				out, _ := ioutil.ReadAll(b)
				Expect(strings.HasPrefix(string(out), fmt.Sprintf("accelerator %s not found", invalidAcceleratorName))).Should(BeTrue())
			})
		})

		When("adds reconcile flag", func() {
			It("Should add the requestedAt annotation", func() {
				updateCmd.SetArgs([]string{acceleratorName, "--reconcile"})
				updateCmd.Execute()
				reconciledAcc, _ := clientset.FakeAcceleratorV1Alpha1().Accelerators("default").Get(context.Background(), "test", v1.GetOptions{})
				out, _ := ioutil.ReadAll(b)
				Expect(strings.HasPrefix(string(out), fmt.Sprintf("accelerator %s updated successfully", acceleratorName))).Should(BeTrue())
				Expect(reconciledAcc.ObjectMeta.Annotations["requestedAt"]).ShouldNot(BeNil())
			})
		})
	})

})
