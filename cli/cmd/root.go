package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"goph_keeper/cli/application"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cli",
	Short:   "short",
	Long:    ``,
	PreRunE: application.Init,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("root cmd worked every time")
	},
	PostRun: application.ShutDown,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringP(Login, "l", "", "login name")
	rootCmd.PersistentFlags().StringP(Pass, "p", "", "password")
	rootCmd.PersistentFlags().StringP(Meta, "m", "", "metadata for stored data")
	rootCmd.PersistentFlags().StringP(Name, "n", "", "data name")
}
