package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andreiz53/web-scraper/api"
	"github.com/andreiz53/web-scraper/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Loaded environment file")

	store, err := db.NewSqliteStorage(os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}
	app := api.NewServer(os.Getenv("LISTEN_ADDR"), store)
	app.Start()
}
