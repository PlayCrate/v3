package storage

import (
	"context"
	"fmt"

	"github.com/kattah7/v3/models"
)

func (s *PostgresStore) DeleteSeasonLB() error {
	query := `DELETE FROM season_lb;`
	_, err := s.db.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) InsertSeasonLB(account *models.SeasonLBAccount) error {
	if account.RobloxID == 0 {
		return fmt.Errorf("robloxId cannot be empty")
	}

	if account.SeasonMain == 0 && account.SeasonEvent == 0 {
		return fmt.Errorf("season main and event cannot be empty")
	}

	query := `
	INSERT INTO season_lb (robloxId, season_main, season_event) 
	VALUES ($1, $2, $3) 
	ON CONFLICT (robloxId) 
	DO UPDATE SET season_main = $2, season_event = $3`
	_, err := s.db.Exec(context.Background(), query, account.RobloxID, account.SeasonMain, account.SeasonEvent)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetSeasonLB() (*models.GetSeasonLB, error) {
	mainQuery := `
	SELECT robloxId, season_main
	FROM season_lb
	WHERE season_main != 0
	ORDER BY season_main DESC
	LIMIT 50;
	`

	// Query to retrieve the top 50 records for season_event
	eventQuery := `
	SELECT robloxId, season_event
	FROM season_lb
	WHERE season_event != 0
	ORDER BY season_event DESC
	LIMIT 50;
	`

	// Execute the mainQuery to retrieve top season_main records
	mainRows, mainErr := s.db.Query(context.Background(), mainQuery)
	if mainErr != nil {
		return nil, mainErr
	}
	defer mainRows.Close()

	// Execute the eventQuery to retrieve top season_event records
	eventRows, eventErr := s.db.Query(context.Background(), eventQuery)
	if eventErr != nil {
		return nil, eventErr
	}
	defer eventRows.Close()

	// Create a GetSeasonLB instance to store the results
	seasonLB := models.NewGetSeasonLB()

	// Populate the season_main results
	for mainRows.Next() {
		var robloxID, seasonMainCount int64
		if err := mainRows.Scan(&robloxID, &seasonMainCount); err != nil {
			return nil, err
		}
		seasonLB.SeasonMain = append(seasonLB.SeasonMain, struct {
			RobloxID        int64 `json:"robloxId"`
			SeasonMainCount int64 `json:"value"`
		}{
			RobloxID:        robloxID,
			SeasonMainCount: seasonMainCount,
		})
	}

	// Populate the season_event results
	for eventRows.Next() {
		var robloxID, seasonEventCount int64
		if err := eventRows.Scan(&robloxID, &seasonEventCount); err != nil {
			return nil, err
		}
		seasonLB.SeasonEvent = append(seasonLB.SeasonEvent, struct {
			RobloxID         int64 `json:"robloxId"`
			SeasonEventCount int64 `json:"value"`
		}{
			RobloxID:         robloxID,
			SeasonEventCount: seasonEventCount,
		})
	}

	if err := mainRows.Err(); err != nil {
		return nil, err
	}

	if err := eventRows.Err(); err != nil {
		return nil, err
	}

	return seasonLB, nil
}
