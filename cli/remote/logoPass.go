package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"goph_keeper/pkg/common"
	"net/http"
)

// StoreLogoPassRequest - DTO верхнего слоя,
// но для скорости и простоты он используется в других сервисах как сквозная
type StoreLogoPassRequest struct {
	UserLogin    string `json:"user_login"`
	DataName     string `json:"data_name"`
	DataLogin    string `json:"data_login"`
	DataPassword string `json:"data_pass"`
	DataMeta     string `json:"data_meta"`
	Token        string `json:"token"`
}

func (s Sender) SendLogoPass(ctx context.Context, req StoreLogoPassRequest) error {
	resp, err := s.clientHTTP.R().SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader(common.AuthorizationHeaderName, req.Token).
		SetBody(req).
		Post(s.config.SendLogoPassURL)
	if err != nil {
		return fmt.Errorf("could not store logo-pass data. err:%w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return ErrInvalidData
	}
	return nil
}

func (s Sender) GetLogoPass(ctx context.Context, req GetLogoPassRequest, token string) (*GetLogoPassResponse, error) {
	resp, err := s.clientHTTP.R().SetContext(ctx).
		SetHeader(common.AuthorizationHeaderName, token).
		SetBody(req).
		Post(s.config.GetLogoPassURL)
	if err != nil {
		return nil, fmt.Errorf("could not get logo pass data. err:%w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, ErrInvalidData
	}
	var respLogoPass *GetLogoPassResponse
	err = json.Unmarshal(resp.Body(), &respLogoPass)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal logo pass data. err:%w", err)
	}
	return respLogoPass, nil
}
