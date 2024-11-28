package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/semaphoreci/cli/cmd/utils"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Verbose bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "sem",
	Short: "Semaphore 2.0 command line interface",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !Verbose {
			log.SetOutput(ioutil.Discard)
		}
	},

	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := homedir.Dir()

	utils.CheckWithMessage(err, "failed to find home directory")

	// Search config in home directory with name ".sem" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".sem")

	// Touch config file and make sure that it exists
	path := fmt.Sprintf("%s/.sem.yaml", home)

	// #nosec
	_, _ = os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)

	err = viper.ReadInConfig()

	utils.CheckWithMessage(err, "failed to load config file")
}
