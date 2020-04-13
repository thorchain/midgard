package models

// AssetDetails contains asset pool creation date and it's price in Rune.
type AssetDetails struct {
	DateCreated int64   `json:"dateCreated"`
	PriceInRune float64 `json:"priceInRune,string"`
}
