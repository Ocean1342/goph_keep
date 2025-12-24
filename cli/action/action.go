package action

import (
	"context"
	"goph_keeper/cli/remote"
)

type ICase interface {
	StoreLogoPass(ctx context.Context, req remote.StoreLogoPassRequest) error
	GetLogoPass(ctx context.Context, req remote.GetLogoPassRequest) (*remote.GetLogoPassResponse, error)
	SendBinData(ctx context.Context, req remote.StoreBin) error
	GetBinData(ctx context.Context, req remote.GetBinRequest) (*remote.GetBinResponse, error)
}

type Case struct {
}

func New() ICase {
	return &Case{}
}
