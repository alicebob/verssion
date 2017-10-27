DROP TABLE IF EXISTS page CASCADE;

CREATE TABLE page
    ( page text NOT NULL
    , timestamp timestamptz NOT NULL
    , stable_version text NOT NULL 
    );
CREATE INDEX page_page ON page (page, timestamp);

CREATE VIEW updates
AS SELECT page, timestamp, stable_version
    FROM (
        SELECT page, timestamp, stable_version, lag(stable_version) OVER (
            PARTITION BY page ORDER BY timestamp
        ) AS prev
        FROM page
    ) sub
    WHERE prev IS NULL OR stable_version <> prev;
