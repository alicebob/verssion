-- DROP TABLE IF EXISTS curated;
-- DROP TABLE IF EXISTS curated_pages;
-- DROP TABLE IF EXISTS page CASCADE;

CREATE TABLE page
    ( page text NOT NULL
    , timestamp timestamptz NOT NULL
    , stable_version text NOT NULL 
    , homepage text NOT NULL
    );
CREATE INDEX page_page ON page (page, timestamp);

CREATE VIEW updates
AS SELECT page, timestamp, stable_version, homepage
    FROM (
        SELECT page, timestamp, stable_version, homepage, lag(stable_version) OVER (
            PARTITION BY page ORDER BY timestamp
        ) AS prev
        FROM page
    ) sub
    WHERE prev IS NULL OR stable_version <> prev;

CREATE VIEW current
AS SELECT page, timestamp, stable_version, homepage
    FROM (
        SELECT *, rank() OVER (
            PARTITION BY page ORDER BY timestamp DESC
        )
        FROM updates
    ) sub
    WHERE rank=1;

CREATE TABLE curated
    ( id text NOT NULL UNIQUE
    , created timestamptz NOT NULL
    , used int NOT NULL default 0
    , lastused timestamptz NOT NULL
    , lastupdated timestamptz NOT NULL
    , title text NOT NULL default ''
    );

CREATE TABLE curated_pages
    ( curated_id text NOT NULL
    , page text NOT NULL
    , UNIQUE (curated_id, page)
    );
