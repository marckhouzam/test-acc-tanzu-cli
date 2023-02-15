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
	"os/user"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type UiServerBody struct {
	Accelerator string                 `json:"accelerator"`
	Options     map[string]interface{} `json:"options"`
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
	var filename string
	var outputDir string
	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate project from accelerator",
		Long: `Generate a project from an accelerator using provided options and download project artifacts as a ZIP file.

Generation options are provided as a JSON string and should match the metadata options that are specified for the
accelerator used for the generation. The options can include "projectName" which defaults to the name of the accelerator.
This "projectName" will be used as the name of the generated ZIP file.

You can see the available options by using the "tanzu accelerator get <accelerator-name>" command.

Here is an example of an options JSON string that specifies the "projectName" and an "includeKubernetes" boolean flag:

    --options '{"projectName":"test", "includeKubernetes": true}'

You can also provide a file that specifies the JSON string using the --options-file flag.

The generate command needs access to the Application Accelerator server. You can specify the --server-url flag or set
an ACC_SERVER_URL environment variable. If you specify the --server-url flag it will override the ACC_SERVER_URL
environment variable if it is set.
`,
		ValidArgsFunction: SuggestAcceleratorNamesFromUiServer(context.Background()),
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return errors.New("you must specify the name of the accelerator")
			}
			return nil
		},
		Example: "tanzu accelerator generate <accelerator-name> --options '{\"projectName\":\"test\"}'",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !strings.HasSuffix(outputDir, "/") && outputDir != "" {
				outputDir += "/"
			}
			var options map[string]interface{}
			if filename != "" {
				fileBytes, err := ioutil.ReadFile(filename)
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
				options["projectName"] = args[0]
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
			if !strings.HasPrefix(serverUrl, "http://") && !strings.HasPrefix(serverUrl, "https://") {
				return errors.New(fmt.Sprintf("error creating request for %s, the URL needs to include the protocol (\"http://\" or \"https://\")", serverUrl))
			}

			osuser, _ := user.Current()
			provenanceId := uuid.New().String()

			apiPrefix := DetermineApiServerPrefix(serverUrl)
			proxyRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/accelerators/zip?name=%s&source=TanzuCLI&username=%s&id=%s", serverUrl, apiPrefix, args[0], osuser.Username, provenanceId), bytes.NewReader(JsonProxyBodyBytes))
			proxyRequest.Header.Add("Content-Type", "application/json")
			client := &http.Client{}
			resp, err := client.Do(proxyRequest)
			if err != nil {
				return err
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
			zipfile := outputDir + options["projectName"].(string) + ".zip"
			err = ioutil.WriteFile(zipfile, body, 0644)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "zip file %s created\n", zipfile)
			invokedRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/accelerators/invoked?type=download&name=%s&source=TanzuCLI&username=%s&id=%s", serverUrl, apiPrefix, args[0], osuser.Username, provenanceId), nil)
			if err != nil {
				return err
			}
			resp, err = client.Do(invokedRequest)
			if err != nil {
				return err
			}
			if resp.StatusCode == http.StatusNotFound {
				// try the deprecated downloaded endpoint for older servers
				downloadedRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/accelerators/downloaded?name=%s", serverUrl, apiPrefix, args[0]), nil)
				if err != nil {
					return err
				}
				resp, err = client.Do(downloadedRequest)
				if err != nil {
					return err
				}
			}
			if resp.StatusCode == http.StatusNotFound {
				return nil
			} else if resp.StatusCode >= 400 {
				var errorMsg string
				var errorResponse UiErrorResponse
				body, _ := ioutil.ReadAll(resp.Body)
				json.Unmarshal(body, &errorResponse)
				if errorResponse.Detail > "" {
					errorMsg = fmt.Sprintf("there was an error registering download for the accelerator, the server response was: \"%s\"\n", errorResponse.Detail)
				} else {
					errorMsg = fmt.Sprintf("there was an error registering download for the accelerator, the server response code was: \"%v\"\n", resp.StatusCode)
				}
				return fmt.Errorf(errorMsg)
			}

			return nil
		},
	}
	generateCmd.Flags().StringVar(&optionsString, "options", "{}", "options JSON string")
	generateCmd.Flags().StringVar(&filename, "options-file", "", "path to file containing options JSON string")
	generateCmd.Flags().StringVar(&outputDir, "output-dir", "", "directory that the zip file will be written to")
	generateCmd.Flags().StringVar(&uiServer, "server-url", "", "the URL for the Application Accelerator server")
	accServerUrl = EnvVar("ACC_SERVER_URL", "")
	return generateCmd
}
