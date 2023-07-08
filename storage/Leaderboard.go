package storage

import (
	"context"

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
