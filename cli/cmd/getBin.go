package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"goph_keeper/cli/action"
	"goph_keeper/cli/application"
	"goph_keeper/cli/remote"
	"os"
)

const BinPutTo = "binPathTo"

// getBinCmd represents the getBin command
var getBinCmd = &cobra.Command{
	Use:     "getBin",
	Short:   "get bin data",
	Long:    `get bin data`,
	PreRunE: application.Init,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		uLogin, err := getLogin(cmd)
		if err != nil {
			return fmt.Errorf("getLogin error: %v", err)
		}
		user, err := application.App.Storage.GetLocalUserByLogin(ctx, uLogin)
		if err != nil {
			return fmt.Errorf("getLocalUserByLogin error: %v", err)
		}
		name, err := getDataName(cmd)
		if err != nil {
			return fmt.Errorf("getDataName error: %v", err)
		}
		binDataPath, err := cmd.Flags().GetString(BinPutTo)
		if err != nil {
			application.App.Printer.Println("empty path. will used ./")
			binDataPath = "./"
		}

		resp, err := action.New().GetBinData(ctx, remote.GetBinRequest{
			UserLogin: uLogin,
			DataName:  name,
			Token:     user.Token,
		})
		if err != nil {
			return fmt.Errorf("sendBinData error: %v", err)
		}
		//перегнать в файл
		fPath := binDataPath + name
		_, err = os.Create(fPath)
		if err != nil {
			return fmt.Errorf("create file error: %v", err)
		}
		err = os.WriteFile(fPath, []byte(resp.BinData), 0755)
		if err != nil {
			return fmt.Errorf("writeFile error: %v", err)
		}
		return nil
	},
	PostRun: application.ShutDown,
}

func init() {
	getBinCmd.Flags().String(BinPutTo, "", "path where put data")
	rootCmd.AddCommand(getBinCmd)
}
