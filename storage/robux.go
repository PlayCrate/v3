package storage

import (
	"context"

	"github.com/kattah7/v3/models"
)

func (s *PostgresStore) GetRobux() (*models.PlayerDataResponse, error) {
	fullResponse := &models.PlayerDataResponse{
		Other: make([]*models.Account, 0),
	}

	GetRows := func() ([]*models.Account, error) {
		query := `SELECT robloxId, robloxName, robux, time_saved FROM players ORDER BY robux DESC LIMIT $1`
		rows, err := s.db.Query(context.Background(), query, LIMIT)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		accounts := make([]*models.Account, 0)
		for rows.Next() {
			account := &models.Account{}
			if err := rows.Scan(&account.ID, &account.Name, &account.Robux, &account.LastSavedTime); err != nil {
				return nil, err
			}
			accounts = append(accounts, account)
		}
		return accounts, nil
	}

	Robux, _ := GetRows()
	fullResponse.Other = Robux

	return fullResponse, nil
}
