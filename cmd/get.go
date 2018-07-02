package cmd

import (
	"github.com/renderedtext/sem/cmd/handler"
	"github.com/renderedtext/sem/cmd/utils"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "List of resources.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		RunGet(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func RunGet(cmd *cobra.Command, args []string) {
	kind := args[0]

	params := handler.GetParams{}
	handler, err := handler.FindHandler(kind)

	utils.Check(err, err.Error())

	handler.Get(params)
}
