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
	os.Setenv("ACC_SERVER_URL", ts.URL)
	generateCmd := GenerateCmd()
	b := new(bytes.Buffer)
	generateCmd.SetOut(b)
	generateCmd.SetErr(b)
	Context("Generate()", func() {
		When("Executes generate command", func() {
			It("Should return zip file", func() {
				generateCmd.SetArgs([]string{"test-acc"})
				generateCmd.Execute()
				out, err := ioutil.ReadAll(b)
				if err != nil {
					Fail("Error testing generate command")
				}

				Expect(string(out)).Should(Equal("zip file test-acc.zip created\n"))
				os.Remove("test-acc.zip")
			})
		})
	})
})
