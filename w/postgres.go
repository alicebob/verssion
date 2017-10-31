package w

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
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

func (p *Postgres) Current(page string) (*Entry, error) {
	rows, err := p.queryUpdates(`
		WHERE page=$1
		ORDER BY timestamp DESC
		LIMIT 1
    `, page)
	if err != nil || len(rows) == 0 {
		return nil, err
	}
	e := rows[0]
	return &e, nil
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
		e.T = e.T.UTC()
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

func (p *Postgres) Known() ([]string, error) {
	var ps []string
	rows, err := p.conn.Query(`
		SELECT DISTINCT(page)
		FROM updates
		ORDER BY page`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, rows.Err()
}

func (p *Postgres) CreateCurated() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	cid := id.String()
	_, err = p.conn.Exec(`
		INSERT INTO curated (id, created, lastused, lastupdated)
		VALUES ($1, now(), now(), now())`,
		cid,
	)
	return cid, err
}

func (p *Postgres) StoreCurated(cur Curated) error {
	return nil
}

func (p *Postgres) LoadCurated(id string) (*Curated, error) {
	row := p.conn.QueryRow(`
		SELECT id, created, lastused, lastupdated, title
		FROM curated
		WHERE id=$1`,
		id,
	)
	cur := Curated{}
	if err := row.Scan(&cur.ID, &cur.Created, &cur.LastUsed, &cur.LastUpdated, &cur.CustomTitle); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	pg, err := p.curatedPages(id)
	if err != nil {
		return nil, err
	}
	cur.Pages = pg
	return &cur, nil
}

func (p *Postgres) curatedPages(id string) ([]string, error) {
	var ps []string
	rows, err := p.conn.Query(`
		SELECT page
		FROM curated_pages
		WHERE curated_id=$1
		ORDER BY page`,
		id,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, rows.Err()
}

// pages must be unique
func (p *Postgres) CuratedPages(id string, pages []string) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM curated_pages WHERE curated_id=$1`, id); err != nil {
		return err
	}
	for _, p := range pages {
		if _, err := tx.Exec(`INSERT INTO curated_pages (curated_id, page) VALUES ($1, $2)`, id, p); err != nil {
			return err
		}
	}
	if _, err := p.conn.Exec(`UPDATE curated SET lastupdated=now() WHERE id=$1`, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (p *Postgres) CuratedUsed(id string) error {
	_, err := p.conn.Exec(`UPDATE curated SET lastused=now(), used=used+1 WHERE id=$1`, id)
	return err
}

func (p *Postgres) CuratedTitle(id, title string) error {
	_, err := p.conn.Exec(`UPDATE curated SET title=$2, lastupdated=now() WHERE id=$1`, id, title)
	return err
}
