package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"goph_keeper/cli/remote"
	"goph_keeper/pkg/auth"
	"goph_keeper/pkg/common"
	"io"
	"net/http"
)

// StoreLogoPass godoc
// @Summary      Сохранить креды
// @Description  Сохранить креды
// @Tags         register
// @Accept       json
// @Produce      json
// @Param        request body remote.StoreLogoPassRequest true "Параметры запроса"
// @Success 	 200
// @Failure      500
// @Router       /api/v1/creeds [post]
func (h *Handler) StoreLogoPass(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(map[string]interface{}{
		"HANDLER":    "StoreLogoPass",
		"REQUEST_ID": r.Context().Value(middleware.RequestIDKey),
	})
	logger.Info("Handling Store LogoPass")
	ctxUser, ok := (r.Context().Value(common.CtxUser)).(*auth.User)
	if !ok || ctxUser == nil {
		logger.Errorf("could not extract user from context")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("could not define user"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}
	var request *remote.StoreLogoPassRequest
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("could not read body"))
		if err != nil {
			logger.Errorf("could not read body")
		}
		return
	}
	err = json.Unmarshal(bodyBytes, &request)
	if err != nil || request == nil {
		logger.Errorf("could not unmarshal request")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("could not unmarshal request"))
		if err != nil {
			logger.Errorf("could not unmarshal request")
		}
		return
	}
	err = h.Storage.StoreLogoPass(r.Context(), ctxUser.ID, request.DataName, request.DataLogin, request.DataPassword, request.DataMeta)
	if err != nil {
		logger.Errorf("could not store logo pass err:%s", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(fmt.Sprintf("could not store logo pass err:%s", err)))
		if err != nil {
			logger.Errorf("could not store logo pass")
		}
		return
	}
}
