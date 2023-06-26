package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kattah7/v3/models"
	"github.com/robfig/cron"
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

	AuctionExpireList(*models.AuctionAccount) ([]*models.AuctionAccount, error)
	AuctionExpireClaim(*models.AuctionAccount) error
}

type PostgresStore struct {
	cfg *models.Config
	db  *pgxpool.Pool
}

var (
	pgInstance *PostgresStore
	pgOnce     sync.Once
)

func (s *PostgresStore) Close() {
	s.db.Close()
}

func NewPostgresStore(ctx context.Context, cfg *models.Config) (*PostgresStore, error) {
	var err error

	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, cfg.DBConnString)
		if err != nil {
			err = fmt.Errorf("unable to connect to database: %v", err)
			return
		}

		pgInstance = &PostgresStore{
			db:  db,
			cfg: cfg,
		}
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
			listed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			status VARCHAR(255) NOT NULL DEFAULT 'OPEN'
		)`,
		`CREATE TABLE IF NOT EXISTS auction_expired (
			id SERIAL PRIMARY KEY,
			robloxId BIGINT NOT NULL,
			robloxName VARCHAR(255) NOT NULL,
			itemType VARCHAR(255) NOT NULL,
			itemData JSONB NOT NULL
		)`,
	}

	for _, query := range queries {
		_, err := s.db.Exec(context.Background(), query)
		if err != nil {
			return err
		}
	}

	c := cron.New()
	c.AddFunc("* * * * *", func() {
		currentTime := time.Now().Local()
		cutoffDuration := time.Duration(s.cfg.CutOffTime) * time.Second
		cutoffTime := currentTime.Add(cutoffDuration)

		fmt.Println(cutoffTime)

		moveQuery := `
			INSERT INTO auction_expired (robloxId, robloxName, itemType, itemData)
			SELECT robloxId, robloxName, itemType, itemData
			FROM auctions
			WHERE listed < $1 AND status = 'OPEN'
		`
		_, err := s.db.Exec(context.Background(), moveQuery, cutoffTime)
		if err != nil {
			fmt.Println("Failed to move data from auctions to auction_expired:", err)
			return
		}

		deleteQuery := `
			DELETE FROM auctions
			WHERE listed < $1 AND status = 'OPEN'
		`
		_, err = s.db.Exec(context.Background(), deleteQuery, cutoffTime)
		if err != nil {
			fmt.Println("Failed to delete rows from auctions:", err)
			return
		}
	})

	c.Start()

	return nil
}
