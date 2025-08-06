package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

var DB *sql.DB // pointer to sql connection pool

func InitializeDB() {
	err := godotenv.Load() // loading .env file for db data
	if err != nil {
		log.Fatal("Error while opening .env file")
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// forming connection string
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s",
		dbHost, dbUser, dbPassword, dbPort, dbName)

	// opening connection via connection string
	db, err := sql.Open("sqlserver", connectionString)

	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	// Check connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	fmt.Println("✅ Connected to MSSQL database!")

	// Set connection pool settings
	db.SetMaxOpenConns(25)                 // Max number of open connections to the DB
	db.SetMaxIdleConns(10)                 // Max number of idle (unused) connections
	db.SetConnMaxLifetime(time.Minute * 5) // Max lifetime of a single connection

	DB = db
}
