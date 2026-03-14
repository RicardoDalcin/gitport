package git

import (
	"fmt"
	"io"
	"strings"

	"gitport/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/format/packfile"
)

func InfoRefsHandler(c *gin.Context, dbClient *db.Client) {
	repoNameWithSuffix := c.Param("name")
	service := c.Query("service")

	if !strings.HasSuffix(repoNameWithSuffix, ".git") {
		c.JSON(400, gin.H{"error": "repository name must end with .git"})
		return
	}

	repoName := strings.TrimSuffix(repoNameWithSuffix, ".git")

	if service != "git-upload-pack" {
		c.JSON(400, gin.H{"error": "unsupported service"})
		return
	}

	repo, err := dbClient.GetRepository(repoName)
	if err != nil {
		c.JSON(404, gin.H{"error": "repository not found"})
		return
	}

	r, err := git.PlainOpen(repo.Path)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to open repository: %v", err)})
		return
	}

	c.Header("Content-Type", fmt.Sprintf("application/x-%s-advertisement", service))
	c.Header("Cache-Control", "no-cache")

	writer := c.Writer

	serviceLine := fmt.Sprintf("# service=%s\n", service)
	totalLength := 4 + len(serviceLine)
	fmt.Fprintf(writer, "%04x%s", totalLength, serviceLine)

	writer.WriteString("0000")

	refs, err := r.References()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to get references: %v", err)})
		return
	}
	defer refs.Close()

	for {
		ref, err := refs.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("failed to read reference: %v", err)})
			return
		}

		if ref.Name() == plumbing.HEAD && ref.Type() == plumbing.SymbolicReference {
			continue
		}

		hash := ref.Hash().String()
		refName := ref.Name().String()
		line := fmt.Sprintf("%s %s\n", hash, refName)

		length := len(line) + 4
		writer.WriteString(fmt.Sprintf("%04x%s", length, line))
	}

	writer.WriteString("0000")
}

func UploadPackHandler(c *gin.Context, dbClient *db.Client) {
	repoNameWithSuffix := c.Param("name")

	if !strings.HasSuffix(repoNameWithSuffix, ".git") {
		c.JSON(400, gin.H{"error": "repository name must end with .git"})
		return
	}

	repoName := strings.TrimSuffix(repoNameWithSuffix, ".git")

	repo, err := dbClient.GetRepository(repoName)
	if err != nil {
		c.JSON(404, gin.H{"error": "repository not found"})
		return
	}

	r, err := git.PlainOpen(repo.Path)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to open repository: %v", err)})
		return
	}

	c.Header("Content-Type", "application/x-git-upload-pack-result")

	requestData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed to read request: %v", err)})
		return
	}
	_ = requestData

	enc := packfile.NewEncoder(c.Writer, r.Storer, false)

	refs, err := r.References()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to get references: %v", err)})
		return
	}
	defer refs.Close()

	var hashes []plumbing.Hash
	for {
		ref, err := refs.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("failed to read reference: %v", err)})
			return
		}
		if ref.Type() == plumbing.HashReference {
			hashes = append(hashes, ref.Hash())
		}
	}

	_, err = enc.Encode(hashes, 0)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to encode packfile: %v", err)})
		return
	}
}
