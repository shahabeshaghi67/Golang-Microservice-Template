package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/model"
	"github.com/uptrace/bun"
)

type User struct {
	userRepo userRepository
}

func NewUser(
	userRepo userRepository,
) *User {
	return &User{
		userRepo: userRepo,
	}
}

func (u *User) Create(ctx context.Context, user model.User) (model.User, error) {
	err := u.userRepo.InTx(ctx, func(txCtx context.Context, tx bun.Tx) error {
		users, err := u.userRepo.Save(txCtx, tx, user)
		if err != nil {
			return err
		}
		if len(users) == 0 {
			return errors.New("no users saved")
		}
		user = users[0]
		return nil
	})
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *User) GetByID(ctx context.Context, id string) (model.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return model.User{}, errors.New("invalid user ID format")
	}

	var user model.User
	err = u.userRepo.InTx(ctx, func(txCtx context.Context, tx bun.Tx) error {
		var err error
		user, err = u.userRepo.GetByID(txCtx, tx, userID)
		return err
	})
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
