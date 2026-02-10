package database

import (
	"errors"
	"fmt"

	"github.com/Saikatdeb12/TodoApp2/internal/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

type SSLMODETYPE string

const sslmode SSLMODETYPE = "disable"

func migrateUp(db *sqlx.DB) error {
	driver, driErr := postgres.WithInstance(db.DB, &postgres.Config{})

	if driErr != nil {
		return driErr
	}
	m, migErr := migrate.NewWithDatabaseInstance("file://internal/database/migrations", "postgres", driver)

	if migErr != nil {
		return migErr
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func Connect() error {
	// connStr := utils.GoDotEnvVariable("POSTGRESQL_URL")

	db_user := utils.GoDotEnvVariable("DB_USER")
	db_password := utils.GoDotEnvVariable("DB_PASSWORD")
	db_name := utils.GoDotEnvVariable("DB_NAME")
	db_host := utils.GoDotEnvVariable("DB_HOST")
	db_port := utils.GoDotEnvVariable("DB_PORT")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", db_host, db_port, db_user, db_password, db_name, sslmode)
	var err error
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	fmt.Println("Connected to PostgreSQL")
	return migrateUp(DB)
}
