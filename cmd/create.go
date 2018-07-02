package cmd

import (
	"io/ioutil"

	"github.com/renderedtext/sem/cmd/handler"
	"github.com/renderedtext/sem/cmd/utils"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource from a file.",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		RunCreate(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	desc := "Filename, directory, or URL to files to use to create the resource"
	createCmd.PersistentFlags().StringP("file", "f", "", desc)
}

func RunCreate(cmd *cobra.Command, args []string) {
	path, err := cmd.Flags().GetString("file")

	utils.Check(err, "Path not provided")

	data, err := ioutil.ReadFile(path)

	utils.Check(err, "Failed to read from resource file.")

	resource, err := parse(data)

	utils.Check(err, "Failed to parse resource file.")

	apiVersion := resource["apiVersion"].(string)
	kind := resource["kind"].(string)

	json_resource, err := yaml.YAMLToJSON(data)

	utils.Check(err, "Failed to parse resource file.")

	params := handler.CreateParams{ApiVersion: apiVersion, Resource: json_resource}
	handler, err := handler.FindHandler(kind)

	utils.Check(err, err.Error())

	handler.Create(params)
}

func parse(data []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	// fmt.Print(string(data))

	err := yaml.Unmarshal(data, &m)

	return m, err
}
