package postgrestor

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/Schalure/urlalias/internal/app/aliasmaker"
	"github.com/Schalure/urlalias/internal/app/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgreStor struct {
	db *sql.DB
}

func NewPostgreStor(dbConnectionString string) aliasmaker.Storager {

	db, err := sql.Open("pgx", dbConnectionString)
	if err != nil {
		log.Panicln(err)
	}

	return &PostgreStor{
		db: db,
	}
}

// ------------------------------------------------------------
//
//	Save pair "shortKey, longURL" to db
//	This is interfase method of "Storager" interface
//	Input:
//		urlAliasNode *repositories.AliasURLModel
//	Output:
//		error - if not nil, can not save "urlAliasNode" because duplicate key
func (s *PostgreStor) Save(urlAliasNode *storage.AliasURLModel) error {

	panic("save to PostgreSQL no implemented")
}

// ------------------------------------------------------------
//
//	Find "urlAliasNode models.AliasURLModel" by short key
//	This is interfase method of "Storager" interface
//	Input:
//		shortKey string
//	Output:
//		*repositories.AliasURLModel
//		error - if can not find "urlAliasNode" by short key
func (s *PostgreStor) FindByShortKey(shortKey string) (*storage.AliasURLModel, error) {

	panic("find by short key from PostgreSQL no implemented")
}

// ------------------------------------------------------------
//
//	Find "urlAliasNode models.AliasURLModel" by long URL
//	This is interfase method of "Storager" interface
//	Input:
//		longURL string
//	Output:
//		*repositories.AliasURLModel
//		error - if can not find "urlAliasNode" by long URL
func (s *PostgreStor) FindByLongURL(longURL string) (*storage.AliasURLModel, error) {

	panic("find by long URL from PostgreSQL no implemented")
}

// ------------------------------------------------------------
//
//	Check connection to DB
//	This is interfase method of "Storager" interface
//	Output:
//		bool - true: connection is
//			   false: connection isn't
//		error - if can not find "urlAliasNode" by long URL
func (s *PostgreStor) IsConnected() bool {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return false
	}
	return true
}

// ------------------------------------------------------------
//
//	Close connection to DB
//	This is interfase method of "Storager" interface
//	Output:
//		error
func (s *PostgreStor) Close() error {

	return s.db.Close()
}
