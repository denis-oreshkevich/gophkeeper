package server

import (
	"encoding/json"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func (c *Controller) HandlePostClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client, err := c.userClient.RegisterClient(ctx)

	if err != nil {
		logger.Log.Error("userClient.RegisterClient", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(client)
	if err != nil {
		logger.Log.Error("json.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

func (c *Controller) HandlePutClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("io.ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var client domain.Client
	err = json.Unmarshal(body, &client)
	if err != nil {
		logger.Log.Error("json.Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = c.syncSvc.UpdateClientLastSyncTms(ctx, client)
	if err != nil {
		logger.Log.Error("syncSvc.UpdateClientLastSyncTms", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
