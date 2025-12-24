package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const (
	Login = "login"
	Pass  = "pass"
	Meta  = "meta"
	Name  = "name"
)

func getLogin(cmd *cobra.Command) (string, error) {
	login := cmd.Flag(Login).Value.String()
	if login == "" {
		return "", fmt.Errorf("login is required")
	}
	return login, nil
}

func getPass(cmd *cobra.Command) (string, error) {
	pass := cmd.Flag(Pass).Value.String()
	if pass == "" {
		return "", fmt.Errorf("login is required")
	}
	return pass, nil
}

func getMeta(cmd *cobra.Command) (string, error) {
	meta := cmd.Flag(Meta).Value.String()
	if meta == "" {
		return "", fmt.Errorf("meta is required")
	}
	return meta, nil
}

func getDataName(cmd *cobra.Command) (string, error) {
	name := cmd.Flag(Name).Value.String()
	if name == "" {
		return "", fmt.Errorf("data name is required")
	}
	return name, nil
}
