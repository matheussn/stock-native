CREATE TABLE institution (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    name             TEXT NOT NULL,
    address          TEXT,
    cnpj             TEXT,
    responsible_name TEXT,
    phone            TEXT,
    is_active        INTEGER NOT NULL DEFAULT 1,
    created_at       TEXT NOT NULL DEFAULT (datetime('now'))
);

ALTER TABLE movement
    ADD COLUMN institution_id INTEGER REFERENCES institution(id);
