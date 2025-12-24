/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"goph_keeper/cli/application"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Регистрирует нового пользователя",
	Long: `Команда Регистрирует нового пользователя. Если нет связи с сервером, то получается ошибка.
	При этом добавлять данные в хранилище можно.`,
	PreRunE: application.Init,
	RunE: func(cmd *cobra.Command, args []string) error {
		app := application.App
		login, err := getLogin(cmd)
		if err != nil {
			return fmt.Errorf("get login failed: %v", err)
		}
		pass, err := getPass(cmd)
		if err != nil {
			return fmt.Errorf("get password failed: %v", err)
		}
		phrase := cmd.Flag("phrase").Value.String()
		if phrase == "" {
			return fmt.Errorf("phrase is required")
		}
		err = app.Auth.Register(context.TODO(), login, pass, phrase)
		if err != nil {
			return fmt.Errorf("register failed: %v", err)
		}
		app.Printer.Println("successfully register")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().String("phrase", "", "phrase")
}
