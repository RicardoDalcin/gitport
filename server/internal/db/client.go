package db

import (
	"database/sql"
	"fmt"
	"time"

	"gitport/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	db *sql.DB
}

func NewClient(dbPath string) (*Client, error) {
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	client := &Client{db: db}

	// Initialize database schema
	if err := client.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return client, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) initSchema() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS repositories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		path TEXT NOT NULL UNIQUE,
		visibility TEXT NOT NULL CHECK(visibility IN ('private', 'public', 'internal')),
		description TEXT,
		default_branch TEXT,
		created_at TEXT NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_repositories_name ON repositories(name);
	CREATE INDEX IF NOT EXISTS idx_repositories_visibility ON repositories(visibility);
	`

	if _, err := c.db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func (c *Client) CreateRepository(repo *models.Repository) error {
	query := `
		INSERT INTO repositories (name, path, visibility, description, default_branch, created_at)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	// Store time as RFC3339 string
	createdAtStr := repo.CreatedAt.Format(time.RFC3339)

	_, err := c.db.Exec(
		query,
		repo.Name,
		repo.Path,
		string(repo.Visibility),
		repo.Description,
		repo.DefaultBranch,
		createdAtStr,
	)

	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	return nil
}

// GetRepository retrieves a repository by name
func (c *Client) GetRepository(name string) (*models.Repository, error) {
	query := `
		SELECT name, path, visibility, description, default_branch, created_at
		FROM repositories
		WHERE name = ?
	`

	var repo models.Repository
	var visibilityStr string
	var createdAtStr string

	err := c.db.QueryRow(query, name).Scan(
		&repo.Name,
		&repo.Path,
		&visibilityStr,
		&repo.Description,
		&repo.DefaultBranch,
		&createdAtStr,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("repository %s not found", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	repo.Visibility = models.Visibility(visibilityStr)
	repo.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	return &repo, nil
}

// ListRepositories retrieves all repositories
func (c *Client) ListRepositories() ([]*models.Repository, error) {
	query := `
		SELECT name, path, visibility, description, default_branch, created_at
		FROM repositories
		ORDER BY created_at DESC
	`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}
	defer rows.Close()

	var repos []*models.Repository
	for rows.Next() {
		var repo models.Repository
		var visibilityStr string
		var createdAtStr string

		if err := rows.Scan(
			&repo.Name,
			&repo.Path,
			&visibilityStr,
			&repo.Description,
			&repo.DefaultBranch,
			&createdAtStr,
		); err != nil {
			return nil, fmt.Errorf("failed to scan repository: %w", err)
		}

		repo.Visibility = models.Visibility(visibilityStr)
		repo.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}

		repos = append(repos, &repo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating repositories: %w", err)
	}

	return repos, nil
}

// UpdateRepository updates an existing repository
func (c *Client) UpdateRepository(repo *models.Repository) error {
	query := `
		UPDATE repositories
		SET path = ?, visibility = ?, description = ?, default_branch = ?
		WHERE name = ?
	`

	result, err := c.db.Exec(
		query,
		repo.Path,
		string(repo.Visibility),
		repo.Description,
		repo.DefaultBranch,
		repo.Name,
	)

	if err != nil {
		return fmt.Errorf("failed to update repository: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("repository %s not found", repo.Name)
	}

	return nil
}

// DeleteRepository deletes a repository by name
func (c *Client) DeleteRepository(name string) error {
	query := `DELETE FROM repositories WHERE name = ?`

	result, err := c.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("repository %s not found", name)
	}

	return nil
}

// RepositoryExists checks if a repository exists by name
func (c *Client) RepositoryExists(name string) (bool, error) {
	query := `SELECT COUNT(*) FROM repositories WHERE name = ?`

	var count int
	err := c.db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check repository existence: %w", err)
	}

	return count > 0, nil
}
