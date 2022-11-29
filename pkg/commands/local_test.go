package commands

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
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
	Context("LocalGenerateCmd()", func() {
		When("Executes generate-from-local command with accelerator name", func() {
			It("Should send accelerator name and create project directory", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					r.ParseMultipartForm(100 << 20)
					Expect(r.FormValue("accelerator_name")).Should(Equal("acc"))
					zipWriter := zip.NewWriter(w)
					zipWriter.Close()
				}))
				defer ts.Close()
				generateCmd := LocalGenerateCmd()
				b := new(bytes.Buffer)
				generateCmd.SetOut(b)
				generateCmd.SetErr(b)
				generateCmd.SetArgs([]string{"--accelerator-name", "acc", "--server-url", ts.URL})
				defer os.RemoveAll("acc")
				generateCmd.Execute()
				out, err := ioutil.ReadAll(b)
				if err != nil {
					Fail("Error testing generate command")
				}

				Expect(string(out)).Should(Equal("generated project acc\n"))
			})
		})

		When("Executes generate-from-local command with local accelerator", func() {
			It("Should send accelerator bytes and create project directory", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					r.ParseMultipartForm(100 << 20)
					fileAcc, _, _ := r.FormFile("accelerator")

					gzr, err := gzip.NewReader(fileAcc)
					if err != nil {
						panic(err)
					}
					defer gzr.Close()
					tr := tar.NewReader(gzr)
					var header *tar.Header
					names := []string{}
					for header, err = tr.Next(); err == nil; header, err = tr.Next() {
						if header.Typeflag == tar.TypeReg {
							names = append(names, header.Name)
						}
					}

					Expect(names).Should(ContainElements("accelerator.yaml", "inner/foo.txt"))

					Expect(fileAcc).ShouldNot(BeNil())
					zipWriter := zip.NewWriter(w)
					zipWriter.Close()
				}))
				defer ts.Close()
				generateCmd := LocalGenerateCmd()
				b := new(bytes.Buffer)
				generateCmd.SetOut(b)
				generateCmd.SetErr(b)
				generateCmd.SetArgs([]string{"--accelerator-path", "acc=./testdata/test-acc", "--server-url", ts.URL})
				defer os.RemoveAll("acc")
				generateCmd.Execute()
				out, err := ioutil.ReadAll(b)
				if err != nil {
					Fail("Error testing generate command")
				}

				Expect(string(out)).Should(Equal("generated project acc\n"))
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
				defer ts.Close()
				generateCmd := LocalGenerateCmd()
				b := new(bytes.Buffer)
				generateCmd.SetOut(b)
				generateCmd.SetErr(b)
				generateCmd.SetArgs([]string{"--accelerator-name", "acc", "--fragment-names", "f1",
					"--fragment-names", "f2", "--fragment-paths", "f3=testdata/test-acc",
					"--fragment-paths", "f4=testdata/test-acc", "--server-url", ts.URL})
				defer os.RemoveAll("acc")
				generateCmd.Execute()
				out, err := ioutil.ReadAll(b)
				if err != nil {
					Fail("Error testing generate command")
				}

				Expect(string(out)).Should(Equal("generated project acc\n"))
			})
		})

		When("Executes generate-from-local command with output directory", func() {
			It("Should send accelerator name and create project directory", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					r.ParseMultipartForm(100 << 20)
					Expect(r.FormValue("accelerator_name")).Should(Equal("acc"))
					zipWriter := zip.NewWriter(w)
					zipWriter.Create("acc/")
					fileWriter, _ := zipWriter.Create("acc/test")
					fileWriter.Write([]byte("hello"))
					zipWriter.Close()
				}))
				defer ts.Close()
				generateCmd := LocalGenerateCmd()
				b := new(bytes.Buffer)
				generateCmd.SetOut(b)
				generateCmd.SetErr(b)
				generateCmd.SetArgs([]string{"--accelerator-name", "acc", "--output-dir", "output", "--server-url", ts.URL})
				defer os.RemoveAll("output")
				generateCmd.Execute()
				out, err := ioutil.ReadAll(b)
				if err != nil {
					Fail("Error testing generate command")
				}

				Expect(string(out)).Should(Equal("generated project acc\n"))
				Expect("output/test").Should(BeARegularFile())
			})
		})

		When("Executes generate-from-local command with --force", func() {
			It("Should send accelerator name and overwrite output directory", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					r.ParseMultipartForm(100 << 20)
					Expect(r.FormValue("accelerator_name")).Should(Equal("acc"))
					zipWriter := zip.NewWriter(w)
					zipWriter.Create("acc/")
					fileWriter, _ := zipWriter.Create("acc/test")
					fileWriter.Write([]byte("hello"))
					zipWriter.Close()
				}))
				defer ts.Close()
				os.MkdirAll("existing-dir/subpath", 0755)
				defer os.RemoveAll("existing-dir")
				generateCmd := LocalGenerateCmd()
				b := new(bytes.Buffer)
				generateCmd.SetOut(b)
				generateCmd.SetErr(b)
				generateCmd.SetArgs([]string{"--accelerator-name", "acc", "--output-dir", "existing-dir", "--force", "--server-url", ts.URL})
				generateCmd.Execute()
				out, err := ioutil.ReadAll(b)
				if err != nil {
					Fail("Error testing generate command")
				}

				Expect(string(out)).Should(Equal("generated project acc\n"))
				Expect("existing-dir/subpath").ShouldNot(BeAnExistingFile())
				Expect("existing-dir/test").Should(BeARegularFile())
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
				defer ts.Close()
				os.MkdirAll("existing-dir/subpath", 0755)
				defer os.RemoveAll("existing-dir")
				generateCmd := LocalGenerateCmd()
				generateCmd.SetArgs([]string{"--accelerator-name", "existing-dir", "--server-url", ts.URL})
				err := generateCmd.Execute()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).Should(ContainSubstring("path existing-dir is not empty, use --force to overwrite"))
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
				defer ets500.Close()
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
				ets503 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusServiceUnavailable)
				}))
				defer ets503.Close()
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
				ets404 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}))
				defer ets404.Close()
				generateCmd := LocalGenerateCmd()
				generateCmd.SetArgs([]string{"--accelerator-name", "test-missing", "--server-url", ets404.URL})
				err := generateCmd.Execute()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).Should(Equal("one of the accelerators or fragments was not found\n"))
			})
		})
	})
})
