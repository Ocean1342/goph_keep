package action

import (
	"context"
	"errors"
	"fmt"
	"goph_keeper/cli/application"
	"goph_keeper/cli/auth"
	"goph_keeper/cli/remote"
	"os"
)

func (c *Case) SendBinData(ctx context.Context, req remote.StoreBin) error {
	app := application.App
	user, err := app.Auth.Auth(ctx, req.UserLogin)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotRegistered) {
			return fmt.Errorf("user: `%v` not registered. please, run register command", req.UserLogin)
		}
		return fmt.Errorf("unknown error. please try again later. err:%v", err)
	}
	req.Token = user.Token
	data, err := os.ReadFile(req.BinDataPath)
	if err != nil {
		return fmt.Errorf("read file error. please try again later. err:%v", err)
	}
	cryptedData, err := app.Crypter.Encode(ctx, data, user)
	if err != nil {
		return fmt.Errorf("encode error:%v", err)
	}
	err = app.Storage.StoreBin(ctx, user.ID, req.DataName, cryptedData, req.DataMeta)
	if err != nil {
		return fmt.Errorf("store bin error:%v", err)
	}
	return nil
}
