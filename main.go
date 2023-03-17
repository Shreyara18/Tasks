package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	apiURL       = "https://randomuser.me/api/"
	dbHost       = "localhost"
	dbPort       = "5432"
	dbUser       = "postgres"
	dbPassword   = "1234"
	dbName       = "TASK"
	defaultLimit = 5
)

type userRecord struct {
	ID        int
	Gender    string
	Title     string
	FirstName string
	LastName  string
	Street    string
	City      string
	State     string
	Country   string
	Postcode  int
	Email     string
	Phone     string
	Picture   string
}

var db *sql.DB

func init() {
	var err error
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database")
}

func main() {
	var err error
	db, err = sql.Open("postgres", "postgres://user:password@localhost/mydatabase?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err.Error())
	}
	defer db.Close()

	err = fetchAndStoreUserRecords()
	if err != nil {
		log.Fatalf("Failed to fetch and store user records: %s", err.Error())
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}

	// API endpoint to fetch user records
	router.GET("/api/users", fetchUserRecords)

	// Page to display user records and export them as CSV
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})
	router.GET("/export", exportUserRecords)

	// Page to edit user record
	router.GET("/edit/:id", editUserRecord)
	router.POST("/edit/:id", saveUserRecord)

	router.Run(":8080")
}

func fetchUserRecords(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	rows, err := db.Query("SELECT id, gender, title, first_name, last_name, street, city, state, country, postcode, email, phone, picture FROM user_records LIMIT $1 OFFSET $2",
		limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user records"})
		return
	}
	defer rows.Close()

	records := make([]userRecord, 0)
	for rows.Next() {
		record := userRecord{}
		err := rows.Scan(&record.ID, &record.Gender, &record.Title, &record.FirstName, &record.LastName, &record.Street,
			&record.City, &record.State, &record.Country, &record.Postcode, &record.Email, &record.Phone, &record.Picture)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user records"})
			return
		}
		records = append(records, record)
	}
	c.JSON(http.StatusOK, records)
}

func exportUserRecords(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/csv")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=user_records.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	rows, err := db.Query("SELECT gender, title, first_name, last_name, street, city, state, country, postcode, email, phone FROM user_records")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user records"})
		return
	}
	defer rows.Close()

	// Write header row
	err = writer.Write([]string{"Gender", "Title", "First Name", "Last Name", "Street", "City", "State", "Country", "Postcode", "Email", "Phone"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export user records"})
		return
	}

	// Write data rows
	for rows.Next() {
		var gender, title, firstName, lastName, street, city, state, country, email, phone string
		var postcode int
		err := rows.Scan(&gender, &title, &firstName, &lastName, &street, &city, &state, &country, &postcode, &email, &phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export user records"})
			return
		}
		err = writer.Write([]string{gender, title, firstName, lastName, street, city, state, country, strconv.Itoa(postcode), email, phone})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export user records"})
			return
		}
	}
}
func editUserRecord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	row := db.QueryRow("SELECT id, gender, title, first_name, last_name, street, city, state, country, postcode, email, phone FROM user_records WHERE id = $1", id)

	record := userRecord{}
	err = row.Scan(&record.ID, &record.Gender, &record.Title, &record.FirstName, &record.LastName, &record.Street,
		&record.City, &record.State, &record.Country, &record.Postcode, &record.Email, &record.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user record"})
		return
	}

	c.HTML(http.StatusOK, "edit.tmpl", gin.H{
		"record": record,
	})
}
func saveUserRecord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	gender := c.PostForm("gender")
	title := c.PostForm("title")
	firstName := c.PostForm("first_name")
	lastName := c.PostForm("last_name")
	street := c.PostForm("street")
	city := c.PostForm("city")
	state := c.PostForm("state")
	country := c.PostForm("country")
	postcode, _ := strconv.Atoi(c.PostForm("postcode"))
	email := c.PostForm("email")
	phone := c.PostForm("phone")

	_, err = db.Exec("UPDATE user_records SET gender = $1, title = $2, first_name = $3, last_name = $4, street = $5, city = $6, state = $7,country = $8, postcode = $9, email = $10, phone = $11 WHERE id = $12",
		gender, title, firstName, lastName, street, city, state, country, postcode, email, phone, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user record"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}
func fetchAndStoreUserRecords() error {
	for i := 0; i < 100; i++ {
		resp, err := http.Get("https://randomuser.me/api/")
		if err != nil {
			return fmt.Errorf("failed to fetch user record: %s", err.Error())
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read user record response body: %s", err.Error())
		}

		record := userRecord{}
		err = json.Unmarshal(body, &record)
		if err != nil {
			return fmt.Errorf("failed to unmarshal user record: %s", err.Error())
		}

		// Store record in Postgres DB
		_, err = db.Exec("INSERT INTO user_records (gender, title, first_name, last_name, street, city, state, country, postcode, email, phone, picture_large) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			record.Gender, record.Title, record.FirstName, record.LastName, record.Street, record.City, record.State, record.Country, record.Postcode, record.Email, record.Phone, record.Picture)
		if err != nil {
			return fmt.Errorf("failed to insert user record: %s", err.Error())
		}
	}
	return nil
}
