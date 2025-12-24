package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"goph_keeper/pkg/common"
	"net/http"
)

type StoreBin struct {
	UserLogin   string `json:"user_login"`
	DataName    string `json:"data_name"`
	BinDataPath string `json:"bin_data_path"`
	BinData     string `json:"bin_data"`
	DataMeta    string `json:"data_meta"`
	Token       string `json:"token"`
}

type GetBinRequest struct {
	UserLogin string `json:"user_login"`
	DataName  string `json:"data_name"`
	Token     string `json:"token"`
}

type GetBinResponse struct {
	BinData  string `json:"bin"`
	DataMeta string `json:"meta"`
}

func (s Sender) SendBinData(ctx context.Context, req StoreBin) error {
	resp, err := s.clientHTTP.R().SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader(common.AuthorizationHeaderName, req.Token).
		SetBody(req).
		Post(s.config.SendBin)
	if err != nil {
		return fmt.Errorf("could not request bin data. err:%w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		fmt.Println(err)
		return ErrInvalidData
	}
	return nil
}

func (s Sender) GetBinData(ctx context.Context, req GetBinRequest) (*GetBinResponse, error) {
	resp, err := s.clientHTTP.R().SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader(common.AuthorizationHeaderName, req.Token).
		SetBody(req).
		Post(s.config.GetBin)
	if err != nil {
		return nil, fmt.Errorf("could not request bin data. err:%w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, ErrInvalidData
	}
	var binResp *GetBinResponse
	err = json.Unmarshal(resp.Body(), &binResp)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal bin data. err:%w", err)
	}
	return binResp, nil
}
