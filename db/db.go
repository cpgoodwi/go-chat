package db

import (
	"database/sql"
	"errors"
	"log"

	_ "modernc.org/sqlite" // TODO: change this to Pocketbase or MySQL or Postgres for production

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/file"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	// TODO: put db in separate file or access from storage solution and not in project source
	db, err := sql.Open("sqlite", "go-chat.db")
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}

func (d *Database) MigrateDB(mType string) error {
	instance, err := sqlite.WithInstance(d.db, &sqlite.Config{})
	if err != nil {
		return err
	}

	fSrc, err := (&file.File{}).Open("./db/migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("file", fSrc, "sqlite", instance)
	if err != nil {
		log.Println("Error from NewWithInstance")
		return err
	}

	switch mType {
	case "UP":
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("No new migrations")
				return nil
			}

			log.Println("Error from UP")

			return err
		}
	case "DOWN":
		if err := m.Down(); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("No new migrations")
				return nil
			}

			log.Println("Error from DOWN")

			return err
		}
	default:
		return errors.New("only migrate UP or DOWN")
	}

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}
