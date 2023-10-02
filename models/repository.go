package models

type RepositoryURL interface{
	Save(urlAlias AliasURLModel) (*AliasURLModel, error)
	FindByShortKey(shortKey string) (*AliasURLModel, error)
	FindByLongURL(longURL string) (*AliasURLModel, error)
}