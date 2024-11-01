package sqldb

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var Db *sql.DB

func InitDB() (*sql.DB, error) {
	var err error
	log.Println("dsnnn", viper.GetString("dsn"))
	Db, err = sql.Open("mysql", viper.GetString("dsn"))
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
		return nil, err
	}

	// Create the database if it does not exist
	_, err = Db.Exec("CREATE DATABASE IF NOT EXISTS cetec")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
		return nil, err
	}

	// Switch to the cetec database
	_, err = Db.Exec(viper.GetString("Database"))
	if err != nil {
		log.Fatalf("Failed to select database: %v", err)
		return nil, err
	}

	// Create tables and insert initial data if they do not exist
	createTablesAndInsertData()
	return Db, nil
}

func createTablesAndInsertData() {
	tableQueries := []string{
		// Create Person table
		`CREATE TABLE IF NOT EXISTS person (
            id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255),
            age INT
        );`,

		// Insert initial data into Person table
		`INSERT IGNORE INTO person(id, name, age) VALUES
            (1, "mike", 31),
            (2, "John", 20),
            (3, "Joseph", 20);`,

		// Create Phone table
		`CREATE TABLE IF NOT EXISTS phone (
            id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
            number VARCHAR(255),
            person_id INT
        );`,

		// Insert initial data into Phone table
		`INSERT IGNORE INTO phone(id, person_id, number) VALUES
            (1, 1, "444-444-4444"),
            (8, 2, "123-444-7777"),
            (3, 3, "445-222-1234");`,

		// Create Address table
		`CREATE TABLE IF NOT EXISTS address (
            id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
            city VARCHAR(255),
            state VARCHAR(255),
            street1 VARCHAR(255),
            street2 VARCHAR(255),
            zip_code VARCHAR(255)
        );`,

		// Insert initial data into Address table
		`INSERT IGNORE INTO address(id, city, state, street1, street2, zip_code) VALUES
            (1, "Eugene", "OR", "111 Main St", "", "98765"),
            (2, "Sacramento", "CA", "432 First St", "Apt 1", "22221"),
            (3, "Austin", "TX", "213 South 1st St", "", "78704");`,

		// Create Address_join table
		`CREATE TABLE IF NOT EXISTS address_join (
            id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
            person_id INT,
            address_id INT
        );`,

		// Insert initial data into Address_join table
		`INSERT IGNORE INTO address_join(id, person_id, address_id) VALUES
            (1, 1, 3),
            (2, 2, 1),
            (3, 3, 2);`,
	}

	for _, query := range tableQueries {
		if _, err := Db.Exec(query); err != nil {
			log.Fatalf("Failed to execute query: %v\nQuery: %s", err, query)
		}
	}

	fmt.Println("Database and tables are set up successfully.")
}
