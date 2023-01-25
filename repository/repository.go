package repository

import (
	"context"

	"github.com/jscastaneda-esp/rest-ws-go/models"
)

type Repository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	InsertPost(ctx context.Context, post *models.Post) error
	GetPostById(ctx context.Context, id string) (*models.Post, error)
	ListPosts(ctx context.Context, page uint64, rowsFetch uint64) ([]*models.Post, error)
	UpdatePost(ctx context.Context, post *models.Post) error
	DeletePost(ctx context.Context, id string, userId string) error
	Close() error
}

var implementation Repository

func SetRepository(repository Repository) {
	implementation = repository
}

func Close() error {
	return implementation.Close()
}
