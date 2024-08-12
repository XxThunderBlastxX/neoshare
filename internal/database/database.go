package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/XxThunderBlastxX/neoshare/internal/config"
)

// ConnectDB connects to a postgresql database
func ConnectDB(dbConfig *config.DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("☹️Failed to connect to database: %w", err)
	}

	fmt.Println("🎉 Database connected !!!")

	return db, nil
}
