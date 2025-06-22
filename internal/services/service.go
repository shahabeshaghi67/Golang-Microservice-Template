package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/model"
)

type userRepository interface {
	InTx(ctx context.Context, operation func(context.Context, bun.Tx) error) error
	Save(context.Context, bun.Tx, ...model.User) ([]model.User, error)
	GetByID(ctx context.Context, tx bun.Tx, userID uuid.UUID) (model.User, error)
}
