package commands

import (
	"bytes"
	"encoding/json"
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

	ets500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorResponse := UiErrorResponse{
			Status: 500,
			Title:  "MissingOptions",
			Detail: "Invalid options provided. Expected but not present: [one, second]",
		}
		json, _ := json.Marshal(errorResponse)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, string(json))
	}))
	ets503 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	ets404 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
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

		When("Executes generate and response it's a 500", func() {
			It("Should output error messasge", func() {
				generateCmd.SetArgs([]string{"test-acc", "--server-url", ets500.URL})
				err := generateCmd.Execute()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).Should(Equal("there was an error generating the accelerator, " +
					"the server response was: \"Invalid options provided. Expected but not present: " +
					"[one, second]\"\n"))
			})
		})

		When("Executes generate and response is a 503", func() {
			It("Should output error messasge", func() {
				generateCmd.SetArgs([]string{"test-acc", "--server-url", ets503.URL})
				err := generateCmd.Execute()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).Should(Equal("there was an error generating the accelerator, " +
					"the server response code was: \"503\"\n"))
			})
		})

		When("Executes generate and response is a 404", func() {
			It("Should output error messasge", func() {
				generateCmd.SetArgs([]string{"test-missing", "--server-url", ets404.URL})
				err := generateCmd.Execute()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).Should(Equal("accelerator test-missing not found\n"))
			})
		})
	})
})
