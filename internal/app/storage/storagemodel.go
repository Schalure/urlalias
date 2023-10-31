package storage

// Storage model for long URL and their alias keys
type AliasURLModel struct {
	ID       uint64	`json:"uuid"`
	ShortKey string	`json:"short_key"`
	LongURL  string `json:"original-url"`
}
