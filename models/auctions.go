package models

import (
	"encoding/json"
	"time"
)

type AuctionAccount struct {
	Payload    string          `json:"payload,omitempty"`
	UID        int64           `json:"id"`
	ID         int64           `json:"robloxId"`
	Name       string          `json:"robloxName"`
	ItemType   string          `json:"itemType"`
	ItemData   json.RawMessage `json:"itemData"`
	PriceType  string          `json:"priceType"`
	Price      int64           `json:"startPrice"`
	ListedDate time.Time       `json:"listedDate"`
}

func NewItem(ID int64, Name string, ItemType string, ItemData json.RawMessage, PriceType string, Price int64) *AuctionAccount {
	return &AuctionAccount{
		ID:        ID,
		Name:      Name,
		ItemType:  ItemType,
		ItemData:  ItemData,
		PriceType: PriceType,
		Price:     Price,
	}
}
