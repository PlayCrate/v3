package storage

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kattah7/v3/models"
)

const LIMIT = 100

type Storage interface {
	GetSecrets() (*PlayerDataResponse, error)
	GetEggs() (*PlayerDataResponse, error)
	GetBubbles() (*PlayerDataResponse, error)
	GetPower() (*PlayerDataResponse, error)
	GetRobux() (*PlayerDataResponse, error)
	GetPlaytime() (*PlayerDataResponse, error)

	InsertAccounts(*models.Account) error
	Close()
}

// type Player map[string]*PlayerData

type PlayerDataResponse struct {
	F2P    []*models.Account `json:"f2p,omitempty"`
	NonF2P []*models.Account `json:"nof2p,omitempty"`
	Other  []*models.Account `json:"other,omitempty"`
}

// type PlayerData struct {
// 	Daily   []*models.Account `json:"daily"`
// 	Weekly  []*models.Account `json:"weekly"`
// 	Monthly []*models.Account `json:"monthly"`
// 	All     []*models.Account `json:"all"`
// }

type PostgresStore struct {
	db *pgxpool.Pool
}

var (
	pgInstance *PostgresStore
	pgOnce     sync.Once
)

func (s *PostgresStore) Close() {
	s.db.Close()
}

func NewPostgresStore(ctx context.Context, connString string) (*PostgresStore, error) {
	var err error

	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			err = fmt.Errorf("unable to connect to database: %v", err)
			return
		}

		pgInstance = &PostgresStore{db}
	})

	if err != nil {
		return nil, err
	}

	return pgInstance, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateTables()
}

func (s *PostgresStore) CreateTables() error {
	query :=
		`CREATE TABLE IF NOT EXISTS players (
			id SERIAL PRIMARY KEY,
			robloxId INT NOT NULL UNIQUE,
			robloxName VARCHAR(255) NOT NULL,
			secrets INT NOT NULL DEFAULT 0,
			eggs INT NOT NULL DEFAULT 0,
			bubbles INT NOT NULL DEFAULT 0,
			power INT NOT NULL DEFAULT 0,
			robux INT NOT NULL DEFAULT 0,
			playtime INT NOT NULL DEFAULT 0,
			time_saved TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`

	_, err := s.db.Exec(context.Background(), query)

	return err
}

func (s *PostgresStore) InsertAccounts(acc *models.Account) error {
	query := `
    INSERT INTO players (robloxId, robloxName, secrets, eggs, bubbles, power, robux, playtime, time_saved)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    ON CONFLICT (robloxId) DO UPDATE SET
        robloxName = EXCLUDED.robloxName,
        secrets = EXCLUDED.secrets,
        eggs = EXCLUDED.eggs,
        bubbles = EXCLUDED.bubbles,
        power = EXCLUDED.power,
        robux = EXCLUDED.robux,
        playtime = EXCLUDED.playtime,
        time_saved = EXCLUDED.time_saved
`

	_, err := s.db.Exec(context.Background(), query, acc.ID, acc.Name, acc.Secrets, acc.Eggs, acc.Bubbles, acc.Power, acc.Robux, acc.Playtime, acc.LastSavedTime)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}
	log.Println("Insert successful", acc.Name)

	return nil
}
