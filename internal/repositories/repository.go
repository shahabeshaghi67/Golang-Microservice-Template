package repositories

import (
	"context"

	"github.com/uptrace/bun"
)

type dbRepo struct {
	conn *bun.DB
}

func (r *dbRepo) InTx(ctx context.Context, operation func(context context.Context, tx bun.Tx) error) error {
	return r.conn.RunInTx(ctx, nil, operation)
}
