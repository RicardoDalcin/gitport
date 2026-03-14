package repos

import (
	"gitport/internal/db"
	"gitport/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type RepositoryItem struct {
	Name          string            `json:"name"`
	Path          string            `json:"path"`
	Visibility    models.Visibility `json:"visibility"`
	CreatedAt     time.Time         `json:"createdAt"`
	Description   string            `json:"description,omitempty"`
	DefaultBranch string            `json:"defaultBranch,omitempty"`
}

type ListReposResponse struct {
	Repositories []*models.Repository `json:"repositories"`
}

func ListReposHandler(c *gin.Context, db *db.Client) {
	repos, err := db.ListRepositories()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(repos) == 0 {
		repos = []*models.Repository{}
	}

	c.JSON(200, ListReposResponse{Repositories: repos})
}
