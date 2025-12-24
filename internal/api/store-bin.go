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

// StoreBin godoc
// @Summary      Сохранить bin
// @Description  Сохранить bin
// @Tags         register
// @Accept       json
// @Produce      json
// @Param        request body remote.StoreBin true "Параметры запроса"
// @Success      200
// @Failure      401
// @Failure      500
// @Router       /api/v1/creeds [post]
func (h *Handler) StoreBin(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(map[string]interface{}{
		"HANDLER":    "StoreBin",
		"REQUEST_ID": r.Context().Value(middleware.RequestIDKey),
	})
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
	var request *remote.StoreBin
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
	err = h.Storage.StoreBin(r.Context(), ctxUser.ID, request.DataName, request.BinData, request.DataMeta)
	if err != nil {
		logger.Errorf("could not store bin err:%s", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(fmt.Sprintf("could not store bin err:%s", err)))
		if err != nil {
			logger.Errorf("could not store bin")
		}
		return
	}
}
