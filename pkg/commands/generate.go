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
	var optionsString string
	var filepath string
	var outputDir string
	var generateCmd = &cobra.Command{
		Use:               "generate",
		Short:             "Generate project from accelerator",
		Long:              `Generate a project from an accelerator and download project artifacts as a ZIP file`,
		ValidArgsFunction: SuggestAcceleratorNamesFromUiServer(context.Background()),
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
				return errors.New("error marshalling proxy body")
			}
			proxyRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/api/accelerators/zip?name=%s", uiServer, args[0]), bytes.NewReader(JsonProxyBodyBytes))
			proxyRequest.Header.Add("Content-Type", "application/json")
			if err != nil {
				return errors.New("error creating proxy request")
			}
			resp, err := client.Do(proxyRequest)
			if err != nil {
				return errors.New("error proxying engine invocation")
			}

			if resp.StatusCode >= 400 {
				var errorMsg string
				if resp.StatusCode == http.StatusNotFound {
					errorMsg = fmt.Sprintf("the accelerator was not found\n")
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
	defaultUiServerUrl := EnvVar("ACC_SERVER_URL", "http://localhost:8877")
	generateCmd.Flags().StringVar(&optionsString, "options", "", "Enter options string")
	generateCmd.Flags().StringVar(&filepath, "options-file", "", "Enter file path with json body")
	generateCmd.Flags().StringVar(&outputDir, "output-dir", "", "Directory where the zip file should be written")
	generateCmd.Flags().StringVar(&uiServer, "server-url", defaultUiServerUrl, "The App Accelerator server URL, this will override ACC_SERVER_URL env variable")
	return generateCmd
}
