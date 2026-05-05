package dto

type Calculate struct {
	ItemsOrdered int `json:"items_ordered"`
}

type CalculateResult struct {
	ItemsOrdered int    `json:"items_ordered"`
	TotalItems   int    `json:"total_items"`
	Packs        []Pack `json:"packs"`
}
