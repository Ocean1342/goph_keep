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

const BinPath = "storeBin"

// storeBinCmd represents the storeBin command
var storeBinCmd = &cobra.Command{
	Use:     "storeBin",
	Short:   "store bin data",
	Long:    `store bin data max size 5mb`,
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
		meta, err := getMeta(cmd)
		if err != nil {
			return fmt.Errorf("getMeta error: %v", err)
		}
		name, err := getDataName(cmd)
		if err != nil {
			return fmt.Errorf("getDataName error: %v", err)
		}
		binDataPath, err := cmd.Flags().GetString(BinPath)
		if err != nil {
			return fmt.Errorf("get binData: %w", err)
		}
		if _, err = os.ReadFile(binDataPath); err != nil {
			return fmt.Errorf("read bin file err: %w", err)
		}

		err = action.New().SendBinData(ctx, remote.StoreBin{
			UserLogin:   uLogin,
			DataName:    name,
			BinDataPath: binDataPath,
			DataMeta:    meta,
			Token:       user.Token,
		})
		if err != nil {
			return fmt.Errorf("sendBinData error: %v", err)
		}
		application.App.Printer.Println("successfully stored bin data")
		return nil
	},
	PostRun: application.ShutDown,
}

func init() {
	rootCmd.AddCommand(storeBinCmd)
	storeBinCmd.Flags().String(BinPath, "", "Stored bin")
}
