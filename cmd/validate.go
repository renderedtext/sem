package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

type YamlBody struct {
	YamlDefinition string `json:"yaml_definition"`
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a pipeline yaml file",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewBaseClientFromConfig()

		if len(args) != 1 {
			fmt.Println("Usage: sem validate [PATH_TO_YAML_FILE]")
			os.Exit(1)
		}

		fileContent, err := os.ReadFile(args[0])
		utils.Check(err)

		body := YamlBody{
			YamlDefinition: string(fileContent),
		}

		json_body, err := json.Marshal(body)
		utils.Check(err)

		response, status, err := c.Post("yaml", json_body)
		utils.Check(err)

		if status == 200 {
			fmt.Println("Valid!")
		} else {
			fmt.Println("Not valid!")
		}
		fmt.Println(string(response))
	},
}

func init() {
	RootCmd.AddCommand(validateCmd)
}
