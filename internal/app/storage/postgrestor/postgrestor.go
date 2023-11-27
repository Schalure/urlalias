package postgrestor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Schalure/urlalias/internal/app/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dbConnectionString string) (*Storage, error) {

	db, err := sql.Open("pgx", dbConnectionString)
	if err != nil {
		log.Panicln(err)
	}

	if _, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users(
		user_id serial PRIMARY KEY
		);
	`); err != nil {
		return nil, err
	}

	if _, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS aliases(
		id serial PRIMARY KEY,
		user_id integer NOT NULL REFERENCES users(user_id),
		original_url text NOT NULL UNIQUE,
		short_key varchar(9) NOT NULL,
		is_deleted boolean NOT NULL DEFAULT false
		);
	`); err != nil {
		return nil, err
	}

	s := Storage{
		db: db,
	}

	s.GetLastShortKey()

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) CreateUser() (uint64, error) {

	lastID := 0
	err := s.db.QueryRow(`insert into users default values returning user_id`).Scan(&lastID)
	if err != nil {
		return 0, errors.New("can't create new user")
	}
	fmt.Println(lastID)
	return uint64(lastID), nil
}

// ------------------------------------------------------------
//
//	Save pair "shortKey, longURL" to db
//	This is interfase method of "Storager" interface
//	Input:
//		urlAliasNode *repositories.AliasURLModel
//	Output:
//		error - if not nil, can not save "urlAliasNode" because duplicate key
func (s *Storage) Save(urlAliasNode *models.AliasURLModel) error {

	_, err := s.db.Exec(`INSERT INTO aliases(user_id, original_url, short_key) VALUES($1, $2, $3);`, urlAliasNode.UserID, urlAliasNode.LongURL, urlAliasNode.ShortKey)

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
func (s *Storage) SaveAll(urlAliasNodes []models.AliasURLModel) error {

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, node := range urlAliasNodes {

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		_, err := tx.ExecContext(ctx, `insert into aliases(user_id, original_url, short_key) VALUES($1, $2, $3);`, node.UserID, node.LongURL, node.ShortKey)
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
func (s *Storage) FindByShortKey(shortKey string) *models.AliasURLModel {

	var aliasNode = new(models.AliasURLModel)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, `SELECT id, user_id, original_url, short_key, is_deleted FROM aliases WHERE short_key = $1;`, shortKey)
	if err := row.Scan(&aliasNode.ID, &aliasNode.UserID, &aliasNode.LongURL, &aliasNode.ShortKey, &aliasNode.DeletedFlag); err != nil {
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
func (s *Storage) FindByLongURL(longURL string) *models.AliasURLModel {

	var aliasNode = new(models.AliasURLModel)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, `SELECT id, user_id, original_url, short_key, is_deleted FROM aliases WHERE original_url=$1;`, longURL)
	if err := row.Scan(&aliasNode.ID, &aliasNode.UserID, &aliasNode.LongURL, &aliasNode.ShortKey, &aliasNode.DeletedFlag); err != nil {
		return nil
	}
	return aliasNode
}

func (s *Storage) FindByUserID(ctx context.Context, userID uint64) ([]models.AliasURLModel, error) {

	rows, err := s.db.QueryContext(ctx, `select original_url, short_key from aliases where user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []models.AliasURLModel
	var node models.AliasURLModel

	for rows.Next() {
		err = rows.Scan(&node.LongURL, &node.ShortKey)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nodes, nil
}


// ------------------------------------------------------------
//
//	Mark aliases like "deleted" by aliasesID
func (s *Storage) MarkDeleted(ctx context.Context, aliasesID []uint64) error {

	var parametrs []string
	var argsID []interface{}
	for i, ID := range aliasesID {
		parametrs = append(parametrs, fmt.Sprintf("id=$%d", i + 1))
		argsID = append(argsID, ID)
	}

	_, err := s.db.Exec(`update aliases set is_deleted = true where (` + strings.Join(parametrs, " OR ") + `);`, argsID...)
	return err
}


// ------------------------------------------------------------
//
//	Get the last saved key
func (s *Storage) GetLastShortKey() string {

	var shortKey string
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, `select short_key from aliases where id=(select max(id) from aliases);`)
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
func (s *Storage) IsConnected() bool {

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
func (s *Storage) Close() error {

	return s.db.Close()
}
