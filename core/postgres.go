package core

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const DBURL = "postgresql:///w"

type Postgres struct {
	conn *pgxpool.Pool
}

var _ DB = &Postgres{}

func NewPostgres(url string) (*Postgres, error) {
	if url == "" {
		url = DBURL
	}
	conn, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	p := &Postgres{
		conn: conn,
	}
	return p, nil
}

func (p *Postgres) Last(page string) (*Page, error) {
	row := p.conn.QueryRow(context.Background(), `
		SELECT page, timestamp, stable_version, homepage
		FROM page
		WHERE page=$1
		ORDER BY timestamp DESC
		LIMIT 1
	`,
		page,
	)
	res := Page{}
	if err := row.Scan(&res.Page, &res.T, &res.StableVersion, &res.Homepage); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	res.T = res.T.UTC()
	return &res, nil
}

func (p *Postgres) Current(limit int, order SortBy) ([]Page, error) {
	q := " ORDER BY " + order.OrderBy()
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT %d", limit)
	}
	return p.queryCurrent(q)
}

func (p *Postgres) CurrentIn(pages ...string) ([]Page, error) {
	if len(pages) == 0 {
		return nil, nil
	}
	return p.queryCurrent(`
		WHERE page = ANY($1)
		ORDER BY timestamp DESC
    `, pages)
}

// History of a list of pages. Newest first.
func (p *Postgres) History(pages ...string) ([]Page, error) {
	if len(pages) == 0 {
		return nil, nil
	}
	return p.queryUpdates(`
		WHERE page = ANY($1)
		ORDER BY timestamp DESC
    `, pages)
}

func (p *Postgres) queryCurrent(where string, args ...interface{}) ([]Page, error) {
	return p.queryPages("current", where, args...)
}

func (p *Postgres) queryUpdates(where string, args ...interface{}) ([]Page, error) {
	return p.queryPages("updates", where, args...)
}

func (p *Postgres) queryPages(table, where string, args ...interface{}) ([]Page, error) {
	var es []Page
	rows, err := p.conn.Query(context.Background(), `
		SELECT page, timestamp, stable_version, homepage
		FROM `+table+where, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var e Page
		if err := rows.Scan(&e.Page, &e.T, &e.StableVersion, &e.Homepage); err != nil {
			return nil, err
		}
		e.T = e.T.UTC()
		es = append(es, e)
	}
	return es, rows.Err()
}

func (p *Postgres) Store(e Page) error {
	_, err := p.conn.Exec(context.Background(), `
	INSERT INTO page
		(page, timestamp, stable_version, homepage)
	VALUES
		($1, $2, $3, $4)
`, e.Page, e.T, e.StableVersion, e.Homepage)
	return err
}

func (p *Postgres) Known() ([]string, error) {
	var ps []string
	rows, err := p.conn.Query(context.Background(), `
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
	_, err = p.conn.Exec(context.Background(), `
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
	row := p.conn.QueryRow(context.Background(), `
		SELECT created, lastused, lastupdated, title
		FROM curated
		WHERE id=$1`,
		id,
	)
	cur := Curated{}
	if err := row.Scan(&cur.Created, &cur.LastUsed, &cur.LastUpdated, &cur.CustomTitle); err != nil {
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
	rows, err := p.conn.Query(context.Background(), `
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
func (p *Postgres) CuratedSetPages(id string, pages []string) error {
	tx, err := p.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	if _, err := tx.Exec(context.Background(), `DELETE FROM curated_pages WHERE curated_id=$1`, id); err != nil {
		return err
	}
	for _, p := range pages {
		if _, err := tx.Exec(context.Background(), `INSERT INTO curated_pages (curated_id, page) VALUES ($1, $2)`, id, p); err != nil {
			return err
		}
	}

	res, err := tx.Exec(context.Background(), `UPDATE curated SET lastupdated=now() WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		tx.Rollback(context.Background())
		return ErrCuratedNotFound
	}
	return tx.Commit(context.Background())
}

func (p *Postgres) CuratedSetUsed(id string) error {
	res, err := p.conn.Exec(context.Background(), `UPDATE curated SET lastused=now(), used=used+1 WHERE id=$1`, id)
	if res.RowsAffected() == 0 {
		return ErrCuratedNotFound
	}
	return err
}

func (p *Postgres) CuratedSetTitle(id, title string) error {
	res, err := p.conn.Exec(context.Background(), `UPDATE curated SET title=$2, lastupdated=now() WHERE id=$1`, id, title)
	if res.RowsAffected() == 0 {
		return ErrCuratedNotFound
	}
	return err
}

func (p *Postgres) UpdateViews(ctx context.Context) error {
	if _, err := p.conn.Exec(ctx, `REFRESH MATERIALIZED VIEW CONCURRENTLY current`); err != nil {
		return err
	}

	return nil
}
