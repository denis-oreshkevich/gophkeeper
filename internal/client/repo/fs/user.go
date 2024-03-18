package fs

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/client/repo"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/logger"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
	"os"
	"sync"
)

type UserRepository struct {
	filename string
	mx       sync.Mutex
}

func NewUserRepository(filename string) *UserRepository {
	return &UserRepository{
		filename: filename,
	}
}

func (r *Repository) FindUser(ctx context.Context) (model.User, error) {
	if _, err := os.Stat(r.userRepo.filename); errors.Is(err, os.ErrNotExist) {
		logger.Log.Info("user file is not exist")
		return model.User{}, repo.ErrItemNotFound
	}

	file, err := os.OpenFile(r.userRepo.filename, os.O_RDONLY, 0666)
	if err != nil {
		return model.User{}, fmt.Errorf("os.OpenFile: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var usr model.User
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &usr)
		if err != nil {
			return model.User{}, fmt.Errorf("json.Unmarshal: %w", err)
		}
	}
	if err := scanner.Err(); err != nil {
		return model.User{}, fmt.Errorf("scanner.Err: %w", err)
	}
	return usr, nil
}

func (r *Repository) CreateUser(ctx context.Context, usr model.User) (model.User, error) {
	r.userRepo.mx.Lock()
	defer r.userRepo.mx.Unlock()

	var create bool
	if _, err := os.Stat(r.userRepo.filename); errors.Is(err, os.ErrNotExist) {
		logger.Log.Info("user file is not exist")
		create = true
	}
	if !create {
		return model.User{}, repo.ErrUserAlreadyExist
	}
	file, err := os.OpenFile(r.userRepo.filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return model.User{}, fmt.Errorf("os.OpenFile: %w", err)
	}
	defer file.Close()

	bytes, err := json.Marshal(usr)
	if err != nil {
		return model.User{}, fmt.Errorf("json.Unmarshal: %w", err)
	}
	_, err = file.Write(bytes)
	if err != nil {
		return model.User{}, fmt.Errorf("file.Write: %w", err)
	}
	if _, err = file.WriteString("\n"); err != nil {
		return model.User{}, fmt.Errorf("file.WriteString: %w", err)
	}
	return usr, nil
}

func (r *Repository) FindUserByLogin(ctx context.Context, login string) (model.User, error) {
	r.userRepo.mx.Lock()
	defer r.userRepo.mx.Unlock()

	var notExist bool
	if _, err := os.Stat(r.userRepo.filename); errors.Is(err, os.ErrNotExist) {
		logger.Log.Info("user file is not exist")
		notExist = true
	}
	if notExist {
		return model.User{}, repo.ErrItemNotFound
	}

	file, err := os.OpenFile(r.userRepo.filename, os.O_RDONLY, 0666)
	if err != nil {
		return model.User{}, fmt.Errorf("os.OpenFile: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var usr model.User
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &usr)
		if err != nil {
			return model.User{}, fmt.Errorf("json.Unmarshal: %w", err)
		}
	}
	if err := scanner.Err(); err != nil {
		return model.User{}, fmt.Errorf("scanner.Err: %w", err)
	}
	return usr, nil
}
