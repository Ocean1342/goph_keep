/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"goph_keeper/cli/action"
	"goph_keeper/cli/application"
	"goph_keeper/cli/remote"

	"github.com/spf13/cobra"
)

// getLogoPassCmd represents the getLogoPass command
var getLogoPassCmd = &cobra.Command{
	Use:     "getLogoPass",
	Short:   "get stored creeds",
	Long:    ``,
	PreRunE: application.Init,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		uLogin, err := getLogin(cmd)
		if err != nil {
			return fmt.Errorf("getLogin error: %v", err)
		}
		name, err := getDataName(cmd)
		if err != nil {
			return fmt.Errorf("could get data name. err:%v", err)
		}
		req := remote.GetLogoPassRequest{UserLogin: uLogin, DataName: name}
		resp, err := action.New().GetLogoPass(ctx, req)
		if err != nil {
			return fmt.Errorf("logoPass error: %v", err)
		}
		application.App.Printer.Println(
			fmt.Sprintf("log: %v pass: %v info:%v", resp.UserLogin, resp.UserPass, resp.Meta),
		)
		return nil
	},
	PostRun: application.ShutDown,
}

func init() {
	rootCmd.AddCommand(getLogoPassCmd)
}
