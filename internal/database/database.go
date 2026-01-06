package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

type Service struct {
	Database *pgxpool.Pool
}

var (
	dbURI = os.Getenv("DATABASE_URI")
	//dbName = os.Getenv("DATABASE_NAME")
)

func New() *Service {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	return &Service{
		Database: db,
	}

}
