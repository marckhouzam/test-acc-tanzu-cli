/*
Copyright 2021 VMware, Inc. All Rights Reserved.
*/
package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
		Run: func(cmd *cobra.Command, args []string) {
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
					log.Fatal(err.Error())
				}
				optionsString = string(fileBytes)
			}
			err := json.Unmarshal([]byte(optionsString), &options)
			if err != nil {
				log.Fatal(err.Error())
			}
			uiServerBody := UiServerBody{
				Accelerator: args[0],
				Options:     options,
			}
			JsonProxyBodyBytes, err := json.Marshal(uiServerBody)
			if err != nil {
				log.Fatal(err, "Error marshalling proxy body")
				return
			}
			proxyRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/api/accelerators/zip?name=%s", uiServer, args[0]), bytes.NewReader(JsonProxyBodyBytes))
			proxyRequest.Header.Add("Content-Type", "application/json")
			if err != nil {
				log.Fatal(err, "Error creating proxy request")
				return
			}
			resp, err := client.Do(proxyRequest)
			if err != nil {
				log.Fatal(err, "Error proxying engine invocation")
				return
			}
			body, _ := ioutil.ReadAll(resp.Body)
			zipfile := outputDir + projectName.ProjectName + ".zip"
			err = ioutil.WriteFile(zipfile, body, 0644)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "zip file %s created\n", zipfile)
		},
	}
	defaultUiServerUrl := EnvVar("ACC_SERVER_URL", "http://localhost:8877")
	generateCmd.Flags().StringVar(&optionsString, "options", "", "Enter options string")
	generateCmd.Flags().StringVar(&filepath, "options-file", "", "Enter file path with json body")
	generateCmd.Flags().StringVar(&outputDir, "output-dir", "", "Directory where the zip file should be written")
	generateCmd.Flags().StringVar(&uiServer, "server-url", defaultUiServerUrl, "The App Accelerator server URL, this will override ACC_SERVER_URL env variable")
	return generateCmd
}