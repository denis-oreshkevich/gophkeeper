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

func (c *Controller) HandlePostCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("io.ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var card domain.Card
	err = json.Unmarshal(body, &card)
	if err != nil {
		logger.Log.Error("json.Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = c.crudSvc.SaveCard(ctx, card)
	if err != nil {
		logger.Log.Error("crudSvc.SaveCard", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (c *Controller) HandleDeleteCardByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	err := c.crudSvc.DeleteCardByID(ctx, id)
	if err != nil {
		logger.Log.Error("crudSvc.DeleteCardByID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (c *Controller) HandleGetUserCards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cards, err := c.crudSvc.FindCardsByUserID(ctx)
	if err != nil {
		logger.Log.Error("crudSvc.FindCardsByUserID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(cards) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	result, err := json.Marshal(cards)
	if err != nil {
		logger.Log.Error("json.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (c *Controller) HandleGetCardByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	card, err := c.crudSvc.FindCardByID(ctx, id)
	if err != nil {
		logger.Log.Error("crudSvc.FindCardByID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(card)
	if err != nil {
		logger.Log.Error("json.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (c *Controller) HandlePostSyncCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("io.ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var cardSync domain.CardSync
	err = json.Unmarshal(body, &cardSync)
	if err != nil {
		logger.Log.Error("json.Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(cardSync.Cards) == 0 {
		log.Debug("cardSync.Cards length for sync = 0")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	cards, err := c.syncSvc.SyncCard(ctx, &cardSync)
	if err != nil {
		logger.Log.Error("syncSvc.SyncCard", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(cards)
	if err != nil {
		logger.Log.Error("json.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(result)
}
