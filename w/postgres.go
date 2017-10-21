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

// Load the latest we have. Could be nil.
func (p *Postgres) Load(page string) (*Entry, error) {
	row := p.conn.QueryRow(`
    SELECT title, revision, timestamp, stable_version
    FROM page
    WHERE title=$1
	ORDER BY revision DESC
	LIMIT 1
    `, page)
	var e Entry
	if err := row.Scan(&e.Page, &e.Revision, &e.T, &e.StableVersion); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (p *Postgres) Store(e Entry) error {
	_, err := p.conn.Exec(`
	INSERT INTO page
		(title, revision, timestamp, stable_version)
	VALUES
		($1, $2, $3, $4)
`, e.Page, e.Revision, e.T, e.StableVersion)
	return err
}
