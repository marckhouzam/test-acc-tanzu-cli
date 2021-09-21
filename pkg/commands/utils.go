package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	acceleratorv1alpha1 "github.com/pivotal/acc-controller/api/v1alpha1"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu-private/tanzu-cli-apps-plugins/pkg/cli-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
		fmt.Fprintf(cmd.OutOrStderr(), "Error getting accelerators from %s, check that the ACC_SERVER_URL"+
			" env variable is set with the correct value, or use the --from-context flag to get the accelerators from your current context\n", url)
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

func SuggestAcceleratorNamesFromUiServer(ctx context.Context) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		suggestions := []string{}
		uiServerUrl := EnvVar("ACC_SERVER_URL", "http://localhost:8877")
		if cmd.Flags().Changed("server-url") {
			uiServerUrl, _ = cmd.Flags().GetString("server-url")
		}
		var response UiAcceleratorList
		resp, err := http.Get(uiServerUrl + "/api/accelerators")
		if err != nil {
			return suggestions, cobra.ShellCompDirectiveError
		}
		defer resp.Body.Close()
		jsonBody, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(jsonBody, &response)
		for _, accelerator := range response.Embedded.Accelerators {
			suggestions = append(suggestions, accelerator.Name)
		}
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	}
}

func SuggestAcceleratorNamesFromConfig(ctx context.Context, c *cli.Config) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		suggestions := []string{}
		accelerators := &acceleratorv1alpha1.AcceleratorList{}
		err := c.List(ctx, accelerators, client.InNamespace(cmd.Flag("namespace").Value.String()))
		if err != nil {
			return suggestions, cobra.ShellCompDirectiveError
		}
		for _, accelerator := range accelerators.Items {
			suggestions = append(suggestions, accelerator.Name)
		}
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	}
}
