package api

import (
	"goph_keeper/internal/storage"
	"goph_keeper/pkg/auth"
)

type Handler struct {
	Storage storage.Storage
	Auth    auth.Auth
}

func New(storage storage.Storage, auth auth.Auth) *Handler {
	return &Handler{
		Storage: storage,
		Auth:    auth,
	}
}
