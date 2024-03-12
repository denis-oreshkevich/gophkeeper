package server

import (
	"encoding/json"
	"errors"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain/service"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type AuthUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ValidationErrEntry struct {
	Field  string   `json:"field"`
	Errors []string `json:"errors"`
}

func NewValidationErr(field string, errs []string) ValidationErrEntry {
	return ValidationErrEntry{
		Field:  field,
		Errors: errs,
	}
}

func (c *Controller) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, ok := readAndValidateUser(w, r)
	if !ok {
		logger.Log.Debug("readAndValidateUser is not ok")
		return
	}
	_, err := c.userClient.Register(ctx, u.Login, u.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExist) {
			logger.Log.Debug("register user", zap.Error(err))
			w.WriteHeader(http.StatusConflict)
			return
		}
		logger.Log.Error("register user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token, err := c.userClient.Login(ctx, u.Login, u.Password)
	if err != nil {
		logger.Log.Error("userClient.Login", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(AuthorizationHeaderName, token)
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, ok := readAndValidateUser(w, r)
	if !ok {
		logger.Log.Debug("readAndValidateUser is not ok")
		return
	}
	token, err := c.userClient.Login(ctx, u.Login, u.Password)
	if err != nil {
		logger.Log.Error("svc.Login", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(AuthorizationHeaderName, token)
	w.WriteHeader(http.StatusOK)
}

func readAndValidateUser(w http.ResponseWriter, r *http.Request) (AuthUser, bool) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("io.ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return AuthUser{}, false
	}
	var u AuthUser
	if err = json.Unmarshal(body, &u); err != nil {
		logger.Log.Error("json.Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return AuthUser{}, false
	}
	valResp, err := validateUser(u)
	if err != nil {
		logger.Log.Debug("u is not valid", zap.Error(err))
		bytes, err := json.Marshal(valResp)
		if err != nil {
			logger.Log.Error("validation json.Marshal", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return AuthUser{}, false
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bytes)
		return AuthUser{}, false
	}
	return u, true
}

var ErrUserValidation = errors.New("user validation error")

func validateUser(user AuthUser) ([]ValidationErrEntry, error) {
	login := strings.Trim(user.Login, " ")
	isEV := len(login) > 0
	pswd := strings.Trim(user.Password, " ")
	isPV := len(pswd) > 0
	var valErrors = make([]ValidationErrEntry, 0)
	if !isEV {
		validationErr := NewValidationErr("login", []string{"login is not valid"})
		valErrors = append(valErrors, validationErr)
	}
	if !isPV {
		validationErr := NewValidationErr("password", []string{"password is not valid"})
		valErrors = append(valErrors, validationErr)
	}
	if len(valErrors) == 0 {
		return nil, nil
	}

	return valErrors, ErrUserValidation
}
