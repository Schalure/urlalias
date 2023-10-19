package storage

// Storage model for long URL and their alias keys
type AliasURLModel struct {
	ID       uint64
	ShortKey string
	LongURL  string
}
