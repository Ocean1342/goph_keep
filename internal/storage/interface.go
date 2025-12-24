package storage

import (
	"context"
)

type Storage interface {
	GetUserByLogin(ctx context.Context, login string) (*User, error)
	CreateUser(ctx context.Context, login, password, phrase string) (*User, error)
	StoreLogoPass(ctx context.Context, userID int, dataName, sLogo, sPass, meta string) error
	GetLogoPass(ctx context.Context, userID int, dataName string) (*LogoPass, error)
	StoreBin(ctx context.Context, userID int, dataName, bin, meta string) error
	GetBin(ctx context.Context, userID int, dataName string) (*Bin, error)
}
