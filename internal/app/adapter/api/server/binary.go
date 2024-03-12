package server

import (
	"encoding/json"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func (c *Controller) HandlePostBinary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("io.ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var binary domain.Binary
	err = json.Unmarshal(body, &binary)
	if err != nil {
		logger.Log.Error("json.Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = c.crudSvc.SaveBinary(ctx, &binary)
	if err != nil {
		logger.Log.Error("crudSvc.SaveBinary", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (c *Controller) HandleDeleteBinaryByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	err := c.crudSvc.DeleteBinaryByID(ctx, id)
	if err != nil {
		logger.Log.Error("crudSvc.DeleteBinaryByID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (c *Controller) HandleGetUserBinaries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	binaries, err := c.crudSvc.FindBinariesByUserID(ctx)
	if err != nil {
		logger.Log.Error("crudSvc.FindBinariesByUserID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(binaries) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	result, err := json.Marshal(binaries)
	if err != nil {
		logger.Log.Error("json.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (c *Controller) HandleGetBinaryByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	binary, err := c.crudSvc.FindBinaryByID(ctx, id)
	if err != nil {
		logger.Log.Error("crudSvc.FindBinaryByID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(binary)
	if err != nil {
		logger.Log.Error("json.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (c *Controller) HandlePostSyncBinary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("io.ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var binarySync domain.BinarySync
	err = json.Unmarshal(body, &binarySync)
	if err != nil {
		logger.Log.Error("json.Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(binarySync.Binaries) == 0 {
		log.Debug("binarySync.Binaries length for sync = 0")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	binaries, err := c.syncSvc.SyncBinary(ctx, &binarySync)
	if err != nil {
		logger.Log.Error("syncSvc.SyncBinary", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(binaries)
	if err != nil {
		logger.Log.Error("json.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(result)
}
