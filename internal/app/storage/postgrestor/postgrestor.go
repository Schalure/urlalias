package postgrestor

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/Schalure/urlalias/internal/app/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgreStor struct {
	db *sql.DB
}

func NewPostgreStor(dbConnectionString string) (*PostgreStor, error) {

	db, err := sql.Open("pgx", dbConnectionString)
	if err != nil {
		log.Panicln(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS aliases(
		id serial PRIMARY KEY,
		originalURL text NOT NULL,
		shortKey text NOT NULL
		);`)

	if err != nil {
		return nil, err
	}

	s := PostgreStor{
		db: db,
	}

	s.GetLastShortKey()

	return &PostgreStor{
		db: db,
	}, nil
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

	_, err := s.db.Exec(`INSERT INTO aliases(originalURL, shortKey) VALUES($1, $2);`, urlAliasNode.LongURL, urlAliasNode.ShortKey)

	if err != nil {
		return err
	}
	return nil
}

// ------------------------------------------------------------
//
//	Save array of pairs "shortKey, longURL" to db
//	This is interfase method of "Storager" interface
//	Input:
//		urlAliasNode []repositories.AliasURLModel
//	Output:
//		error - if not nil, can not save "[]storage.AliasURLModel"
func (s *PostgreStor) SaveAll(urlAliasNodes []storage.AliasURLModel) error {

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, node := range urlAliasNodes {

		ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
		defer cancel()
		_, err := tx.ExecContext(ctx, `insert into aliases(originalURL, shortKey) VALUES($1, $2);`, node.LongURL, node.ShortKey)
			// sql.Named("long_url", node.LongURL),
			// sql.Named("short_url", node.ShortKey))
		if err != nil {
			return err
		}
	}
	return tx.Commit()
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
func (s *PostgreStor) FindByShortKey(shortKey string) *storage.AliasURLModel {

	var aliasNode = new(storage.AliasURLModel)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, `SELECT id, originalURL, shortKey FROM aliases WHERE shortKey = $1;`, shortKey)
	if err := row.Scan(&aliasNode.ID, &aliasNode.LongURL, &aliasNode.ShortKey); err != nil {
		return nil
	}
	return aliasNode
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
func (s *PostgreStor) FindByLongURL(longURL string) *storage.AliasURLModel {

	var aliasNode = new(storage.AliasURLModel)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, `SELECT id, originalURL, shortKey FROM aliases WHERE originalURL=$1;`, longURL)
	if err := row.Scan(&aliasNode.ID, &aliasNode.LongURL, &aliasNode.ShortKey); err != nil {
		return nil
	}
	return aliasNode
}

// ------------------------------------------------------------
//
//	Get the last saved key
//	This is interfase method of "Storager" interface
//	Output:
//		string - last saved key
func (s *PostgreStor) GetLastShortKey() string {

	var shortKey string
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, `select shortkey from aliases where id=(select max(id) from aliases);`)
	if err := row.Scan(&shortKey); err != nil {
		return ""
	}
	return shortKey
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
