package storage

import (
	"context"
	"fmt"

	"github.com/kattah7/v3/models"
)

func (s *PostgresStore) InsertGhostHunt(account *models.GhostHuntAccount) (*models.GhostHuntSerial, error) {
	if account.RobloxID == 0 {
		return nil, fmt.Errorf("robloxId cannot be empty")
	}

	// Check the current row count in the table
	countQuery := "SELECT count(*) FROM ghost_hunt_top_25"
	var rowCount int
	err := s.db.QueryRow(context.Background(), countQuery).Scan(&rowCount)
	if err != nil {
		return nil, err
	}

	if rowCount >= 25 {
		return nil, fmt.Errorf("row limit reached. Cannot insert more than 25 rows")
	}

	// If the row count is less than 25, proceed with the insertion
	insertQuery := `INSERT INTO ghost_hunt_top_25 (robloxId) VALUES ($1)`
	_, err = s.db.Exec(context.Background(), insertQuery, account.RobloxID)
	if err != nil {
		return nil, err
	}

	return &models.GhostHuntSerial{Serial: int64(rowCount + 1)}, nil
}
