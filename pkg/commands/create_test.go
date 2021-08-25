package commands

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/acc-controller/api/clientset/fake"
)

var _ = Describe("command get", func() {
	acceleratorName := "test"
	namespace := "default"
	clientset := fake.NewSimpleClientset()
	createCmd := CreateCmd(clientset.FakeAcceleratorV1Alpha1())
	b := new(bytes.Buffer)
	createCmd.SetOut(b)
	createCmd.SetErr(b)
	Context("create()", func() {
		When("calling create command", func() {
			It("Should create the accelerator", func() {
				createCmd.SetArgs([]string{acceleratorName, "--git-repository", "https://www.test.com", "--git-branch", "main"})
				createCmd.ExecuteContext(context.Background())
				out, err := ioutil.ReadAll(b)
				if err != nil {
					Fail("Error testing GET command")
				}
				Expect(string(out)).Should(Equal(fmt.Sprintf("created accelerator %s in namespace %s\n", acceleratorName, namespace)))
			})
		})

	})

})
