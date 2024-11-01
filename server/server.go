package server

import (
	"database/sql"
	"net/http"

	db "infilon_project/sqldb"

	"github.com/gin-gonic/gin"
)

type PersonInfo struct {
	Name        string `json:"name"`
	Age         string `json:"age"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	State       string `json:"state"`
	Street1     string `json:"street1"`
	Street2     string `json:"street2"`
	ZipCode     string `json:"zip_code"`
}

func GetPersonInfo(c *gin.Context) {
	personID := c.Param("person_id")
	var info PersonInfo

	query := `
    SELECT p.name,p.age, ph.number, a.city, a.state, a.street1, a.street2, a.zip_code
    FROM person p
    JOIN phone ph ON p.id = ph.person_id
    JOIN address_join aj ON p.id = aj.person_id
    JOIN address a ON aj.address_id = a.id
    WHERE p.id = ?
    `

	err := db.Db.QueryRow(query, personID).Scan(&info.Name, &info.Age, &info.PhoneNumber, &info.City, &info.State, &info.Street1, &info.Street2, &info.ZipCode)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, info)
}

func CreatePerson(c *gin.Context) {
	var newPerson PersonInfo
	if err := c.ShouldBindJSON(&newPerson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Start a transaction
	tx, err := db.Db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Insert person
	personRes, err := tx.Exec("INSERT INTO person(name, age) VALUES(?, ?)", newPerson.Name, newPerson.Age) // Assume age 25 for example
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert person"})
		return
	}
	personID, _ := personRes.LastInsertId()

	// Insert phone
	_, err = tx.Exec("INSERT INTO phone(number, person_id) VALUES(?, ?)", newPerson.PhoneNumber, personID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert phone"})
		return
	}

	// Insert address
	addressRes, err := tx.Exec("INSERT INTO address(city, state, street1, street2, zip_code) VALUES(?, ?, ?, ?, ?)",
		newPerson.City, newPerson.State, newPerson.Street1, newPerson.Street2, newPerson.ZipCode)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert address"})
		return
	}
	addressID, _ := addressRes.LastInsertId()

	// Insert into address_join
	_, err = tx.Exec("INSERT INTO address_join(person_id, address_id) VALUES(?, ?)", personID, addressID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate address"})
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person created successfully"})
}

func Connect() {
	r := gin.Default()
	r.GET("/person/:person_id/info", GetPersonInfo)
	r.POST("/person/create", CreatePerson)

	r.Run(":8080")
}
