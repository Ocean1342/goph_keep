package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"goph_keeper/cli/remote"
	"goph_keeper/pkg/auth"
	"goph_keeper/pkg/common"
	"io"
	"net/http"
)

// GetBin godoc
// @Summary      Получить bin
// @Description  Получить bin
// @Tags         register
// @Accept       json
// @Produce      json
// @Param        request body remote.GetBinRequest true "Параметры запроса"
// @Success      200  {object} storage.Bin
// @Failure      500
// @Router       /api/v1/bin [post]
func (h *Handler) GetBin(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(map[string]interface{}{
		"HANDLER":    "GetBin",
		"REQUEST_ID": r.Context().Value(middleware.RequestIDKey),
	})
	logger.Info("Handling Store Bin")
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

	var request *remote.GetBinRequest
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

	data, err := h.Storage.GetBin(r.Context(), ctxUser.ID, request.DataName)
	if err != nil {
		logger.Errorf("could not get data. err:%v", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("could not unmarshal request"))
		if err != nil {
			logger.Errorf("could not unmarshal request")
		}
		return
	}

	resp, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("could not marshal data")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not unmarshal request"))
		return
	}
	w.Write(resp)

}
