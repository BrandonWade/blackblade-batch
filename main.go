package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BlueMonday/go-scryfall"
)

// var (
// 	db     *sqlx.DB
// 	logger *logrus.Logger
// )

// func init() {
// 	dbUsername := os.Getenv("DB_USERNAME")
// 	dbPassword := os.Getenv("DB_PASSWORD")
// 	dbDatabase := os.Getenv("DB_DATABASE")
// 	dbHost := os.Getenv("DB_HOST")
// 	dbPort := os.Getenv("DB_PORT")
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, dbDatabase)

// 	// Connect to MySQL
// 	var err error
// 	db, err = sqlx.Connect("mysql", dsn)
// 	if err != nil {
// 		log.Fatalf("error connecting to db: %s\n", err.Error())
// 	}

// 	logger = logrus.New()
// }

func main() {
	// defer db.Close()

	// cardRepository := repositories.NewCardRepository(db)
	// cardService := services.NewCardService(logger, cardRepository)
	// batchRunner := runner.NewBatchRunner(logger, cardService)

	// // Start the service to fetch cards from the Scryfall API
	// batchRunner.Run()

	ctx := context.Background()
	client, _ := scryfall.NewClient()

	result, err := client.ListCards(ctx, scryfall.ListCardsOptions{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("%v\n", len(result.Cards))
}
