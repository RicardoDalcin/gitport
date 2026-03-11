package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file:gitport.db?cache=shared&mode=memory")

	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	fmt.Println("Database opened successfully")
	db.Exec(`
	CREATE TABLE IF NOT EXISTS repositories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		path TEXT NOT NULL,
		visibility TEXT NOT NULL,
		created_at TEXT NOT NULL
	)
	`)
	fmt.Println("Table created successfully")

	db.Exec(`INSERT INTO repositories (name, path, visibility, created_at) VALUES ('test', 'test', 'public', '2026-03-11 12:00:00')`)
	fmt.Println("Data inserted successfully")
	rows, err := db.Query("SELECT * FROM repositories")
	if err != nil {
		log.Fatalf("failed to query database: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var path string
		var visibility string
		var created_at string
		err = rows.Scan(&id, &name, &path, &visibility, &created_at)
		if err != nil {
			log.Fatalf("failed to scan row: %v", err)
		}
		fmt.Printf("id: %d, name: %s, path: %s, visibility: %s, created_at: %s\n", id, name, path, visibility, created_at)
	}

	fmt.Println("Hello, World!")
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})
	router.Run(":8080")
}
