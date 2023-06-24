package models

import "encoding/json"

type AuctionAccount struct {
	UID      int64           `json:"id"`
	ID       int64           `json:"robloxId"`
	Name     string          `json:"robloxName"`
	ItemType string          `json:"itemType"`
	ItemData json.RawMessage `json:"itemData"`
	Price    int64           `json:"startPrice"`
}

func NewItem(ID int64, Name string, ItemType string, ItemData json.RawMessage, Price int64) *AuctionAccount {
	return &AuctionAccount{
		ID:       ID,
		Name:     Name,
		ItemType: ItemType,
		ItemData: ItemData,
		Price:    Price,
	}
}
