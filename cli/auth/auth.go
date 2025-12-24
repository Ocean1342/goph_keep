package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"goph_keeper/cli/entity"
	"goph_keeper/cli/remote"
	"goph_keeper/cli/storage"
)

var (
	ErrUserNotRegistered = errors.New("user not registered")
)

// Auth - сервис отвечает за авторизацию и регистрацию на удалённом сервере.
//
//	является логической прокладкой для инкапсуляции запросов к серверу
//	и поведения если сервер недоступен
type Auth struct {
	store  storage.Repository
	sender remote.Requester
}

func New(store storage.Repository, sender remote.Requester) (*Auth, error) {
	return &Auth{store: store, sender: sender}, nil
}

// Register - регистрация пользователя на удалённом сервере и сохранение данных в локальном хранилище.
func (a Auth) Register(ctx context.Context, login, pass, phrase string) error {
	//сохранить в локальный стор
	err := a.store.StoreInitialUserData(ctx, login, pass, phrase)
	if err != nil {
		return fmt.Errorf("failed store user: %w", err)
	}
	//попытаться сохранить на сервер
	token, err := a.sender.Register(ctx, login, pass, phrase)
	if err != nil {
		return fmt.Errorf("failed send register data: %w", err)
	}
	if token != "" {
		err = a.store.StoreToken(ctx, login, token)
		if err != nil {
			return fmt.Errorf("failed store token: %w", err)
		}
	}
	return nil
}

// Auth - авторизация пользователя по логину
func (a Auth) Auth(ctx context.Context, login string) (*entity.User, error) {
	u, err := a.store.GetLocalUserByLogin(ctx, login)
	if err != nil || u == nil {
		logrus.Errorf("get local user by login err: %v", err)
		return nil, ErrUserNotRegistered
	}
	if u.IsAuthenticated {
		return u, nil
	}
	//a.sender.Auth()
	return u, nil
}
