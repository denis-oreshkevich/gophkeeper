package fs

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain/service"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/logger"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type ClientRepository struct {
	filename string
	mx       sync.Mutex
}

func NewClientRepository(filename string) *ClientRepository {
	return &ClientRepository{
		filename: filename,
	}
}

func (r *ClientRepository) FindClient(ctx context.Context) (domain.Client, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	file, err := os.OpenFile(r.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return domain.Client{}, fmt.Errorf("os.OpenFile: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var found bool
	var cl domain.Client
	for scanner.Scan() {
		text := scanner.Text()
		if text != "" {
			err := json.Unmarshal(scanner.Bytes(), &cl)
			if err != nil {
				return domain.Client{}, fmt.Errorf("json.Unmarshal: %w", err)
			}
			found = true
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return domain.Client{}, fmt.Errorf("scanner.Err: %w", err)
	}
	if !found {
		return domain.Client{}, service.ErrItemNotFound
	}
	return cl, nil
}

func (r *Repository) CreateClient(ctx context.Context, client domain.Client) (domain.Client, error) {
	r.clientRepo.mx.Lock()
	defer r.clientRepo.mx.Unlock()

	var create bool
	if _, err := os.Stat(r.userRepo.filename); errors.Is(err, os.ErrNotExist) {
		logger.Log.Info("user file is not exist")
		create = true
	}
	if create {
		file, err := os.OpenFile(r.userRepo.filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return domain.Client{}, fmt.Errorf("os.OpenFile: %w", err)
		}
		defer file.Close()

		bytes, err := json.Marshal(client)
		if err != nil {
			return domain.Client{}, fmt.Errorf("json.Unmarshal: %w", err)
		}
		_, err = file.Write(bytes)
		if err != nil {
			return domain.Client{}, fmt.Errorf("file.Write: %w", err)
		}
		if _, err = file.WriteString("\n"); err != nil {
			return domain.Client{}, fmt.Errorf("file.WriteString: %w", err)
		}
	}
	return client, nil
}
func (r *Repository) UpdateClientLastSyncTmsByID(ctx context.Context, id string, syncTms time.Time) error {
	r.clientRepo.mx.Lock()
	defer r.clientRepo.mx.Unlock()

	file, err := os.OpenFile(r.clientRepo.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	searchStr := fmt.Sprintf(`"id":"%s"`, id)

	tmp, err := os.CreateTemp("", r.clientRepo.filename)
	if err != nil {
		return fmt.Errorf("os.CreateTemp: %w", err)
	}
	defer tmp.Close()

	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, searchStr) {
			var cl domain.Client
			if err := json.Unmarshal(scanner.Bytes(), &cl); err != nil {
				return fmt.Errorf("json.Unmarshal: %w", err)
			}
			cl.SyncTms = syncTms
			bytes, err := json.Marshal(cl)
			if err != nil {
				return fmt.Errorf("json.Marshal: %w", err)
			}
			text = strings.Replace(text, text, string(bytes), 1)
		}
		_, err := io.WriteString(tmp, text)
		if err != nil {
			return fmt.Errorf("io.WriteString: %w", err)
		}

		if _, err = file.WriteString("\n"); err != nil {
			return fmt.Errorf("file.WriteString: %w", err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner.Err: %w", err)
	}
	err = os.Rename(tmp.Name(), r.clientRepo.filename)
	if err != nil {
		return fmt.Errorf("os.Rename: %w", err)
	}
	return nil
}
func (r *Repository) FindClientByID(ctx context.Context, id string) (domain.Client, error) {
	r.clientRepo.mx.Lock()
	defer r.clientRepo.mx.Unlock()

	file, err := os.OpenFile(r.clientRepo.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return domain.Client{}, fmt.Errorf("os.OpenFile: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	searchStr := fmt.Sprintf(`"id":"%s"`, id)

	var found bool
	var cl domain.Client
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, searchStr) {
			err := json.Unmarshal(scanner.Bytes(), &cl)
			if err != nil {
				return domain.Client{}, fmt.Errorf("json.Unmarshal: %w", err)
			}
			found = true
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return domain.Client{}, fmt.Errorf("scanner.Err: %w", err)
	}
	if !found {
		return domain.Client{}, service.ErrItemNotFound
	}
	return cl, nil
}
func (r *Repository) FindClientsByUserID(ctx context.Context,
	userID string) ([]domain.Client, error) {
	r.clientRepo.mx.Lock()
	defer r.clientRepo.mx.Unlock()

	file, err := os.OpenFile(r.clientRepo.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	searchStr := fmt.Sprintf(`"userID":"%s"`, userID)

	var found bool
	var res []domain.Client
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, searchStr) {
			var cl domain.Client
			err := json.Unmarshal(scanner.Bytes(), &cl)
			if err != nil {
				return nil, fmt.Errorf("json.Unmarshal: %w", err)
			}
			found = true
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner.Err: %w", err)
	}
	if !found {
		return nil, service.ErrItemNotFound
	}
	return res, nil
}
