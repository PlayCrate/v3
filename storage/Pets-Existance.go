package storage

import (
	"context"
	"fmt"

	"github.com/kattah7/v3/models"
)

func (s *PostgresStore) InsertPetsExistance(pet *models.PetsExistance) error {
	if pet.RobloxID == 0 {
		return fmt.Errorf("robloxId cannot be empty")
	}

	if pet.Pets == nil {
		return fmt.Errorf("pets cannot be empty")
	}

	deleteQuery := `DELETE FROM pets_exist WHERE robloxId = $1`
	_, err := s.db.Exec(context.Background(), deleteQuery, pet.RobloxID)
	if err != nil {
		return err
	}

	for _, v := range pet.Pets {
		for petId, petCount := range v {
			if petId == "" {
				return fmt.Errorf("pet name cannot be empty")
			}

			if petCount == 0 {
				return fmt.Errorf("pet count cannot be empty")
			}

			query := `INSERT INTO pets_exist (robloxId, petId, petCount) VALUES ($1, $2, $3) ON CONFLICT (robloxId, petId) DO UPDATE SET petCount = $3`

			_, err := s.db.Exec(context.Background(), query, pet.RobloxID, petId, petCount)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *PostgresStore) GetPetsExistance() ([]*models.GetPetsExistance, error) {
	query := `SELECT petId, petCount FROM pets_exist ORDER BY petCount DESC`

	rows, err := s.db.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pets []*models.GetPetsExistance

	for rows.Next() {
		var petId string
		var petCount int

		err := rows.Scan(&petId, &petCount)
		if err != nil {
			return nil, err
		}

		var existingPet *models.GetPetsExistance
		for _, pet := range pets {
			if pet.PetID == petId {
				existingPet = pet
				break
			}
		}

		if existingPet != nil {
			existingPet.PetCount += petCount
		} else {
			pet := &models.GetPetsExistance{
				PetID:    petId,
				PetCount: petCount,
			}
			pets = append(pets, pet)
		}
	}

	return pets, nil
}

func (s *PostgresStore) DeletePetsExistence(pet *models.PetsExistance) error {
	if pet.RobloxID == 0 {
		return fmt.Errorf("robloxId cannot be empty")
	}

	deleteQuery := `DELETE FROM pets_exist WHERE robloxId = $1`
	del, err := s.db.Exec(context.Background(), deleteQuery, pet.RobloxID)
	if err != nil {
		return err
	}

	if del.RowsAffected() == 0 {
		return fmt.Errorf("robloxId %d does not exist", pet.RobloxID)
	}

	return nil
}
