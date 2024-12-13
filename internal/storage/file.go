package storage

import (
	"fmt"
	"log/slog"
	"os"
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
	path = fmt.Sprintf("%s/%s", s.dirName, name)
	s.log.Info("file with path:", path)

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
	value, exists := s.openFiles[name]
	if !exists {
		return errorlist.ErrNoFileWithSuchName
	}

	delete(s.openFiles, name)
	s.mu.Unlock()

	return value.Close()
}

func (s *FileStorage) OpenFile(name string) error {
	path := fmt.Sprintf("%s/%s", s.dirName, name)
	s.log.Info("file with path:", path)

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.openFiles[name] = file
	s.mu.Unlock()

	return nil
}

func (s *FileStorage) ReadFileChunk(name string, buffer []byte) (bytesRead int, err error) {
	file, exists := s.openFiles[name]
	if !exists {
		return bytesRead, errorlist.ErrNoFileWithSuchName
	}

	bytesRead, err = file.Read(buffer)
	if err != nil {
		return bytesRead, err
	}

	return bytesRead, nil
}
