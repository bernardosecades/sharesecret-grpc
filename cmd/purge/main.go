package main

import (
	_ "github.com/bernardosecades/sharesecret/cmd"
	"github.com/bernardosecades/sharesecret/internal/storage/mysql"

	"fmt"
	"log"
	"os"
)

func main() {

	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	secretRepository := mysql.NewMySQLSecretRepository(dbName, dbUser, dbPass, dbHost, dbPort)
	r, err := secretRepository.RemoveSecretsExpired()

	if err != nil {
		log.Fatal("Error to try to remove expired secrets", err)
	}

	fmt.Println("Secrets deleted:")
	fmt.Println(r)
}
