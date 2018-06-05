package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var create = `
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE FUNCTION auto_set_updated() RETURNS trigger AS $auto_set_updated$
	BEGIN
		NEW.updated := NOW();
		RETURN NEW;
	END;
$auto_set_updated$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS collections (
	id BIGSERIAL PRIMARY KEY,
	parent_collection_id BIGINT REFERENCES collections(id) ON UPDATE CASCADE ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
	full_name TEXT UNIQUE,
	name TEXT NOT NULL,
	title TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS collections_created_idx ON features(created);
CREATE INDEX IF NOT EXISTS collections_updated_idx ON features(updated);
CREATE INDEX IF NOT EXISTS collections_deleted_idx ON features(deleted);

CREATE TRIGGER collections_auto_set_updated BEFORE UPDATE ON collections
	FOR EACH ROW EXECUTE PROCEDURE auto_set_updated();

CREATE TABLE IF NOT EXISTS features (
	id UUID PRIMARY KEY,
	collection_id BIGINT REFERENCES collections(id) NOT NULL ON UPDATE CASCADE ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
	geometry GEOMETRY(GEOMETRY, 4326) NOT NULL,
	properties JSONB NOT NULL,
	created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS features_collection_id_idx ON features(id);
CREATE INDEX IF NOT EXISTS features_geometry_idx ON features USING GIST(geometry);
CREATE INDEX IF NOT EXISTS collections_created_idx ON features(created);
CREATE INDEX IF NOT EXISTS collections_updated_idx ON features(updated);
CREATE INDEX IF NOT EXISTS collections_deleted_idx ON features(deleted);

CREATE TRIGGER features_auto_set_updated BEFORE UPDATE ON features
	FOR EACH ROW EXECUTE PROCEDURE auto_set_updated();
`

var drop = `
DROP TABLE features;
DROP TABLE collections;
`

// Migrate updates the database
func Migrate(db *sql.DB) error {
	sqlxDB := sqlx.NewDb(db, driverName)
	_, err := sqlxDB.Exec(create)

	return err
}

// Drop all the data
func Drop(db *sql.DB) error {
	sqlxDB := sqlx.NewDb(db, driverName)
	_, err := sqlxDB.Exec(drop)
	return err
}
