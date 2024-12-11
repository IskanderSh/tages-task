package storage

import (
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/IskanderSh/tages-task/pkg/errorlist"
)

type FileStorage struct {
	log       *slog.Logger
	dirName   string
	openFiles map[string]*os.File
	mu        *sync.Mutex
}

func NewFileStorage(log *slog.Logger, dirName string) *FileStorage {
	return &FileStorage{
		log:       log,
		dirName:   dirName,
		openFiles: make(map[string]*os.File),
		mu:        &sync.Mutex{},
	}
}

func (s *FileStorage) CreateFile(name string) (path string, err error) {
	path = filepath.Join(s.dirName, name)
	s.log.Debug("creating file with path: ", path)

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	s.openFiles[name] = file
	s.mu.Unlock()

	return path, nil
}

func (s *FileStorage) SaveFileChunk(name string, data []byte) error {
	s.mu.Lock()
	file, exists := s.openFiles[name]
	if !exists {
		return errorlist.ErrNoFileWithSuchName
	}
	s.mu.Unlock()

	_, err := file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *FileStorage) CloseFile(name string) error {
	s.mu.Lock()
	file, exists := s.openFiles[name]
	if !exists {
		return errorlist.ErrNoFileWithSuchName
	}
	s.mu.Unlock()

	return file.Close()
}