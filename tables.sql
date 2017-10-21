DROP TABLE IF EXISTS page;

CREATE TABLE page
    ( title text NOT NULL
    , revision int NOT NULL
    , timestamp timestamptz NOT NULL
    , stable_version text NOT NULL 
    , UNIQUE(title, revision)
    )
