package w

import (
	"github.com/jackc/pgx"
)

const DBURL = "postgresql:///w"

type Postgres struct {
	conn *pgx.Conn
}

var _ DB = &Postgres{}

func NewPostgres(url string) (*Postgres, error) {
	if url == "" {
		url = DBURL
	}
	cc, err := pgx.ParseURI(url)
	if err != nil {
		return nil, err
	}
	conn, err := pgx.Connect(cc)
	if err != nil {
		return nil, err
	}

	p := &Postgres{
		conn: conn,
	}
	return p, nil
}

// Bunch of recent changes. Just to have something
func (p *Postgres) Recent() ([]Entry, error) {
	return p.queryUpdates(`
		ORDER BY timestamp DESC
    `)
}

// History of a single page. Newest first.
func (p *Postgres) History(page string) ([]Entry, error) {
	return p.queryUpdates(`
		WHERE title=$1
		ORDER BY timestamp DESC
    `, page)
}

func (p *Postgres) queryUpdates(where string, args ...interface{}) ([]Entry, error) {
	var es []Entry
	rows, err := p.conn.Query(`
		SELECT title, timestamp, stable_version
		FROM updates`+where, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.Page, &e.T, &e.StableVersion); err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	return es, rows.Err()
}

func (p *Postgres) Store(e Entry) error {
	_, err := p.conn.Exec(`
	INSERT INTO page
		(title, timestamp, stable_version)
	VALUES
		($1, $2, $3)
`, e.Page, e.T, e.StableVersion)
	return err
}
