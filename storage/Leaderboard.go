package storage

import (
	"context"

	"github.com/kattah7/v3/models"
)

func (s *PostgresStore) GetSpecificPlayer(robloxId int64) (*models.Account, error) {
	query := `SELECT robloxId, robloxName, bubbles, time_saved FROM players WHERE robloxId = $1`
	row := s.db.QueryRow(context.Background(), query, robloxId)

	account := &models.Account{}
	if err := row.Scan(&account.ID, &account.Name, &account.Bubbles, &account.LastSavedTime); err != nil {
		return nil, err
	}
	return account, nil
}
