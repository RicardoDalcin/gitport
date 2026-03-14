package main

import (
	"log"

	"gitport/internal/db"
	"gitport/internal/handlers/git"
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

	router.GET("/repos/:name/info/refs", func(c *gin.Context) {
		git.InfoRefsHandler(c, db)
	})
	router.POST("/repos/:name.git/git-upload-pack", func(c *gin.Context) {
		git.UploadPackHandler(c, db)
	})

	router.Run(":8080")
}
