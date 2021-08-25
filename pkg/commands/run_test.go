package commands

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("command run", func() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Test String")
	}))
	runCmd := RunCmd(ts.URL)
	b := new(bytes.Buffer)
	runCmd.SetOut(b)
	runCmd.SetErr(b)
	Context("Run()", func() {
		When("Executes run command", func() {
			It("Should return zip file", func() {
				runCmd.SetArgs([]string{"test-acc"})
				runCmd.Execute()
				out, err := ioutil.ReadAll(b)
				if err != nil {
					Fail("Error testing GET command")
				}

				Expect(string(out)).Should(Equal("zip file test-acc.zip created"))
				os.Remove("test-acc.zip")
			})
		})
	})
})
