package main

import (
	"log"

	"gitport/internal/db"
	"gitport/internal/handlers/repos"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := db.NewClient("file:gitport.db?cache=shared&mode=memory")

	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	router := gin.Default()
	router.POST("/v1/repos", func(c *gin.Context) {
		repos.NewRepoHandler(c, db)
	})
	router.GET("/v1/repos", func(c *gin.Context) {
		repos.ListReposHandler(c, db)
	})
	router.Run(":8080")
}
