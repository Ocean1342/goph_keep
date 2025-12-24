/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Version string
	Date    string
)

// buildDataCmd represents the buildData command
var buildDataCmd = &cobra.Command{
	Use:   "buildData",
	Short: "Возвращает информацию о дате сборки и версии",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s, date: %s \n", Version, Date)
	},
}

func init() {
	rootCmd.AddCommand(buildDataCmd)
}
