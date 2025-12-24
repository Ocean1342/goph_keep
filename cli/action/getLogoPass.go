package action

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"goph_keeper/cli/application"
	"goph_keeper/cli/auth"
	"goph_keeper/cli/remote"
)

func (c *Case) GetLogoPass(ctx context.Context, req remote.GetLogoPassRequest) (*remote.GetLogoPassResponse, error) {
	app := application.App
	user, err := app.Auth.Auth(ctx, req.UserLogin)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotRegistered) {
			return nil, fmt.Errorf("user: `%v` not registered. please, run register command", req.UserLogin)
		}
		return nil, fmt.Errorf("unknown error. please try again later. err:%v", err)
	}
	//пытаемся извлечь данные с удалённого сервера
	var resp *remote.GetLogoPassResponse
	resp, err = app.Remote.GetLogoPass(ctx, req, user.Token)
	if err != nil {
		logrus.Errorf("get logo pass error. err:%v", err)
		//пытаемся получить из локального хранилища
		resp = &remote.GetLogoPassResponse{}
		resp.UserLogin, resp.UserPass, resp.Meta, err = app.Storage.GetLogoPass(ctx, user.ID, req.DataName)
		if err != nil {
			return nil, fmt.Errorf("could not get data. err:%v", err)
		}
	}
	var decryptedLogo string
	var decryptedPass string
	if resp.UserLogin != "" {
		decodeBytes, err := hex.DecodeString(resp.UserLogin)
		if err != nil {
			return nil, fmt.Errorf("could not hex decode logo. err:%v", err)
		}
		decryptedLogo, err = app.Crypter.Decode(ctx, decodeBytes, user)
		if err != nil {
			return nil, fmt.Errorf("could not decrypt data. err:%v", err)
		}
	}
	if resp.UserPass != "" {
		decodeBytes, err := hex.DecodeString(resp.UserPass)
		if err != nil {
			return nil, fmt.Errorf("could not hex decode pass. err:%v", err)
		}
		decryptedPass, err = app.Crypter.Decode(ctx, decodeBytes, user)
		if err != nil {
			return nil, fmt.Errorf("could not decrypt data. err:%v", err)
		}
	}
	return &remote.GetLogoPassResponse{UserLogin: decryptedLogo, UserPass: decryptedPass, Meta: resp.Meta}, nil
}
