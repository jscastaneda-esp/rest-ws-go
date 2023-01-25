package repository

import (
	"context"

	"github.com/jscastaneda-esp/rest-ws-go/models"
)

func InsertPost(ctx context.Context, post *models.Post) error {
	return implementation.InsertPost(ctx, post)
}

func GetPostById(ctx context.Context, id string) (*models.Post, error) {
	return implementation.GetPostById(ctx, id)
}

func ListPosts(ctx context.Context, page uint64, rowsFetch uint64) ([]*models.Post, error) {
	return implementation.ListPosts(ctx, page, rowsFetch)
}

func UpdatePost(ctx context.Context, post *models.Post) error {
	return implementation.UpdatePost(ctx, post)
}

func DeletePost(ctx context.Context, id string, userId string) error {
	return implementation.DeletePost(ctx, id, userId)
}
