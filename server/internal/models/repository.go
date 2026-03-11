package models

import "time"

type Visibility string

const (
	VisibilityPrivate  Visibility = "private"
	VisibilityPublic   Visibility = "public"
	VisibilityInternal Visibility = "internal"
)

type Repository struct {
	Name          string     `json:"name"`
	Path          string     `json:"path"`
	Visibility    Visibility `json:"visibility"`
	CreatedAt     time.Time  `json:"createdAt"`
	Description   string     `json:"description,omitempty"`
	DefaultBranch string     `json:"defaultBranch,omitempty"`
}
