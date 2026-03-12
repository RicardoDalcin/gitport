package repos

import (
	"gitport/internal/db"
	"gitport/internal/models"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type NewRepoBody struct {
	Name        string            `json:"name"`
	Visibility  models.Visibility `json:"visibility"`
	Description string            `json:"description"`
}

func NewRepoHandler(c *gin.Context, db *db.Client) {
	body := NewRepoBody{}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	basePath := filepath.Join(os.Getenv("GIT_REPO_BASE_DIR"))
	path := filepath.Join(basePath, "repos", body.Name)

	err := db.CreateRepository(&models.Repository{
		Path:          path,
		Name:          body.Name,
		Visibility:    body.Visibility,
		Description:   body.Description,
		CreatedAt:     time.Now(),
		DefaultBranch: "main",
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Repository created successfully"})
}
