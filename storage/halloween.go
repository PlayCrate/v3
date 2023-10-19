package storage

import (
	"context"
	// "encoding/json"
	"fmt"

	"github.com/kattah7/v3/models"
)

func (s *PostgresStore) InsertHalloweenLB(account *models.HalloweenAccount) error {
	if account.RobloxID == 0 {
		return fmt.Errorf("robloxId cannot be empty")
	}

	if account.Houses == 0 && account.Candies == 0 {
		return fmt.Errorf("houses and candy cannot be empty")
	}

	query := `
	INSERT INTO halloween_lb (robloxId, houses, candies) 
	VALUES ($1, $2, $3) 
	ON CONFLICT (robloxId) 
	DO UPDATE SET houses = $2, candies = $3`
	_, err := s.db.Exec(context.Background(), query, account.RobloxID, account.Houses, account.Candies)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetHalloweenLB() (*models.GetHalloweenLB, error) {
	housesQuery := `
	SELECT robloxId, houses
	FROM halloween_lb
	WHERE houses != 0
	ORDER BY houses DESC
	LIMIT 50;
	`

	candiesQuery := `
	SELECT robloxId, candies
	FROM halloween_lb
	WHERE candies != 0
	ORDER BY candies DESC
	LIMIT 50;
	`

	housesRows, housesErr := s.db.Query(context.Background(), housesQuery)
	if housesErr != nil {
		return nil, housesErr
	}
	defer housesRows.Close()

	candiesRows, candiesErr := s.db.Query(context.Background(), candiesQuery)
	if candiesErr != nil {
		return nil, candiesErr
	}
	defer candiesRows.Close()

	halloweenLB := models.NewGetHalloweenLB()
	for housesRows.Next() {
		var robloxID, houseCount int64
		if err := housesRows.Scan(&robloxID, &houseCount); err != nil {
			return nil, err
		}
		halloweenLB.HousesCount = append(halloweenLB.HousesCount, struct {
			RobloxID   int64 `json:"robloxId"`
			HouseCount int64 `json:"value"`
		}{
			RobloxID:   robloxID,
			HouseCount: houseCount,
		})
	}

	for candiesRows.Next() {
		var robloxID, candiesCount int64
		if err := candiesRows.Scan(&robloxID, &candiesCount); err != nil {
			return nil, err
		}
		halloweenLB.CandiesCount = append(halloweenLB.CandiesCount, struct {
			RobloxID   int64 `json:"robloxId"`
			CandyCount int64 `json:"value"`
		}{
			RobloxID:   robloxID,
			CandyCount: candiesCount,
		})
	}

	if err := housesRows.Err(); err != nil {
		return nil, err
	}

	if err := candiesRows.Err(); err != nil {
		return nil, err
	}

	return halloweenLB, nil
}
