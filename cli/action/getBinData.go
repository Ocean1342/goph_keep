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

func (c *Case) GetBinData(ctx context.Context, req remote.GetBinRequest) (*remote.GetBinResponse, error) {
	app := application.App
	user, err := app.Auth.Auth(ctx, req.UserLogin)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotRegistered) {
			return nil, fmt.Errorf("user: `%v` not registered. please, run register command", req.UserLogin)
		}
		return nil, fmt.Errorf("unknown error. please try again later. err:%v", err)
	}
	resp, err := app.Remote.GetBinData(ctx, req)
	if err != nil {
		resp = &remote.GetBinResponse{}
		logrus.Errorf("get bin data error. err:%v", err)
		//пытаемся получить из локального хранилища
		cryptedBin, dataMeta, err := app.Storage.GetBin(ctx, user.ID, req.DataName)
		if err != nil {
			return nil, fmt.Errorf("could not get data. err:%v", err)
		}
		resp.DataMeta = dataMeta
		resp.BinData = cryptedBin
	}
	if resp.BinData == "" {
		return nil, fmt.Errorf("no data found for user `%v`", req.UserLogin)
	}
	decodeBytes, err := hex.DecodeString(resp.BinData)
	if err != nil {
		return nil, fmt.Errorf("could not hex decode logo. err:%v", err)
	}
	decryptedBin, err := app.Crypter.Decode(ctx, decodeBytes, user)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt data. err:%v", err)
	}
	resp.BinData = decryptedBin
	return resp, nil
}
