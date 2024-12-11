package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/IskanderSh/tages-task/internal/models"
	"github.com/IskanderSh/tages-task/pkg/errorlist"
	"github.com/IskanderSh/tages-task/pkg/utils"
)

type Service struct {
	log         *slog.Logger
	fileStorage FileStorage
	metaStorage MetaStorage
}

func NewService(
	log *slog.Logger,
	fileStorage FileStorage,
	metaStorage MetaStorage,
) *Service {
	return &Service{
		log:         log,
		fileStorage: fileStorage,
		metaStorage: metaStorage,
	}
}

type FileStorage interface {
	CreateFile(name string) (path string, err error)
	SaveFileChunk(name string, data []byte) error
	CloseFile(name string) error
}

type MetaStorage interface {
	GetByName(ctx context.Context, filename string) (*models.MetaInfo, error)
	Create(ctx context.Context, value models.MetaInfo) error
}

func (s *Service) UploadFile(
	ctx context.Context,
	req models.UploadFileRequest,
	counter int,
) (fileName string, err error) {
	name := fileName

	// need to create new file, if this is first chunk of gRPC stream
	if counter == 0 {
		name, err = s.createNewFile(ctx, req.FileName)
		if err != nil {
			return "", err
		}
	}

	if err = s.fileStorage.SaveFileChunk(name, req.Content); err != nil {
		return "", err
	}

	return name, nil
}

func (s *Service) FinishUpload(name string) error {
	return s.fileStorage.CloseFile(name)
}

// CreateNewFile creates file on disk, it will create file with new suffix id,
// if file with such name is already exists
func (s *Service) createNewFile(ctx context.Context, fileName string) (newName string, err error) {
	var name = fileName

	metaInfo, err := s.metaStorage.GetByName(ctx, fileName)
	if err != nil && !errors.Is(err, errorlist.ErrNotFound) {
		return "", err
	}

	// if file with such name is already exists, create file with suffix
	if metaInfo != nil {
		lastAddedID, found := utils.UniqueIDFromFileName(metaInfo.FileName)
		if found {
			name = fmt.Sprintf("%s (%d)", name, lastAddedID)
		}
	}

	path, err := s.fileStorage.CreateFile(metaInfo.FileName)
	if err != nil {
		return "", err
	}

	// set new name and path for new file
	metaInfo.FileName = name
	metaInfo.Path = path

	// save meta info for new file
	if err = s.metaStorage.Create(ctx, *metaInfo); err != nil {
		return "", err
	}

	return name, nil
}
