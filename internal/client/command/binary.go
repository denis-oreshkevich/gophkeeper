package command

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/client/service"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/config"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/logger"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
	"github.com/google/uuid"
	"io"
	"os"
)

func DoFile(ctx context.Context, conf *config.Config, clientService *service.ClientService,
	user model.User) error {
	switch conf.Action {
	case config.ActionGet:
		byID, fErr := clientService.FindBinaryByID(ctx, conf.ID)
		if fErr != nil {
			return fmt.Errorf("clientService.FindBinaryByID: %w", fErr)
		}

		dec, bsErr := base64.StdEncoding.DecodeString(byID.Data)
		if bsErr != nil {
			return fmt.Errorf("base64.StdEncoding.DecodeString: %w", bsErr)
		}

		f, osErr := os.Create(byID.Name)
		if osErr != nil {
			return fmt.Errorf("os.Create: %w", osErr)
		}
		defer f.Close()

		if _, wErr := f.Write(dec); wErr != nil {
			return fmt.Errorf("f.Write: %w", wErr)
		}
		if sErr := f.Sync(); sErr != nil {
			return fmt.Errorf("f.Sync: %w", sErr)
		}
		logger.Log.Info("Success")
	case config.ActionSave:
		file, err := os.Open(conf.Filename)
		if err != nil {
			return fmt.Errorf("os.Open: %w", err)
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			return fmt.Errorf("file.Stat: %w", err)
		}

		bf := make([]byte, stat.Size())
		_, err = bufio.NewReader(file).Read(bf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("bufio.NewReader(file).Read: %w", err)
		}
		bStr := base64.StdEncoding.EncodeToString(bf)
		id := uuid.New()
		binary := model.Binary{
			ID:          id.String(),
			Name:        conf.Filename,
			Data:        bStr,
			New:         conf.IsNew,
			UserID:      user.ID,
			Status:      model.StatusActive,
			ModifiedTms: stat.ModTime().UTC(),
		}
		err = clientService.SaveBinary(ctx, &binary)
		if err != nil {
			return fmt.Errorf("clientService.SaveBinary: %w", err)
		}
		logger.Log.Info(fmt.Sprintf("saved file id = %s", id.String()))
	case config.ActionDelete:
		err := clientService.DeleteBinaryByID(ctx, conf.ID)
		if err != nil {
			return fmt.Errorf("clientService.DeleteBinaryByID: %w", err)
		}
		logger.Log.Info("Success")
	default:
		return fmt.Errorf("action value %s is unsupported", conf.Action)
	}
	return nil
}
