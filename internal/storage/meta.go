package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/IskanderSh/tages-task/config"
	"github.com/IskanderSh/tages-task/internal/models"
	"github.com/IskanderSh/tages-task/pkg/errorlist"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type MetaStorage struct {
	log *slog.Logger
	db  *sqlx.DB
}

func NewMetaStorage(log *slog.Logger, cfg config.MetaStorage) (*MetaStorage, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &MetaStorage{
		log: log,
		db:  db,
	}, nil
}

func (s *MetaStorage) GetByName(ctx context.Context, filename string) (*models.MetaInfo, error) {
	var result models.MetaInfo
	err := s.db.QueryRowContext(ctx, queryGetByID, filename).
		Scan(&result.ID, &result.FileName, &result.Path, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorlist.ErrNotFound
		}
		return nil, err
	}

	return &result, nil
}

func (s *MetaStorage) Create(ctx context.Context, value models.MetaInfo) error {
	_, err := s.db.ExecContext(
		ctx,
		queryCreate,
		value.FileName,
		value.Path,
		value.CreatedAt,
		value.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *MetaStorage) FetchAll(ctx context.Context) (data []models.MetaInfo, err error) {
	rows, err := s.db.QueryContext(ctx, queryFetchAll)
	if err != nil {
		return data, err
	}

	for rows.Next() {
		var value models.MetaInfo
		if err = rows.Scan(&value.ID, &value.FileName, &value.Path, &value.CreatedAt, &value.UpdatedAt); err != nil {
			s.log.Warn("error scanning values to meta info:", err)
			continue
		}

		data = append(data, value)
	}

	return data, nil
}
