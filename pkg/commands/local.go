/*
Copyright 2022-2023 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/denormal/go-gitignore"
	"github.com/spf13/cobra"
)

// kvPair represents a single key-value pair
type kvPair struct {
	key   string
	value string
}

func (p *kvPair) isEmpty() bool {
	return len(p.key) == 0 || len(p.value) == 0
}

// pairValue implements the cobra Value interface, making it usable as a flag
type pairValue struct {
	pair *kvPair
}

func newPairValue(val kvPair, p *kvPair) *pairValue {
	pairV := new(pairValue)
	pairV.pair = p
	*pairV.pair = val
	return pairV
}

func (s *pairValue) Type() string {
	return `"key=value" pair`
}

func (s *pairValue) String() string {
	if s.pair == nil || s.pair.isEmpty() {
		return ""
	}
	return s.pair.key + "=" + s.pair.value
}

func (s *pairValue) Set(val string) error {
	kv := strings.SplitN(val, "=", 2)
	if len(kv) != 2 {
		return fmt.Errorf("%s must be formatted as key=value", val)
	}
	*s.pair = kvPair{kv[0], kv[1]}
	return nil
}

func LocalGenerateCmd() *cobra.Command {
	var uiServer string
	var accServerUrl string
	var optionsString string
	var optionsFilename string
	var acceleratorName string
	var outputDirectory string
	var localAccelerator kvPair
	var fragmentNames []string
	var localFragments map[string]string
	var forceOverwrite bool
	var localGenerateCommand = &cobra.Command{
		Use:   "generate-from-local",
		Short: "Generate project from a combination of registered and local artifacts",
		Long: `Generate a project from a combination of local files and registered accelerators/fragments using provided 
options and download project artifacts as a ZIP file.

Options values are provided as a JSON object and should match the declared options that are specified for the
accelerator used for the generation. The options can include "projectName" which defaults to the name of the accelerator.
This "projectName" will be used as the name of the generated ZIP file.

Here is an example of an options JSON string that specifies the "projectName" and an "includeKubernetes" boolean flag:

    --options '{"projectName":"test", "includeKubernetes": true}'

You can also provide a file that specifies the JSON string using the --options-file flag.

The generate-from-local command needs access to the Application Accelerator server. You can specify the --server-url flag or set
an ACC_SERVER_URL environment variable. If you specify the --server-url flag it will override the ACC_SERVER_URL
environment variable if it is set.
`,
		Example: `tanzu accelerator generate-from-local --accelerator-path java-rest=workspace/java-rest --fragment-paths java-version=workspace/version --fragment-names tap-workload --options '{"projectName":"test"}'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// build a form body
			requestBody := &bytes.Buffer{}
			bodyWriter := multipart.NewWriter(requestBody)

			var defaultProjectName string
			if !localAccelerator.isEmpty() {
				defaultProjectName = localAccelerator.key
				accFolderName := localAccelerator.value
				fileWriter, err := bodyWriter.CreateFormFile("accelerator", defaultProjectName+".tar.gz")

				err = tarToWriter(accFolderName, fileWriter)
				if err != nil {
					return err
				}
			} else if acceleratorName != "" {
				defaultProjectName = acceleratorName
				field, err := bodyWriter.CreateFormField("accelerator_name")
				if err != nil {
					return err
				}
				_, err = field.Write([]byte(acceleratorName))
				if err != nil {
					return err
				}
				if err != nil {
					return err
				}
			} else {
				return errors.New("no accelerator, you must provide --accelerator-name or --accelerator-path")
			}

			for _, fragmentName := range fragmentNames {
				field, err := bodyWriter.CreateFormField("fragment_names")
				if err != nil {
					return err
				}
				_, err = field.Write([]byte(fragmentName))
				if err != nil {
					return err
				}
			}

			for fragmentName, fragmentFolderName := range localFragments {
				fileWriter, err := bodyWriter.CreateFormFile("fragment_"+fragmentName, fragmentName+".tar.gz")
				err = tarToWriter(fragmentFolderName, fileWriter)
				if err != nil {
					return err
				}
			}

			if optionsFilename != "" {
				fileBytes, err := ioutil.ReadFile(optionsFilename)
				if err != nil {
					return err
				}
				optionsString = string(fileBytes)
			}

			var options map[string]interface{}
			if optionsFilename != "" {
				fileBytes, err := ioutil.ReadFile(optionsFilename)
				if err != nil {
					return err
				}
				optionsString = string(fileBytes)
			}
			err := json.Unmarshal([]byte(optionsString), &options)
			if err != nil {
				return errors.New("invalid options provided, must be valid JSON")
			}
			if _, found := options["projectName"]; !found {
				options["projectName"] = defaultProjectName
			}
			projectName := options["projectName"].(string)

			client := &http.Client{}

			optionsField, err := bodyWriter.CreateFormField("options")
			if err != nil {
				return err
			}
			err = json.NewEncoder(optionsField).Encode(options)
			if err != nil {
				return err
			}

			// Close the body writer
			bodyWriter.Close()

			serverUrl := accServerUrl
			if uiServer != "" {
				serverUrl = uiServer
			}
			if serverUrl == "" {
				return errors.New("no server URL provided, you must provide --server-url option or set ACC_SERVER_URL environment variable")
			}
			if !strings.HasPrefix(serverUrl, "http://") && !strings.HasPrefix(serverUrl, "https://") {
				return errors.New(fmt.Sprintf("error creating request for %s, the URL needs to include the protocol (\"http://\" or \"https://\")", serverUrl))
			}

			apiPrefix := DetermineApiServerPrefix(serverUrl)
			proxyRequest, _ := http.NewRequest("POST", fmt.Sprintf("%s/%s/accelerators/zip", serverUrl, apiPrefix), requestBody)
			proxyRequest.Header.Add("Content-Type", bodyWriter.FormDataContentType())
			resp, err := client.Do(proxyRequest)
			if err != nil {
				return err
			}

			if resp.StatusCode >= 300 {
				var errorMsg string
				if resp.StatusCode == http.StatusNotFound {
					errorMsg = fmt.Sprintf("one of the accelerators or fragments was not found\n")
				} else {
					var errorResponse UiErrorResponse
					body, _ := ioutil.ReadAll(resp.Body)
					json.Unmarshal(body, &errorResponse)
					if errorResponse.Detail > "" {
						errorMsg = fmt.Sprintf("there was an error generating the accelerator, the server response was: \"%s\"\n", errorResponse.Detail)
					} else {
						errorMsg = fmt.Sprintf("there was an error generating the accelerator, the server response code was: \"%v\"\n", resp.StatusCode)
					}
				}
				return fmt.Errorf(errorMsg)
			}

			body, _ := ioutil.ReadAll(resp.Body)
			zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
			if err != nil {
				return err
			}

			targetDirectory := outputDirectory
			if outputDirectory == "" {
				targetDirectory = projectName
			}

			if forceOverwrite {
				err := os.RemoveAll(targetDirectory)
				if err != nil {
					return errors.New(fmt.Sprintf("could not remove %s", targetDirectory))
				}
			} else {
				if _, err := os.Stat(targetDirectory); !errors.Is(err, os.ErrNotExist) {
					// directory exists
					if empty, _ := isEmpty(targetDirectory); !empty {
						return errors.New(fmt.Sprintf("path %s is not empty, use --force to overwrite", targetDirectory))
					}
				}
			}
			for _, f := range zipReader.File {
				err = extractFile(f, targetDirectory)
				if err != nil {
					return err
				}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "generated project %s\n", projectName)
			return nil
		},
	}
	localGenerateCommand.Flags().StringVar(&optionsString, "options", "{}", "options JSON string")
	localGenerateCommand.Flags().StringVar(&optionsFilename, "options-file", "", "path to file containing options JSON string")
	localGenerateCommand.Flags().StringVar(&uiServer, "server-url", "", "the URL for the Application Accelerator server")
	localGenerateCommand.Flags().StringVarP(&outputDirectory, "output-dir", "o", "", "the directory that the project will be created in (defaults to the project name)")
	localGenerateCommand.Flags().StringVar(&acceleratorName, "accelerator-name", "", "name of the registered accelerator to use")
	localGenerateCommand.Flags().StringSliceVar(&fragmentNames, "fragment-names", []string{}, "names of the registered fragments to use")
	localGenerateCommand.Flags().Var(newPairValue(kvPair{}, &localAccelerator), "accelerator-path", "key value pair of the name and path to the directory containing the accelerator")
	localGenerateCommand.Flags().StringToStringVar(&localFragments, "fragment-paths", map[string]string{}, "key value pairs of the name and path to the directory containing each fragment")
	localGenerateCommand.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "force clean and rewrite of output-dir")
	localGenerateCommand.MarkFlagsMutuallyExclusive("options", "options-file")
	localGenerateCommand.MarkFlagsMutuallyExclusive("accelerator-path", "accelerator-name")
	accServerUrl = EnvVar("ACC_SERVER_URL", "")
	return localGenerateCommand
}

func extractFile(f *zip.File, targetDirectory string) error {
	filePaths := strings.Split(f.Name, "/")[1:]
	path := filepath.Join(append([]string{targetDirectory}, filePaths...)...)

	if f.FileInfo().IsDir() {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return errors.New(fmt.Sprintf("could not create directory %s", path))
		}
	} else {
		// create directories to the file
		os.MkdirAll(filepath.Dir(path), 0755)

		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return errors.New("error creating subdirectories in generated project")
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return errors.New(fmt.Sprintf("could not open file %s", f.Name))
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return errors.New(fmt.Sprintf("could not open file %s", f.Name))
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return nil
}

func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

// tarToWriter takes a source and a writer and walks sourceDir writing each file
// found to the tar writer
func tarToWriter(sourceDir string, writer io.Writer) error {

	// ensure the sourceDir actually exists before trying to tar it
	if _, err := os.Stat(sourceDir); err != nil {
		return fmt.Errorf("cannot find directory %v", sourceDir)
	}

	gzw := gzip.NewWriter(writer)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	cleanSourceDir := filepath.Clean(sourceDir)

	ignore, err := gitignore.NewRepository(cleanSourceDir)
	if err != nil {
		return err
	}

	// walk path
	return filepath.Walk(sourceDir, func(file string, fi os.FileInfo, err error) error {

		// return on any error
		if err != nil {
			return err
		}

		// exclude .git directory and its contents
		if fi.IsDir() && fi.Name() == ".git" {
			return filepath.SkipDir
		}

		// exclude directories and files in .gitignore
		// don't call Ignore for root path (see https://github.com/denormal/go-gitignore/pull/4)
		if file != sourceDir {
			if match := ignore.Match(file); match != nil && match.Ignore() {
				if fi.IsDir() {
					return filepath.SkipDir
				} else {
					return nil
				}
			}
		}

		// return on non-regular files
		if !fi.Mode().IsRegular() {
			return nil
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}

		// update the name to correctly reflect the desired destination when untaring, accounting for Windows paths
		// which use "\" instead of "/"
		fileDestination := strings.TrimPrefix(filepath.Clean(file), cleanSourceDir+string(filepath.Separator))
		header.Name = filepath.ToSlash(fileDestination)

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		// manually close here after each file operation; deferring would cause each file close
		// to wait until all operations have completed.
		f.Close()

		return nil
	})
}
