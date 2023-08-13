package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"io"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kattah7/v3/models"
	"github.com/robfig/cron"

	"github.com/redis/go-redis/v9"
)

const LIMIT = 100

type Storage interface {
	Close()

	GetSecrets() (*models.PlayerDataResponse, error)
	GetEggs() (*models.PlayerDataResponse, error)
	GetBubbles() (*models.PlayerDataResponse, error)
	GetPower() (*models.PlayerDataResponse, error)
	GetRobux() (*models.PlayerDataResponse, error)
	GetPlaytime() (*models.PlayerDataResponse, error)
	GetSpecificPlayer(int64) (*models.AccountLookup, error)
	InsertAccounts(*models.Account) error

	ListAuction(*models.AuctionAccount) error
	RemoveAuction(*models.AuctionAccount) error
	GetAuctions() ([]*models.AuctionAccount, error)
	PurchaseAuction(*models.AuctionAccount) error
	GetAuctionClaims(*models.AuctionAccount) ([]*models.AuctionAccount, error)
	AuctionClaim(*models.AuctionAccount) error
	AuctionUnlist(*models.AuctionAccount) error
	GetAuctionListing(*models.AuctionAccount) ([]*models.AuctionAccount, error)

	InsertPetsExistance(*models.PetsExistance) error
	GetPetsExistance() ([]*models.GetPetsExistance, error)
	DeletePetsExistence(*models.PetsExistance) error
}

type PostgresStore struct {
	cfg *models.Config
	db  *pgxpool.Pool
	rdb *redis.Client
}

var (
	pgInstance *PostgresStore
	pgOnce     sync.Once
)

func (s *PostgresStore) Close() {
	s.db.Close()
}

func NewPostgresStore(ctx context.Context, cfg *models.Config, rdb *redis.Client) (*PostgresStore, error) {
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
			rdb: rdb,
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
		`CREATE TABLE IF NOT EXISTS pets_exist (
			id SERIAL PRIMARY KEY,
			robloxId BIGINT NOT NULL,
			petId VARCHAR(255) NOT NULL,
			petCount BIGINT NOT NULL DEFAULT 0,
			CONSTRAINT uc_robloxid_petid UNIQUE (robloxId, petId)
		)`,
	}

	for _, query := range queries {
		_, err := s.db.Exec(context.Background(), query)
		if err != nil {
			return err
		}
	}

	c := cron.New()
	c.AddFunc(s.cfg.Cron, func() {
		currentTime := time.Now().Local()
		cutoffDuration := time.Duration(s.cfg.CutOffTime) * time.Second
		cutoffTime := currentTime.Add(cutoffDuration)

		query := `SELECT robloxId, robloxName, itemData FROM auctions WHERE listed < $1 AND status = 'OPEN'`

		rows, err := s.db.Query(context.Background(), query, cutoffTime)
		if err != nil {
			fmt.Errorf("unable to query database: %v", err)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var robloxId int64
			var robloxName string
			var itemData map[string]interface{}

			err := rows.Scan(&robloxId, &robloxName, &itemData)
			if err != nil {
				fmt.Errorf("unable to query database: %v", err)
				return
			}

			itemData["timestamp"] = time.Now().Unix()
			itemData["message"] = "This item has expired and has been returned to your mailbox."
			itemData["senderId"] = 1
			itemData["senderName"] = "PlayCrate"
			itemData["displayName"] = "PlayCrate"
			itemData["targetId"] = robloxId

			updatedItemData, err := json.Marshal(itemData)
			if err != nil {
				fmt.Errorf("unable to marshal itemData: %v", err)
				return
			}

			body := models.MailboxExpire{
				RobloxName: robloxName,
				RobloxId:   robloxId,
				Type:       "ADD",
				Payload:    json.RawMessage(fmt.Sprintf("[%s]", updatedItemData)),
			}

			jsonData, err := json.Marshal(body)
			if err != nil {
				fmt.Printf("Failed to marshal auction data: %v\n", err)
				return
			}

			var baseURL string
			if s.cfg.Prod {
				baseURL = "https://roblox.kattah.me/mailbox"
			} else {
				baseURL = "https://playcrate-debug.kattah.me/mailbox"
			}

			req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Printf("Failed to create HTTP request: %v\n", err)
				return
			}

			req.Header.Set("authorization", s.cfg.V1Auth)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Failed to send API request: %v\n", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Printf("API request failed with status code: %d\n", resp.StatusCode)
				return
			}

			bodya, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Failed to read API response body: %v\n", err)
				return
			}

			var apiResp models.ApiResponse
			err = json.Unmarshal(bodya, &apiResp)
			if err != nil {
				fmt.Printf("Failed to unmarshal API response body: %v\n", err)
				return
			}

			if apiResp.Success == true {
				deleteQuery := `DELETE FROM auctions WHERE listed < $1 AND status = 'OPEN'`
				_, err = s.db.Exec(context.Background(), deleteQuery, cutoffTime)
				if err != nil {
					fmt.Println("Failed to delete rows from auctions:", err)
					return
				}

				fmt.Println("Successfully sent expired auction to mailbox")
			} else {
				fmt.Println("Failed to send expired auction to mailbox")
			}
		}
	})

	cacheDB := func() {
		pets, err := s.GetPetsExistance()
		if err != nil {
			fmt.Println("Failed to get pets:", err)
		} else {
			if err := cacheData(s, "pets-exist", pets); err != nil {
				fmt.Println(err)
			}
		}

		// Cache eggs data
		eggslb, err := s.GetEggs()
		if err != nil {
			fmt.Println("Failed to get eggs:", err)
		} else {
			if err := cacheData(s, "eggs-lb", eggslb); err != nil {
				fmt.Println(err)
			}
		}

		// Cache secrets data
		secretslb, err := s.GetSecrets()
		if err != nil {
			fmt.Println("Failed to get secrets:", err)
		} else {
			if err := cacheData(s, "secrets-lb", secretslb); err != nil {
				fmt.Println(err)
			}
		}

		// Cache bubbles data
		bubbleslb, err := s.GetBubbles()
		if err != nil {
			fmt.Println("Failed to get bubbles:", err)
		} else {
			if err := cacheData(s, "bubbles-lb", bubbleslb); err != nil {
				fmt.Println(err)
			}
		}

		// Cache power data
		powerlb, err := s.GetPower()
		if err != nil {
			fmt.Println("Failed to get power:", err)
		} else {
			if err := cacheData(s, "power-lb", powerlb); err != nil {
				fmt.Println(err)
			}
		}

		// Cache robux data
		robuxlb, err := s.GetRobux()
		if err != nil {
			fmt.Println("Failed to get robux:", err)
		} else {
			if err := cacheData(s, "robux-lb", robuxlb); err != nil {
				fmt.Println(err)
			}
		}

		// Cache playtime data
		playtimelb, err := s.GetPlaytime()
		if err != nil {
			fmt.Println("Failed to get playtime:", err)
		} else {
			if err := cacheData(s, "playtime-lb", playtimelb); err != nil {
				fmt.Println(err)
			}
		}

		fmt.Println("Successfully updated pets and eggs")
	}

	c.AddFunc("@every 1m", func() {
		cacheDB()
	})

	go cacheDB()

	c.Start()

	return nil
}

func cacheData(s *PostgresStore, key string, data interface{}) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Failed to marshal data: %w", err)
	}

	err = s.rdb.Set(context.Background(), key, dataJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("Failed to set data: %w", err)
	}

	return nil
}
