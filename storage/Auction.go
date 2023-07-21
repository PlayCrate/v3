package storage

import (
	"context"
	"fmt"

	"github.com/kattah7/v3/models"
)

func (s *PostgresStore) ListAuction(item *models.AuctionAccount) error {
	if item.ID == 0 && item.Name == "" {
		return fmt.Errorf("robloxId or robloxName cannot be empty")
	}

	if item.ItemType == "" {
		return fmt.Errorf("itemType cannot be empty")
	}

	if item.ItemType != "EGG" && item.ItemType != "PET" && item.ItemType != "BOOST" && item.ItemType != "POTION" {
		return fmt.Errorf("itemType must be EGG, PET, BOOST, or POTION")
	}

	if item.ItemData == nil {
		return fmt.Errorf("itemData cannot be empty")
	}

	fmt.Println(item.PriceType)
	if item.Price == 0 {
		return fmt.Errorf("price cannot be empty")
	}

	if item.PriceType != "Diamonds" && item.PriceType != "Coins" && item.PriceType != "DarkCoins" && item.PriceType != "Pearls" && item.PriceType != "Candy" && item.PriceType != "Chocolate" {
		return fmt.Errorf("priceType must be Diamonds, Coins, DarkCoins, Pearls, Candy, or Chocolate")
	}

	checkQuery := `SELECT robloxId FROM auctions WHERE robloxId = $1 AND status = 'OPEN'`
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
	INSERT INTO auctions (robloxId, robloxName, itemType, itemData, startPrice, priceType)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err2 := s.db.Exec(context.Background(), query, item.ID, item.Name, item.ItemType, item.ItemData, item.Price, item.PriceType)
	if err2 != nil {
		return fmt.Errorf("unable to insert row: %w", err2)
	}

	return nil
}

func (s *PostgresStore) PurchaseAuction(item *models.AuctionAccount) error {
	if item.UID == 0 {
		return fmt.Errorf("id cannot be empty")
	}

	query := `UPDATE auctions SET status = 'PURCHASED' WHERE id = $1 AND status = 'OPEN'`

	result, err := s.db.Exec(context.Background(), query, item.UID)
	if err != nil {
		return fmt.Errorf("Unable to update row: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("No rows affected")
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
	query := `SELECT id, robloxId, robloxName, itemType, itemData, startPrice, priceType, listed FROM auctions WHERE status = 'OPEN' ORDER BY id DESC`

	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("Unable to query row: %w", err)
	}

	defer rows.Close()

	var auctions []*models.AuctionAccount

	for rows.Next() {
		item := &models.AuctionAccount{}
		err := rows.Scan(&item.UID, &item.ID, &item.Name, &item.ItemType, &item.ItemData, &item.Price, &item.PriceType, &item.ListedDate)
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

func (s *PostgresStore) GetAuctionClaims(item *models.AuctionAccount) ([]*models.AuctionAccount, error) {
	if item.ID == 0 {
		return nil, fmt.Errorf("robloxId cannot be empty")
	}

	query := `SELECT id, robloxId, robloxName, itemType, itemData, startPrice, priceType FROM auctions WHERE status = 'PURCHASED' AND robloxId = $1 ORDER BY id DESC`

	rows, err := s.db.Query(context.Background(), query, item.ID)
	if err != nil {
		return nil, fmt.Errorf("Unable to query row: %w", err)
	}

	defer rows.Close()

	var auctions []*models.AuctionAccount

	for rows.Next() {
		item := &models.AuctionAccount{}
		err := rows.Scan(&item.UID, &item.ID, &item.Name, &item.ItemType, &item.ItemData, &item.Price, &item.PriceType)
		if err != nil {
			return nil, fmt.Errorf("Unable to scan row: %w", err)
		}

		auctions = append(auctions, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating rows: %w", err)
	}

	if len(auctions) == 0 {
		emptyArray := make([]*models.AuctionAccount, 0)

		return emptyArray, nil
	}

	return auctions, nil
}

func (s *PostgresStore) AuctionClaim(item *models.AuctionAccount) error {
	if item.UID == 0 {
		return fmt.Errorf("uid cannot be empty")
	}

	if item.ID == 0 {
		return fmt.Errorf("id cannot be empty")
	}

	checkQuery := `SELECT robloxId FROM auctions WHERE id = $1 AND status = 'PURCHASED'`

	var robloxId int64
	err := s.db.QueryRow(context.Background(), checkQuery, item.UID).Scan(&robloxId)

	if err != nil {
		return fmt.Errorf("Unable to query row: %w", err)
	}

	if robloxId != item.ID {
		return fmt.Errorf("robloxId does not match with id")
	}

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

func (s *PostgresStore) GetAuctionListing(item *models.AuctionAccount) ([]*models.AuctionAccount, error) {
	if item.ID == 0 {
		return nil, fmt.Errorf("robloxId cannot be empty")
	}

	query := `SELECT id, robloxId, robloxName, itemType, itemData, startPrice, priceType, listed FROM auctions WHERE status = 'OPEN' AND robloxId = $1 ORDER BY id DESC`

	rows, err := s.db.Query(context.Background(), query, item.ID)
	if err != nil {
		return nil, fmt.Errorf("Unable to query row: %w", err)
	}

	defer rows.Close()

	var auctions []*models.AuctionAccount

	for rows.Next() {
		item := &models.AuctionAccount{}
		err := rows.Scan(&item.UID, &item.ID, &item.Name, &item.ItemType, &item.ItemData, &item.Price, &item.PriceType, &item.ListedDate)
		if err != nil {
			return nil, fmt.Errorf("Unable to scan row: %w", err)
		}

		auctions = append(auctions, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating rows: %w", err)
	}

	if len(auctions) == 0 {
		emptyArray := make([]*models.AuctionAccount, 0)

		return emptyArray, nil
	}

	return auctions, nil
}

func (s *PostgresStore) AuctionUnlist(item *models.AuctionAccount) error {
	if item.UID == 0 {
		return fmt.Errorf("uid cannot be empty")
	}

	if item.ID == 0 {
		return fmt.Errorf("id cannot be empty")
	}

	checkQuery := `SELECT robloxId FROM auctions WHERE id = $1 AND status = 'OPEN'`

	var robloxId int64
	err := s.db.QueryRow(context.Background(), checkQuery, item.UID).Scan(&robloxId)

	if err != nil {
		return fmt.Errorf("Unable to query row: %w", err)
	}

	if robloxId != item.ID {
		return fmt.Errorf("robloxId does not match with id")
	}

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

// func (s *PostgresStore) AuctionExpireList(item *models.AuctionAccount) ([]*models.AuctionAccount, error) {
// 	if item.ID == 0 {
// 		return nil, fmt.Errorf("robloxId cannot be empty")
// 	}

// 	query := `SELECT * FROM auction_expired WHERE robloxId = $1 ORDER BY id DESC`

// 	rows, err := s.db.Query(context.Background(), query, item.ID)

// 	if err != nil {
// 		return nil, fmt.Errorf("Unable to query row: %w", err)
// 	}

// 	defer rows.Close()

// 	var auctions []*models.AuctionAccount

// 	for rows.Next() {
// 		item := &models.AuctionAccount{}
// 		err := rows.Scan(&item.UID, &item.ID, &item.Name, &item.ItemType, &item.ItemData)
// 		if err != nil {
// 			return nil, fmt.Errorf("Unable to scan row: %w", err)
// 		}

// 		auctions = append(auctions, item)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("Error iterating rows: %w", err)
// 	}

// 	if len(auctions) == 0 {
// 		emptyArray := make([]*models.AuctionAccount, 0)

// 		return emptyArray, nil
// 	}

// 	return auctions, nil
// }

// func (s *PostgresStore) AuctionExpireClaim(item *models.AuctionAccount) error {
// 	if item.UID == 0 {
// 		return fmt.Errorf("uid cannot be empty")
// 	}

// 	if item.ID == 0 {
// 		return fmt.Errorf("id cannot be empty")
// 	}

// 	checkQuery := `SELECT robloxId FROM auction_expired WHERE id = $1`

// 	var robloxId int64
// 	err := s.db.QueryRow(context.Background(), checkQuery, item.UID).Scan(&robloxId)

// 	if err != nil {
// 		return fmt.Errorf("Unable to query row: %w", err)
// 	}

// 	if robloxId != item.ID {
// 		return fmt.Errorf("robloxId does not match with id")
// 	}

// 	query := `DELETE FROM auction_expired WHERE id = $1`
// 	result, err := s.db.Exec(context.Background(), query, item.UID)

// 	if err != nil {
// 		return fmt.Errorf("Unable to delete row: %w", err)
// 	}

// 	if result.RowsAffected() == 0 {
// 		return fmt.Errorf("No rows affected")
// 	}

// 	return nil
// }
