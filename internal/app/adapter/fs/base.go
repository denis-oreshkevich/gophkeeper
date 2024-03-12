package fs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain/service"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type BaseRepository[T domain.Base] struct {
	filename     string
	replFilename string
	base         T
	mx           sync.Mutex
}

func NewBaseRepository[T domain.Base](filename string, replFilename string,
	base T) *BaseRepository[T] {
	return &BaseRepository[T]{
		filename:     filename,
		replFilename: replFilename,
		base:         base,
	}
}

func (r *BaseRepository[T]) save(ctx context.Context, bs T) error {
	bytes, err := json.Marshal(bs)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}
	r.mx.Lock()
	defer r.mx.Unlock()

	file, err := os.OpenFile(r.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %w", err)
	}
	defer file.Close()

	if !bs.IsNew() {

		scanner := bufio.NewScanner(file)
		searchStr := fmt.Sprintf(`"id":"%s"`, bs.GetID())

		tmp, err := os.CreateTemp("", r.replFilename)
		if err != nil {
			return fmt.Errorf("os.CreateTemp: %w", err)
		}
		defer tmp.Close()

		for scanner.Scan() {
			text := scanner.Text()
			if strings.Contains(text, searchStr) {
				var bs domain.Base
				if err := json.Unmarshal(scanner.Bytes(), &bs); err != nil {
					return fmt.Errorf("json.Unmarshal: %w", err)
				}
				if bs.GetStatus() == domain.StatusActive {
					text = strings.Replace(text, text, "", 1)
				}
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
		err = os.Rename(tmp.Name(), r.filename)
		if err != nil {
			return fmt.Errorf("os.Rename: %w", err)
		}
	} else {
		if _, err = file.Write(bytes); err != nil {
			return fmt.Errorf("file.Write: %w", err)
		}
		if _, err = file.WriteString("\n"); err != nil {
			return fmt.Errorf("file.WriteString: %w", err)
		}
	}
	return nil
}

func (r *BaseRepository[T]) findByID(ctx context.Context, id string) (T, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	file, err := os.OpenFile(r.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		var t T
		return t, fmt.Errorf("os.OpenFile: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	searchStr := fmt.Sprintf(`"id":"%s"`, id)

	var found bool
	var bs T
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, searchStr) {
			err := json.Unmarshal(scanner.Bytes(), &bs)
			if err != nil {
				var t T
				return t, fmt.Errorf("json.Unmarshal: %w", err)
			}
			found = true
			break
		}
	}
	if err := scanner.Err(); err != nil {
		var t T
		return t, fmt.Errorf("scanner.Err: %w", err)
	}
	if !found {
		var t T
		return t, service.ErrItemNotFound
	}
	return bs, nil
}

func (r *BaseRepository[T]) findByUserID(ctx context.Context,
	userID string) ([]T, error) {

	r.mx.Lock()
	defer r.mx.Unlock()

	file, err := os.OpenFile(r.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var res []T
	for scanner.Scan() {
		var bs T
		err := json.Unmarshal(scanner.Bytes(), &bs)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}
		res = append(res, bs)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner.Err: %w", err)
	}
	return res, nil
}

func (r *BaseRepository[T]) findActiveModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]T, error) {

	r.mx.Lock()
	defer r.mx.Unlock()

	file, err := os.OpenFile(r.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var res []T
	for scanner.Scan() {
		var bs T
		err := json.Unmarshal(scanner.Bytes(), &bs)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}
		if bs.GetModifiedTms().After(tms) && bs.GetStatus() == domain.StatusActive {
			res = append(res, bs)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner.Err: %w", err)
	}
	return res, nil
}

func (r *BaseRepository[T]) findDeletedModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]T, error) {

	r.mx.Lock()
	defer r.mx.Unlock()

	file, err := os.OpenFile(r.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var res []T
	for scanner.Scan() {
		var bs T
		err := json.Unmarshal(scanner.Bytes(), &bs)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}
		if bs.GetModifiedTms().After(tms) && bs.GetStatus() == domain.StatusDeleted {
			res = append(res, bs)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner.Err: %w", err)
	}
	return res, nil
}

func (r *BaseRepository[T]) deleteByID(ctx context.Context, id string) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	file, err := os.OpenFile(r.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	searchStr := fmt.Sprintf(`"id":"%s"`, id)

	tmp, err := os.CreateTemp("", r.replFilename)
	if err != nil {
		return fmt.Errorf("os.CreateTemp: %w", err)
	}
	defer tmp.Close()

	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, searchStr) {
			var bs T
			if err := json.Unmarshal(scanner.Bytes(), &bs); err != nil {
				return fmt.Errorf("json.Unmarshal: %w", err)
			}
			bs.SetStatus(domain.StatusDeleted)
			bytes, err := json.Marshal(bs)
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
	err = os.Rename(tmp.Name(), r.filename)
	if err != nil {
		return fmt.Errorf("os.Rename: %w", err)
	}
	return nil
}
