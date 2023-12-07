package aliasentity

// Storage model for long URL and their alias keys
type AliasURLModel struct {
	ID          uint64 `json:"uuid" db:"uuid"`
	UserID      uint64 `json:"user_id" db:"user_id"`
	ShortKey    string `json:"short_url" db:"short_url"`
	LongURL     string `json:"original_url" db:"original_url"`
	DeletedFlag bool   `json:"is_deleted" db:"is_deleted"`
}
