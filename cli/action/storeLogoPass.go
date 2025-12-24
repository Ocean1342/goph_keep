package action

import (
	"context"
	"errors"
	"fmt"
	"goph_keeper/cli/application"
	"goph_keeper/cli/auth"
	"goph_keeper/cli/remote"
)

func (c *Case) StoreLogoPass(ctx context.Context, req remote.StoreLogoPassRequest) error {
	app := application.App
	user, err := app.Auth.Auth(ctx, req.UserLogin)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotRegistered) {
			return fmt.Errorf("user: `%v` not registered. please, run register command", req.UserLogin)
		}
		return fmt.Errorf("unknown error. please try again later. err:%v", err)
	}
	req.Token = user.Token
	cryptedLogo, err := app.Crypter.Encode(ctx, []byte(req.DataLogin), user)
	if err != nil {
		return fmt.Errorf("encode error:%v", err)
	}
	cryptedPass, err := app.Crypter.Encode(ctx, []byte(req.DataPassword), user)
	if err != nil {
		return fmt.Errorf("encode error:%v", err)
	}

	err = app.Storage.StoreLogoPass(ctx, user.ID, req.DataName, cryptedLogo, cryptedPass, req.DataMeta)
	if err != nil {
		return fmt.Errorf("store data: %w", err)
	}
	app.Printer.Println("successfully stored")
	return nil
}
