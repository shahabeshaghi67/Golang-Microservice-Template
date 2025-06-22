package handlers

import (
	"context"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/model"
)

type UserService interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	GetByID(ctx context.Context, id string) (model.User, error)
}
