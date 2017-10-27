package w

import (
	"fmt"
	"strings"

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

// History of a list of page. Newest first.
func (p *Postgres) History(pages ...string) ([]Entry, error) {
	var (
		in   []string
		args []interface{}
	)
	for i, p := range pages {
		in = append(in, fmt.Sprintf("$%d", i+1))
		args = append(args, p)
	}
	return p.queryUpdates(`
		WHERE page IN (`+strings.Join(in, ",")+`)
		ORDER BY timestamp DESC
    `, args...)
}

func (p *Postgres) queryUpdates(where string, args ...interface{}) ([]Entry, error) {
	var es []Entry
	rows, err := p.conn.Query(`
		SELECT page, timestamp, stable_version
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
		(page, timestamp, stable_version)
	VALUES
		($1, $2, $3)
`, e.Page, e.T, e.StableVersion)
	return err
}
