package main

import (
	"fmt"
	"log"
	"os"
	"steam-trading/crawler/buff"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cwd, err := os.Getwd()
	fmt.Printf("Current working directory: %s\n", cwd)

	fmt.Printf("Buff Secret: %s\n", os.Getenv("BUFF_SECRET"))

	buffCrawler := &buff.BuffCrawler{}

}
