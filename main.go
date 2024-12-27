package main

import (
	"github.com/joho/godotenv"
	"github.com/oswgg/migrator/cmd"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cmd.Execute()
}
