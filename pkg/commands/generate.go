/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

type UiServerBody struct {
	Accelerator string                 `json:"accelerator"`
	Options     map[string]interface{} `json:"options"`
}

type OptionsProjectName struct {
	ProjectName string `json:"projectName"`
}

type AcceleratorName struct {
	Name string `json:"name"`
}
type Accelerators struct {
	Accelerators []AcceleratorName `json:"accelerators"`
}

type UiAcceleratorList struct {
	Embedded Accelerators `json:"_embedded"`
}

type UiErrorResponse struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func GenerateCmd() *cobra.Command {
	var uiServer string
	var accServerUrl string
	var optionsString string
	var filepath string
	var outputDir string
	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate project from accelerator",
		Long: `Generate a project from an accelerator using provided options and download project artifacts as a ZIP file.

Generation options are provided as a JSON string and should match the metadata options that are specified for the
accelerator used for the generation. The options can include "projectName" which defaults to the name of the accelerator.
This "projectName" will be used as the name of the generated ZIP file.

You can see the available options by using the "tanzu accelerator list <accelerator-name>" command.

Here is an example of an options JSON string that specifies the "projectName" and an "includeKubernetes" boolean flag:

    --options '{"projectName":"test", "includeKubernetes": true}'

You can also provide a file that specifies the JSON string using the --options-file flag.

The generate command needs access to the Application Accelerator server. You can specify the --server-url flag or set
an ACC_SERVER_URL environment variable. If you specify the --server-url flag it will override the ACC_SERVER_URL
environmnet variable if it is set.
`,
		ValidArgsFunction: SuggestAcceleratorNamesFromUiServer(context.Background()),
		Example:           "tanzu accelerator generate <accelerator-name> --options '{\"projectName\":\"test\"}'",
		RunE: func(cmd *cobra.Command, args []string) error {
			if optionsString == "" {
				optionsString = "{\"projectName\": \"" + args[0] + "\"}"
			}
			if !strings.Contains(optionsString, "projectName") {
				optionsString = "{\"projectName\": \"" + args[0] + "\"," + optionsString[1:]
			}
			if !strings.HasSuffix(outputDir, "/") && outputDir != "" {
				outputDir += "/"
			}
			var projectName OptionsProjectName
			json.Unmarshal([]byte(optionsString), &projectName)
			client := &http.Client{}
			var options map[string]interface{}
			if filepath != "" {
				fileBytes, err := ioutil.ReadFile(filepath)
				if err != nil {
					return err
				}
				optionsString = string(fileBytes)
			}
			err := json.Unmarshal([]byte(optionsString), &options)
			if err != nil {
				return err
			}
			uiServerBody := UiServerBody{
				Accelerator: args[0],
				Options:     options,
			}
			JsonProxyBodyBytes, err := json.Marshal(uiServerBody)
			if err != nil {
				return errors.New("error marshalling request body")
			}
			serverUrl := accServerUrl
			if uiServer != "" {
				serverUrl = uiServer
			}
			if serverUrl == "" {
				return errors.New("no server URL provided, you must provide --server-url option or set ACC_SERVER_URL environment variable")
			}
			proxyRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/api/accelerators/zip?name=%s", serverUrl, args[0]), bytes.NewReader(JsonProxyBodyBytes))
			proxyRequest.Header.Add("Content-Type", "application/json")
			if err != nil {
				return errors.New(fmt.Sprintf("error creating request for %s", serverUrl))
			}
			resp, err := client.Do(proxyRequest)
			if err != nil {
				return errors.New(fmt.Sprintf("error invoking %s", serverUrl))
			}

			if resp.StatusCode >= 400 {
				var errorMsg string
				if resp.StatusCode == http.StatusNotFound {
					errorMsg = fmt.Sprintf("accelerator %s not found\n", args[0])
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
			zipfile := outputDir + projectName.ProjectName + ".zip"
			err = ioutil.WriteFile(zipfile, body, 0644)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "zip file %s created\n", zipfile)

			return nil
		},
	}
	generateCmd.Flags().StringVar(&optionsString, "options", "", "options JSON string")
	generateCmd.Flags().StringVar(&filepath, "options-file", "", "path to file containing options JSON string")
	generateCmd.Flags().StringVar(&outputDir, "output-dir", "", "directory that the zip file will be written to")
	generateCmd.Flags().StringVar(&uiServer, "server-url", "", "the URL for the Application Accelerator server")
	accServerUrl = EnvVar("ACC_SERVER_URL", "")
	return generateCmd
}
