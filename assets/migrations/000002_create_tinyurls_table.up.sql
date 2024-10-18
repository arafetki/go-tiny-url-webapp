CREATE TABLE
    IF NOT EXISTS tinyurls (
        short TEXT PRIMARY KEY,
        long TEXT NOT NULL,
        expiry DATE NOT NULL,
        created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );