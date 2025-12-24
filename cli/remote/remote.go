package remote

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"goph_keeper/cli/config"
	"goph_keeper/pkg/auth"
	"net/http"
)

var (
	ErrInvalidData = errors.New("invalid data")
)

type GetLogoPassRequest struct {
	UserLogin string `json:"user_login"`
	DataName  string `json:"data_name"`
}

type GetLogoPassResponse struct {
	UserLogin string `json:"user_login"`
	UserPass  string `json:"user_pass"`
	Meta      string `json:"meta"`
}

// Requester - сервис отвечает за отправку и получение данных на удалённый сервер
type Requester interface {
	Register(ctx context.Context, login, pass, phrase string) (auth.Token, error)
	Auth(ctx context.Context, login, pass string) (auth.Token, error)
	//logopass
	SendLogoPass(ctx context.Context, req StoreLogoPassRequest) error
	GetLogoPass(ctx context.Context, req GetLogoPassRequest, token string) (*GetLogoPassResponse, error)
	//
	SendBinData(ctx context.Context, req StoreBin) error
	GetBinData(ctx context.Context, req GetBinRequest) (*GetBinResponse, error)
}

type Sender struct {
	clientHTTP *resty.Client
	config     *config.Config
}

func New(clientHTTP *resty.Client, cfg *config.Config) (Requester, error) {
	return &Sender{
		clientHTTP: clientHTTP,
		config:     cfg,
	}, nil
}

type SendContract struct {
	//user + data + meta + dataType + token!
}

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s Sender) Auth(_ context.Context, login, pass string) (auth.Token, error) {
	resp, err := s.clientHTTP.R().
		SetHeader("Content-Type", "application/json").
		SetBody(authRequest{
			Login:    login,
			Password: pass,
		}).
		Post(s.config.AuthURL)
	if err != nil {
		return "", fmt.Errorf("could not send register data. err:%w", err)
	}
	authHeader := resp.Header().Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("could not get auth token from service response")
	}
	return auth.Token(authHeader), nil
}

type registerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Phrase   string `json:"phrase"`
}

// Register - регистрирует нового пользователя.
// Данные передаются открыто, тк на сервере используется TLS
func (s Sender) Register(_ context.Context, login, pass, phrase string) (auth.Token, error) {
	resp, err := s.clientHTTP.R().
		SetHeader("Content-Type", "application/json").
		SetBody(registerRequest{
			Login:    login,
			Password: pass,
			Phrase:   phrase,
		}).
		Post(s.config.RegisterURL)
	if err != nil {
		logrus.Errorf("confli err:%v", err)
		return "", fmt.Errorf("could not send register data. err:%w", err)
	}
	//409 приходит в случае уже существующего пользователя
	switch resp.StatusCode() {
	case http.StatusConflict:

		return "", ErrInvalidData
	}
	authHeader := resp.Header().Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("could not get auth token from service response")
	}
	return auth.Token(authHeader), nil
}
