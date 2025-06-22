package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/model"
)

// User is an interface for the repository of Users.
type User struct {
	dbRepo
}

// NewUser is a factory function that returns a new instance of the User Repository interface.
func NewUser(conn *bun.DB) *User {
	return &User{dbRepo{conn: conn}}
}

// GetUserByID returns a User by their IDs.
func (r *User) GetByID(ctx context.Context, tx bun.Tx, id uuid.UUID) (model.User, error) {
	var user model.User
	err := tx.NewSelect().
		Model(&user).
		Where("id = (?)", id).
		Scan(ctx)
	return user, err
}

// Save saves the given users.
func (r *User) Save(ctx context.Context, tx bun.Tx, Users ...model.User) ([]model.User, error) {
	if len(Users) == 0 {
		return nil, errors.New("no users provided to save")
	}
	res := []model.User{}
	res = append(res, Users...)
	_, err := tx.NewInsert().
		Model(&res).
		On("CONFLICT (id) DO UPDATE").
		ExcludeColumn("created_at").
		Value("updated_at", "now()").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save users")
	}
	return res, nil
}
