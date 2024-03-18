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

func (c *Controller) HandlePostText(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("io.ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var text model2.Text
	err = json.Unmarshal(body, &text)
	if err != nil {
		logger.Log.Error("json.Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = c.svc.SaveText(ctx, &text)
	if err != nil {
		logger.Log.Error("svc.SaveText", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (c *Controller) HandleDeleteTextByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	err := c.svc.DeleteTextByID(ctx, id)
	if err != nil {
		logger.Log.Error("svc.DeleteTextByID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (c *Controller) HandleGetUserTexts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	texts, err := c.svc.FindTextsByUserID(ctx)
	if err != nil {
		logger.Log.Error("svc.FindTextsByUserID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(texts) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	result, err := json.Marshal(texts)
	if err != nil {
		logger.Log.Error("json.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (c *Controller) HandleGetTextByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	text, err := c.svc.FindTextByID(ctx, id)
	if err != nil {
		logger.Log.Error("svc.FindTextByID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(text)
	if err != nil {
		logger.Log.Error("json.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (c *Controller) HandlePostSyncText(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("io.ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var textSync model2.TextSync
	err = json.Unmarshal(body, &textSync)
	if err != nil {
		logger.Log.Error("json.Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	texts, err := c.svc.SyncText(ctx, &textSync)
	if err != nil {
		logger.Log.Error("svc.SyncText", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(texts)
	if err != nil {
		logger.Log.Error("json.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(result)
}
