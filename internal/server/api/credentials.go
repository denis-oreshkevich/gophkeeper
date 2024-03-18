package api

import (
	"encoding/json"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/logger"
	model2 "github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
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
	var cred model2.Credentials
	err = json.Unmarshal(body, &cred)
	if err != nil {
		logger.Log.Error("json.Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = c.svc.SaveCredentials(ctx, cred)
	if err != nil {
		logger.Log.Error("svc.SaveCredentials", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (c *Controller) HandleDeleteCredentialsByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	err := c.svc.DeleteCredentialsByID(ctx, id)
	if err != nil {
		logger.Log.Error("svc.DeleteCredentialsByID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (c *Controller) HandleGetUserCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	creds, err := c.svc.FindCredentialsByUserID(ctx)
	if err != nil {
		logger.Log.Error("svc.FindCredentialsByUserID", zap.Error(err))
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
	cred, err := c.svc.FindCredentialsByID(ctx, id)
	if err != nil {
		logger.Log.Error("svc.FindCredentialsByID", zap.Error(err))
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
	var credSync model2.CredSync
	err = json.Unmarshal(body, &credSync)
	if err != nil {
		logger.Log.Error("json.Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	credentials, err := c.svc.SyncCredentials(ctx, &credSync)
	if err != nil {
		logger.Log.Error("svc.SyncCredentials", zap.Error(err))
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
