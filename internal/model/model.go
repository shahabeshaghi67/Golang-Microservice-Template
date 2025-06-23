package model

import (
	"time"
)

// WithTimestamps is a struct that can be embedded in a model in order to add timestamps for
// the update and creation time to a model.
type WithTimestamps struct {
	CreatedAt time.Time `json:"created_at" yaml:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
}
