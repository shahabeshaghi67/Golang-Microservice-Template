package model

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/utils"
)

var (
	userValidator *validator.Validate
)

func init() {
	v := validator.New()
	v.RegisterTagNameFunc(utils.JsonTagName)
	userValidator = v
}

// User holds data for the user (this is a test model).
type User struct {
	ID   uuid.UUID `json:"id" bun:",pk" validate:"required" yaml:"id"`
	Name string    `json:"name" yaml:"name"`
	WithTimestamps
}

// Validate is used to check if struct fields meet data requirements.
func (so *User) Validate() error {
	return userValidator.Struct(so)
}
