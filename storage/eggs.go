package storage

import (
	"context"

	"github.com/kattah7/v3/models"
)

func (s *PostgresStore) GetEggs() (*PlayerDataResponse, error) {
	fullResponse := &PlayerDataResponse{
		F2P:    make([]*models.Account, 0),
		NonF2P: make([]*models.Account, 0),
	}

	GetRows := func(f2p bool) ([]*models.Account, error) {
		query := `SELECT robloxId, robloxName, eggs, time_saved FROM players`
		if f2p {
			query += " WHERE robux = 0"
		}

		query += `
			ORDER BY eggs DESC
			LIMIT $1`

		rows, err := s.db.Query(context.Background(), query, LIMIT)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		accounts := make([]*models.Account, 0)
		for rows.Next() {
			account := &models.Account{}
			if err := rows.Scan(&account.ID, &account.Name, &account.Eggs, &account.LastSavedTime); err != nil {
				return nil, err
			}
			accounts = append(accounts, account)
		}
		return accounts, nil
	}

	allF2P, _ := GetRows(true)
	fullResponse.F2P = allF2P

	allNonF2P, _ := GetRows(false)
	fullResponse.NonF2P = allNonF2P

	return fullResponse, nil
}
