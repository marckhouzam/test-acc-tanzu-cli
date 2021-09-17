package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func EnvVar(key string, defVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defVal
}

type Accelerator struct {
	Name                 string   `json:"name"`
	IconUrl              string   `json:"iconUrl,omitempty"`
	SourceUrl            string   `json:"sourceUrl,omitempty"`
	SpecGitRepositoryUrl string   `json:"specGitRepositoryUrl,omitempty"`
	SourceBranch         string   `json:"sourceBranch,omitempty"`
	SourceTag            string   `json:"sourceTag,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
	Description          string   `json:"description,omitempty"`
	DisplayName          string   `json:"displayName,omitempty"`
	Ready                bool     `json:"ready,omitempty"`
	ArchiveUrl           string   `json:"archiveUrl,omitempty"`
	ArchiveReady         bool     `json:"archiveReady,omitempty"`
	ArchiveMessage       string   `json:"archiveMessage,omitempty"`
}

type Embedded struct {
	Accelerators []Accelerator `json:"accelerators"`
}

type UiAcceleratorsApiResponse struct {
	Emdedded Embedded `json:"_embedded"`
}
type Choice struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

type Option struct {
	Name         string      `json:"name"`
	DefaultValue interface{} `json:"defaultValue" yaml:"defaultValue"`
	Display      bool        `json:"display"`
	DataType     interface{} `json:"dataType" yaml:"dataType"`
	Choices      []Choice    `json:"choices,omitempty"`
}

type OptionsResponse struct {
	Options []Option `json:"options"`
}

func GetAcceleratorsFromUiServer(url string, cmd *cobra.Command) ([]Accelerator, error) {
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("%s/api/accelerators", url))
	if err != nil {
		fmt.Fprintf(cmd.OutOrStderr(), "Error getting accelerators from %s\n", url)
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var uiResponse UiAcceleratorsApiResponse
	defer resp.Body.Close()
	err = json.Unmarshal(body, &uiResponse)
	if err != nil {
		fmt.Fprintf(cmd.OutOrStderr(), "Error unmarshalling response\n")
		log.Fatal(err)
		return nil, err
	}
	return uiResponse.Emdedded.Accelerators, nil
}

func GetAcceleratorOptionsFromUiServer(url string, acceleratorName string, cmd *cobra.Command) ([]Option, error) {
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("%s/api/accelerators/options?name=%s", url, acceleratorName))
	if err != nil {
		fmt.Fprintf(cmd.OutOrStderr(), "Error getting accelerator %s options from %s\n", acceleratorName, url)
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var optionsResponse OptionsResponse
	defer resp.Body.Close()
	err = json.Unmarshal(body, &optionsResponse)
	if err != nil {
		fmt.Fprintf(cmd.OutOrStderr(), "Error unmarshalling response\n")
		log.Fatal(err)
		return nil, err
	}
	return optionsResponse.Options, nil
}
