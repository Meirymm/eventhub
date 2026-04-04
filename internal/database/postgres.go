package database

import (
	"database/sql"
	"eventhub/config"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) *sql.DB {
	// Код снова берет данные из конфига (из .env)
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Database is unreachable! Check Docker: ", err)
	}

	log.Println("Database connected successfully!")
	return db
}
