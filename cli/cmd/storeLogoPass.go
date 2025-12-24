package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"goph_keeper/cli/action"
	"goph_keeper/cli/application"
	"goph_keeper/cli/remote"
)

const StoreLogo = "storeLogo"
const StorePass = "storePass"

// storeLogoPassCmd represents the storeLogoPass command
var storeLogoPassCmd = &cobra.Command{
	Use:     "storeLogoPass",
	Short:   "сохраняет логин и пароль для сайтов",
	Long:    `сохраняет логин и пароль для сайтов`,
	PreRunE: application.Init,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		uLogin, err := getLogin(cmd)
		if err != nil {
			return fmt.Errorf("getLogin error: %v", err)
		}
		dataLogo, err := cmd.Flags().GetString(StoreLogo)
		if err != nil {
			return fmt.Errorf("get store logo: %w", err)
		}
		dataPass, err := cmd.Flags().GetString(StorePass)
		if err != nil {
			return fmt.Errorf("get store pass: %w", err)
		}
		meta, err := getMeta(cmd)
		if err != nil {
			return fmt.Errorf("getMeta error: %v", err)
		}
		name, err := getDataName(cmd)
		if err != nil {
			return fmt.Errorf("getDataName error: %v", err)
		}
		req := remote.StoreLogoPassRequest{
			UserLogin:    uLogin,
			DataName:     name,
			DataLogin:    dataLogo,
			DataPassword: dataPass,
			DataMeta:     meta,
		}

		err = action.New().StoreLogoPass(ctx, req)
		if err != nil {
			return fmt.Errorf("storelogoPass error: %v", err)
		}
		return nil
	},
	PostRun:      application.ShutDown,
	SilenceUsage: false,
}

func init() {
	rootCmd.AddCommand(storeLogoPassCmd)
	storeLogoPassCmd.Flags().String(StoreLogo, "", "Stored login")
	storeLogoPassCmd.Flags().String(StorePass, "", "Stored password")
}
