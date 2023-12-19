PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS sites (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    default_data_source TEXT
);

CREATE TABLE IF NOT EXISTS collections (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    metadata TEXT DEFAULT '{}',
    schema TEXT DEFAULT '[]',
    form_schema TEXT DEFAULT '{}',
    site_id TEXT REFERENCES sites(id) ON DELETE CASCADE,
    UNIQUE (name, site_id)
);

CREATE TABLE IF NOT EXISTS data_sources (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    provider TEXT NOT NULL,
    config TEXT DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS settings (
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    "group" TEXT NOT NULL,
    site_id TEXT REFERENCES sites(id) ON DELETE CASCADE,
    UNIQUE (key, "group", site_id)
);