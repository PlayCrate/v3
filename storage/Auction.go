package storage

import (
	"context"
	"fmt"

	"github.com/kattah7/v3/models"
)

func (s *PostgresStore) ListAuction(item *models.AuctionAccount) error {
	if item.ItemType == "" {
		return fmt.Errorf("itemType cannot be empty")
	}

	if item.ItemType != "EGG" && item.ItemType != "PET" && item.ItemType != "BOOST" && item.ItemType != "POTION" {
		return fmt.Errorf("itemType must be EGG, PET, BOOST, or POTION")
	}

	if item.ItemData == nil {
		return fmt.Errorf("itemData cannot be empty")
	}

	if item.Price == 0 {
		return fmt.Errorf("price cannot be empty")
	}

	checkQuery := `SELECT robloxId FROM auctions WHERE robloxId = $1`
	rows, err := s.db.Query(context.Background(), checkQuery, item.ID)
	if err != nil {
		return fmt.Errorf("Unable to query row: %w", err)
	}

	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		if count >= 5 {
			return fmt.Errorf("Exceeded maximum allowed rows")
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("Error iterating rows: %w", err)
	}

	query := `
	INSERT INTO auctions (robloxId, robloxName, itemType, itemData, startPrice)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err2 := s.db.Exec(context.Background(), query, item.ID, item.Name, item.ItemType, item.ItemData, item.Price)
	if err2 != nil {
		return fmt.Errorf("unable to insert row: %w", err2)
	}

	return nil
}

func (s *PostgresStore) RemoveAuction(item *models.AuctionAccount) error {
	query := `DELETE FROM auctions WHERE id = $1`

	result, err := s.db.Exec(context.Background(), query, item.UID)
	if err != nil {
		return fmt.Errorf("Unable to delete row: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("No rows affected")
	}

	return nil
}

func (s *PostgresStore) GetAuctions() ([]*models.AuctionAccount, error) {
	query := `SELECT id, robloxId, robloxName, itemType, itemData, startPrice FROM auctions ORDER BY id DESC`

	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("Unable to query row: %w", err)
	}

	defer rows.Close()

	var auctions []*models.AuctionAccount

	for rows.Next() {
		item := &models.AuctionAccount{}
		err := rows.Scan(&item.UID, &item.ID, &item.Name, &item.ItemType, &item.ItemData, &item.Price)
		if err != nil {
			return nil, fmt.Errorf("Unable to scan row: %w", err)
		}

		auctions = append(auctions, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating rows: %w", err)
	}

	return auctions, nil
}
