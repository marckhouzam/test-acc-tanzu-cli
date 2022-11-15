package commands

import (
	"archive/zip"
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
	Context("LocalGenerateCmd()", func() {
		When("Executes generate-from-local command with accelerator name", func() {
			It("Should send accelerator name and create project directory", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					r.ParseMultipartForm(100 << 20)
					Expect(r.FormValue("accelerator_name")).Should(Equal("acc"))
					zipWriter := zip.NewWriter(w)
					zipWriter.Close()
				}))
				generateCmd := LocalGenerateCmd()
				b := new(bytes.Buffer)
				generateCmd.SetOut(b)
				generateCmd.SetErr(b)
				generateCmd.SetArgs([]string{"--accelerator-name", "acc", "--server-url", ts.URL})
				generateCmd.Execute()
				out, err := ioutil.ReadAll(b)
				if err != nil {
					Fail("Error testing generate command")
				}

				Expect(string(out)).Should(Equal("generated project acc\n"))
				os.RemoveAll("acc")
			})
		})

		When("Executes generate-from-local command with local accelerator", func() {
			It("Should send accelerator bytes and create project directory", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					r.ParseMultipartForm(100 << 20)
					fileAcc, handlerAcc, _ := r.FormFile("accelerator")
					Expect(fileAcc).ShouldNot(BeNil())
					Expect(handlerAcc.Filename).Should(Equal("test-acc"))
					zipWriter := zip.NewWriter(w)
					zipWriter.Close()
				}))
				generateCmd := LocalGenerateCmd()
				b := new(bytes.Buffer)
				generateCmd.SetOut(b)
				generateCmd.SetErr(b)
				generateCmd.SetArgs([]string{"--accelerator-path", "acc=testdata/test-acc", "--server-url", ts.URL})
				generateCmd.Execute()
				out, err := ioutil.ReadAll(b)
				if err != nil {
					Fail("Error testing generate command")
				}

				Expect(string(out)).Should(Equal("generated project acc\n"))
				os.RemoveAll("acc")
			})
		})

		When("Executes generate-from-local command with a combination of fragments", func() {
			It("Should send accelerator fragments", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					r.ParseMultipartForm(100 << 20)
					Expect(r.Form["fragment_names"]).Should(ContainElements("f1", "f2"))
					f3File, f3Handler, _ := r.FormFile("fragment_f3")
					Expect(f3File).ShouldNot(BeNil())
					Expect(f3Handler.Filename).Should(Equal("test-acc"))
					f4File, f4Handler, _ := r.FormFile("fragment_f4")
					Expect(f4File).ShouldNot(BeNil())
					Expect(f4Handler.Filename).Should(Equal("test-acc"))
					zipWriter := zip.NewWriter(w)
					zipWriter.Close()
				}))
				generateCmd := LocalGenerateCmd()
				b := new(bytes.Buffer)
				generateCmd.SetOut(b)
				generateCmd.SetErr(b)
				generateCmd.SetArgs([]string{"--accelerator-name", "acc", "--fragment-names", "f1",
					"--fragment-names", "f2", "--fragment-paths", "f3=testdata/test-acc",
					"--fragment-paths", "f4=testdata/test-acc", "--server-url", ts.URL})
				generateCmd.Execute()
				out, err := ioutil.ReadAll(b)
				if err != nil {
					Fail("Error testing generate command")
				}

				Expect(string(out)).Should(Equal("generated project acc\n"))
				os.RemoveAll("acc")
			})
		})

		When("Executes generate-from-local command with invalid directory", func() {
			It("Should output error message", func() {
				generateCmd := LocalGenerateCmd()
				generateCmd.SetArgs([]string{"--accelerator-path", "acc=invalid"})
				err := generateCmd.Execute()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).Should(Equal("cannot find directory invalid"))
			})
		})

		When("Executes generate-from-local command and project directory already exists", func() {
			It("Should output error message", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					zipWriter := zip.NewWriter(w)
					zipWriter.Create("existing-dir/")
					zipWriter.Close()
				}))
				os.Mkdir("existing-dir", 0755)
				generateCmd := LocalGenerateCmd()
				generateCmd.SetArgs([]string{"--accelerator-name", "existing-dir", "--server-url", ts.URL})
				err := generateCmd.Execute()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).Should(ContainSubstring("existing-dir already exists, use --force to overwrite"))
				os.RemoveAll("existing-dir")
			})
		})

		When("Executes generate-from-local command without accelerator", func() {
			It("Should output error message", func() {
				generateCmd := LocalGenerateCmd()
				generateCmd.SetArgs([]string{})
				err := generateCmd.Execute()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).Should(Equal("no accelerator, you must provide --accelerator-name or --accelerator-path"))
			})
		})

		When("Executes generate and response it's a 500", func() {
			It("Should output error message", func() {
				generateCmd := LocalGenerateCmd()
				generateCmd.SetArgs([]string{"--accelerator-name", "test-500", "--server-url", ets500.URL})
				err := generateCmd.Execute()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).Should(Equal("there was an error generating the accelerator, " +
					"the server response was: \"Invalid options provided. Expected but not present: " +
					"[one, second]\"\n"))
			})
		})

		When("Executes generate and response is a 503", func() {
			It("Should output error message", func() {
				generateCmd := LocalGenerateCmd()
				generateCmd.SetArgs([]string{"--accelerator-name", "test-503", "--server-url", ets503.URL})
				err := generateCmd.Execute()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).Should(Equal("there was an error generating the accelerator, " +
					"the server response code was: \"503\"\n"))
			})
		})

		When("Executes generate and response is a 404", func() {
			It("Should output error message", func() {
				generateCmd := LocalGenerateCmd()
				generateCmd.SetArgs([]string{"--accelerator-name", "test-missing", "--server-url", ets404.URL})
				err := generateCmd.Execute()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).Should(Equal("one of the accelerators or fragments was not found\n"))
			})
		})
	})
})
