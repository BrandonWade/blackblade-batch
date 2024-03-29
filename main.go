package main

import (
	"fmt"
	"log"
	"os"
	"time"

	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/BrandonWade/blackblade-batch/clients"
	"github.com/BrandonWade/blackblade-batch/repositories"
	"github.com/BrandonWade/blackblade-batch/runner"
	"github.com/BrandonWade/blackblade-batch/services"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

var (
	baseURL string
	db      *sqlx.DB
	client  *scryfall.Client
	logger  *logrus.Logger
)

func init() {
	baseURL = os.Getenv("BASE_SCRYFALL_URL")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, dbDatabase)

	// Connect to MySQL
	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("error connecting to db: %s\n", err.Error())
	}

	client, err = scryfall.NewClient()
	if err != nil {
		log.Fatalf("error creating scryfall client: %s\n", err.Error())
	}

	db.SetConnMaxLifetime(time.Second)

	logger = logrus.New()
}

func main() {
	defer db.Close()

	scryfallClient := clients.NewScryfallClient(baseURL, logger, client)
	cardRepository := repositories.NewCardRepository(logger, db)
	cardService := services.NewCardService(logger, scryfallClient, cardRepository)
	batchRunner := runner.NewBatchRunner(logger, cardService)

	// Start the service to fetch cards from the Scryfall API
	batchRunner.Run()
}
