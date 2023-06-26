package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kattah7/v3/models"
)

const LIMIT = 100

type Storage interface {
	GetSecrets() (*models.PlayerDataResponse, error)
	GetEggs() (*models.PlayerDataResponse, error)
	GetBubbles() (*models.PlayerDataResponse, error)
	GetPower() (*models.PlayerDataResponse, error)
	GetRobux() (*models.PlayerDataResponse, error)
	GetPlaytime() (*models.PlayerDataResponse, error)

	InsertAccounts(*models.Account) error
	Close()

	ListAuction(*models.AuctionAccount) error
	RemoveAuction(*models.AuctionAccount) error
	GetAuctions() ([]*models.AuctionAccount, error)
	PurchaseAuction(*models.AuctionAccount) error
	GetAuctionClaims(*models.AuctionAccount) ([]*models.AuctionAccount, error)

	AuctionClaim(*models.AuctionAccount) error
	AuctionUnlist(*models.AuctionAccount) error

	GetAuctionListing(*models.AuctionAccount) ([]*models.AuctionAccount, error)
}

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
	queries := []string{
		`CREATE TABLE IF NOT EXISTS players (
			id SERIAL PRIMARY KEY,
			robloxId BIGINT NOT NULL UNIQUE,
			robloxName VARCHAR(255) NOT NULL,
			secrets BIGINT NOT NULL DEFAULT 0,
			eggs BIGINT NOT NULL DEFAULT 0,
			bubbles BIGINT NOT NULL DEFAULT 0,
			power BIGINT NOT NULL DEFAULT 0,
			robux BIGINT NOT NULL DEFAULT 0,
			playtime BIGINT NOT NULL DEFAULT 0,
			time_saved TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS auctions (
			id SERIAL PRIMARY KEY,
			robloxId BIGINT NOT NULL,
			robloxName VARCHAR(255) NOT NULL,
			itemType VARCHAR(255) NOT NULL,
			itemData JSONB NOT NULL,
			startPrice BIGINT NOT NULL,
			priceType VARCHAR(255) NOT NULL,
			status VARCHAR(255) NOT NULL DEFAULT 'OPEN'
		)`,
	}

	for _, query := range queries {
		_, err := s.db.Exec(context.Background(), query)
		if err != nil {
			return err
		}
	}

	return nil
}
