package storage

import (
	"context"
	"goph_keeper/cli/entity"
	"goph_keeper/pkg/auth"
)

type IStore interface {
	GetLocalUserByLogin(ctx context.Context, login string) (*entity.User, error)
	StoreInitialUserData(ctx context.Context, login, pass, phrase string) error
	StoreToken(ctx context.Context, login string, token auth.Token) error
	StoreLogoPass(ctx context.Context, userID int, name, storeLogo, storePass, meta string) error
	GetLogoPass(ctx context.Context, userID int, name string) (logo, pass, meta string, err error)
	GetNotSyncedItems(ctx context.Context) ([]CatalogItem, error)
	UpdateCatalogItem(ctx context.Context, userID int, dataID string) error
	StoreBin(ctx context.Context, userID int, name, storeBin, meta string) error
	GetBin(ctx context.Context, userID int, name string) (bin, meta string, err error)
}
