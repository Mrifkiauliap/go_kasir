package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func InitDB() (*sql.DB, error) {
	// Connect to database
	psqlInfo := viper.GetString("DATABASE_URL")
	if psqlInfo == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(25)

	fmt.Println("Successfully connected to database!")
	return db, nil
}
