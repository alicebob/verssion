DROP TABLE IF EXISTS page;

CREATE TABLE page
    ( title text NOT NULL
    , timestamp timestamptz NOT NULL
    , stable_version text NOT NULL 
    );
CREATE INDEX page_title ON page (title, timestamp);
