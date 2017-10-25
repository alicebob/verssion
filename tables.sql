DROP TABLE IF EXISTS page CASCADE;

CREATE TABLE page
    ( title text NOT NULL
    , timestamp timestamptz NOT NULL
    , stable_version text NOT NULL 
    );
CREATE INDEX page_title ON page (title, timestamp);

CREATE VIEW updates
AS SELECT title, timestamp, stable_version
    FROM (
        SELECT title, timestamp, stable_version, lag(stable_version) OVER (
            PARTITION BY title ORDER BY timestamp
        ) AS prev
        FROM page
    ) sub
    WHERE prev IS NULL OR stable_version <> prev;
