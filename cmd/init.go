package cmd

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"time"
)

var commitHash string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Not .env file found")
	}

	fmt.Printf("Build Time: %s\n", time.Now().Format(time.RFC3339))
	fmt.Printf("Version: %s\n", commitHash)
}
