package models

// Storage model for long URL and their alias keys
type AliasURLModel struct {
	ID       uint64 `json:"uuid"`
	UserID   uint64 `json:"user_id"`
	ShortKey string `json:"short_url"`
	LongURL  string `json:"original_url"`
}
