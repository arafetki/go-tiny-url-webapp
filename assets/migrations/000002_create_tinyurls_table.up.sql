CREATE TABLE
    IF NOT EXISTS tinyurls (
        short VARCHAR(7) PRIMARY KEY,
        long TEXT NOT NULL,
        expiry DATE NOT NULL,
        created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );