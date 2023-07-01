package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/kattah7/v3/models"
)

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
