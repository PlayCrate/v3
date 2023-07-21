package storage

import (
	"context"
	"fmt"

	"github.com/kattah7/v3/models"
)

func (s *PostgresStore) GetSpecificPlayer(robloxId int64) (*models.AccountLookup, error) {
	query := `
	SELECT
		robloxId, robloxName, secrets, eggs, bubbles, power, playtime, robux,
		(SELECT COUNT(secrets) + 1 FROM players WHERE secrets > p.secrets) AS secretsRank,
		(SELECT COUNT(eggs) + 1 FROM players WHERE eggs > p.eggs) AS eggsRank,
		(SELECT COUNT(bubbles) + 1 FROM players WHERE bubbles > p.bubbles) AS bubblesRank,
		(SELECT COUNT(power) + 1 FROM players WHERE power > p.power) AS powerRank,
		(SELECT COUNT(playtime) + 1 FROM players WHERE playtime > p.playtime) AS playtimeRank,
		(SELECT COUNT(robux) + 1 FROM players WHERE robux > p.robux) AS robuxRank,
		CASE
			WHEN robux = 0 THEN (SELECT COUNT(secrets) + 1 FROM players WHERE secrets > p.secrets AND robux = 0)
			ELSE NULL
		END AS freeToPlaySecretsRank,
		CASE
			WHEN robux = 0 THEN (SELECT COUNT(eggs) + 1 FROM players WHERE eggs > p.eggs AND robux = 0)
			ELSE NULL
		END AS freeToPlayEggsRank,
		CASE
			WHEN robux = 0 THEN (SELECT COUNT(bubbles) + 1 FROM players WHERE bubbles > p.bubbles AND robux = 0)
			ELSE NULL
		END AS freeToPlayBubblesRank,
		CASE
			WHEN robux = 0 THEN (SELECT COUNT(power) + 1 FROM players WHERE power > p.power AND robux = 0)
			ELSE NULL
		END AS freeToPlayPowerRank
	FROM
		players AS p
	WHERE
		robloxId = $1
`
	row := s.db.QueryRow(context.Background(), query, robloxId)

	account := &models.AccountLookup{}
	err := row.Scan(
		&account.RobloxID, &account.RobloxName,
		&account.Secrets, &account.Eggs, &account.Bubbles, &account.Power, &account.Playtime, &account.Robux,
		&account.SecretsRank, &account.EggsRank, &account.BubblesRank, &account.PowerRank, &account.PlaytimeRank, &account.RobuxRank,
		&account.F2PSecretsRank, &account.F2PEggsRank, &account.F2PBubblesRank, &account.F2PPowerRank,
	)

	if err != nil {
		return nil, err
	}

	return account, nil
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

	return nil
}

func (s *PostgresStore) GetBubbles() (*models.PlayerDataResponse, error) {
	fullResponse := &models.PlayerDataResponse{
		F2P:    make([]*models.Account, 0),
		NonF2P: make([]*models.Account, 0),
	}

	GetRows := func(f2p bool) ([]*models.Account, error) {
		query := `SELECT robloxId, robloxName, bubbles, time_saved FROM players`
		if f2p {
			query += " WHERE robux = 0"
		}

		query += `
			ORDER BY bubbles DESC
			LIMIT $1`

		rows, err := s.db.Query(context.Background(), query, LIMIT)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		accounts := make([]*models.Account, 0)
		for rows.Next() {
			account := &models.Account{}
			if err := rows.Scan(&account.ID, &account.Name, &account.Bubbles, &account.LastSavedTime); err != nil {
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

func (s *PostgresStore) GetEggs() (*models.PlayerDataResponse, error) {
	fullResponse := &models.PlayerDataResponse{
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

func (s *PostgresStore) GetPlaytime() (*models.PlayerDataResponse, error) {
	fullResponse := &models.PlayerDataResponse{
		F2P:    make([]*models.Account, 0),
		NonF2P: make([]*models.Account, 0),
	}

	GetRows := func(f2p bool) ([]*models.Account, error) {
		query := `SELECT robloxId, robloxName, playtime, time_saved FROM players`

		if f2p {
			query += " WHERE robux = 0"
		}

		query += `
			ORDER BY playtime DESC
			LIMIT $1`

		rows, err := s.db.Query(context.Background(), query, LIMIT)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		accounts := make([]*models.Account, 0)
		for rows.Next() {
			account := &models.Account{}
			if err := rows.Scan(&account.ID, &account.Name, &account.Playtime, &account.LastSavedTime); err != nil {
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

func (s *PostgresStore) GetPower() (*models.PlayerDataResponse, error) {
	fullResponse := &models.PlayerDataResponse{
		F2P:    make([]*models.Account, 0),
		NonF2P: make([]*models.Account, 0),
	}

	GetRows := func(f2p bool) ([]*models.Account, error) {
		query := `SELECT robloxId, robloxName, power, time_saved FROM players`

		if f2p {
			query += " WHERE robux = 0"
		}

		query += `
			ORDER BY power DESC
			LIMIT $1`

		rows, err := s.db.Query(context.Background(), query, LIMIT)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		accounts := make([]*models.Account, 0)
		for rows.Next() {
			account := &models.Account{}
			if err := rows.Scan(&account.ID, &account.Name, &account.Power, &account.LastSavedTime); err != nil {
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

func (s *PostgresStore) GetSecrets() (*models.PlayerDataResponse, error) {
	fullResponse := &models.PlayerDataResponse{
		F2P:    make([]*models.Account, 0),
		NonF2P: make([]*models.Account, 0),
	}

	GetRows := func(f2p bool) ([]*models.Account, error) {
		query := `SELECT robloxId, robloxName, secrets, time_saved FROM players`

		if f2p {
			query += " WHERE robux = 0"
		}

		query += `
			ORDER BY secrets DESC
			LIMIT $1`

		rows, err := s.db.Query(context.Background(), query, LIMIT)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		accounts := make([]*models.Account, 0)
		for rows.Next() {
			account := &models.Account{}
			if err := rows.Scan(&account.ID, &account.Name, &account.Secrets, &account.LastSavedTime); err != nil {
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
