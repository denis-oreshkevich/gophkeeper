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

func (c *Controller) HandlePostCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("io.ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var cred domain.Credentials
	err = json.Unmarshal(body, &cred)
	if err != nil {
		logger.Log.Error("json.Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = c.crudSvc.SaveCredentials(ctx, cred)
	if err != nil {
		logger.Log.Error("crudSvc.SaveCredentials", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (c *Controller) HandleDeleteCredentialsByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	err := c.crudSvc.DeleteCredentialsByID(ctx, id)
	if err != nil {
		logger.Log.Error("crudSvc.DeleteCredentialsByID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (c *Controller) HandleGetUserCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	creds, err := c.crudSvc.FindCredentialsByUserID(ctx)
	if err != nil {
		logger.Log.Error("crudSvc.FindCredentialsByUserID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(creds) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	result, err := json.Marshal(creds)
	if err != nil {
		logger.Log.Error("json.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (c *Controller) HandleGetCredentialsByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	cred, err := c.crudSvc.FindCredentialsByID(ctx, id)
	if err != nil {
		logger.Log.Error("crudSvc.FindCredentialsByID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(cred)
	if err != nil {
		logger.Log.Error("json.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (c *Controller) HandlePostSyncCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("io.ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var credSync domain.CredSync
	err = json.Unmarshal(body, &credSync)
	if err != nil {
		logger.Log.Error("json.Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(credSync.Credentials) == 0 {
		log.Debug("credSync.Credentials length for sync = 0")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	credentials, err := c.syncSvc.SyncCredentials(ctx, &credSync)
	if err != nil {
		logger.Log.Error("syncSvc.SyncCredentials", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(credentials)
	if err != nil {
		logger.Log.Error("json.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(result)
}
