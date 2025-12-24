package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"goph_keeper/config"
	"goph_keeper/internal/api"
	"net/http"
)

func Init(cfg *config.Config, handler *api.Handler) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	//r.Use(loggable)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/user/register", handler.UserRegister)
		r.Post("/user/login", handler.UserAuth)
		r.With(handler.Authenticate).Post("/creeds", handler.StoreLogoPass)
		r.With(handler.Authenticate).Post("/creeds/get", handler.GetLogoPass)
		r.With(handler.Authenticate).Post("/bin", handler.StoreBin)
		r.With(handler.Authenticate).Post("/bin/get", handler.GetBin)
	})

	server := &http.Server{
		Addr:    cfg.RunAddr,
		Handler: r,
	}
	logrus.Debugf("server started at %s", server.Addr)

	var err error
	if cfg.TLSEnabled {
		err = server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil {
		logrus.Errorf("could not start server. err: %v", err)
	}
}
