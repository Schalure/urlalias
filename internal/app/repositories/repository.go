package repositories

// Storage model for long URL and their alias keys
type AliasURLModel struct {
	ID       uint64
	ShortKey string
	LongURL  string
}

// Access interface to storage
type RepositoryURL interface {
	Save(s *AliasURLModel) error
	FindByShortKey(shortKey string) (*AliasURLModel, error)
	FindByLongURL(longURL string) (*AliasURLModel, error)
}
