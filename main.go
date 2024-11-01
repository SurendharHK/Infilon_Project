package main

import (
	"infilon_project/server"
	"infilon_project/sqldb"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func main() {

	var err error

	err = LoadConfig()

	if err != nil {
		log.Fatalf("Failed to initialize configurations: %v", err)
	}

	// Initialize the database and check for errors
	sqldb.Db, err = sqldb.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := sqldb.Db.Close(); err != nil {
			log.Printf("Error closing the database: %v", err)
		}
	}()

	// Start the server
	server.Connect()

}

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return err
	}
	return nil
}
