package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

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
	ReadFileChunk(name string, buffer []byte) (bytesRead int, err error)
	OpenFile(name string) error
}

type MetaStorage interface {
	GetByName(ctx context.Context, filename string) (*models.MetaInfo, error)
	Create(ctx context.Context, value models.MetaInfo) error
	FetchAll(ctx context.Context) (data []models.MetaInfo, err error)
}

func (s *Service) UploadFile(
	ctx context.Context,
	req models.UploadFileRequest,
	counter int,
) (fileName string, err error) {
	name := req.FileName

	// need to create new file, if this is first chunk of gRPC stream
	if counter == 0 {
		name, err = s.createNewFile(ctx, req.FileName)
		if err != nil {
			s.log.Error(err.Error())
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

func (s *Service) DownloadFile(name string, counter int, buffer []byte) (bytesRead int, err error) {
	// open file if it is first chunk
	if counter == 0 {
		if err = s.fileStorage.OpenFile(name); err != nil {
			return bytesRead, err
		}
	}

	return s.fileStorage.ReadFileChunk(name, buffer)
}

func (s *Service) GetFiles(ctx context.Context) ([]models.MetaInfo, error) {
	return s.metaStorage.FetchAll(ctx)
}

// CreateNewFile creates file on disk, it will create file with new suffix id,
// if file with such name is already exists
func (s *Service) createNewFile(ctx context.Context, fileName string) (newName string, err error) {
	var name = fileName

	metaInfo, err := s.metaStorage.GetByName(ctx, fileName)
	s.log.Info("metaInfo:", metaInfo)
	if err != nil && !errors.Is(err, errorlist.ErrNotFound) {
		return "", err
	}

	// if file with such name is already exists, create file with suffix
	if metaInfo != nil {
		name = s.nextUniqueFileName(metaInfo.FileName)
	}

	path, err := s.fileStorage.CreateFile(name)
	if err != nil {
		return "", err
	}

	// set new name and path for new file
	metaInfo = &models.MetaInfo{
		FileName:  name,
		Path:      path,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// save meta info for new file
	if err = s.metaStorage.Create(ctx, *metaInfo); err != nil {
		return "", err
	}

	return name, nil
}

func (s *Service) nextUniqueFileName(name string) string {
	nextName := utils.NextUniqueName(name)
	// service has file with same name, but couldn't define new unique id for file
	if nextName == name {
		s.log.Warn("error defining name for new file with existing name: ", name)
	}

	s.log.Info("getting new file name for last one", nextName, name)

	return nextName
}
