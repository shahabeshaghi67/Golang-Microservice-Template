package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/model"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/repositories"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/pkg/testutils"
)

var runGetByID = func(repo *repositories.User, UserID uuid.UUID) (model.User, error) {
	var user model.User
	var err error
	var ctx = context.Background()
	err = repo.InTx(ctx, func(_ context.Context, tx bun.Tx) error {
		user, err = repo.GetByID(ctx, tx, UserID)
		return err
	})
	return user, err
}

func TestUser_GetByID(t *testing.T) {
	testUser := testUsers[0]
	conn, repo := setUpUserRepo()
	defer conn.Close()

	t.Run("fail uuid not found", func(t *testing.T) {
		dummyID := uuid.MustParse(dummyCorrectID)
		_, err := runGetByID(repo, dummyID)
		require.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		user, err := runGetByID(repo, testUser.ID)
		require.NoError(t, err)
		assert.Equal(t, user.Name, testUser.Name)
	})
}

var runSave = func(repo *repositories.User, sos ...model.User) ([]model.User, error) {
	var ctx = context.Background()
	var err error
	err = repo.InTx(ctx, func(_ context.Context, tx bun.Tx) error {
		sos, err = repo.Save(ctx, tx, sos...)
		return err
	})
	return sos, err
}

func TestUser_Save(t *testing.T) {
	conn, repo := setUpUserRepo()
	defer conn.Close()
	testUser := testUsers[0]
	testUser.CreatedAt, testUser.UpdatedAt = time.Time{}, time.Time{}

	t.Run("success create", func(t *testing.T) {
		testUser.ID = uuid.New()
		created, err := runSave(repo, testUser)
		require.NoError(t, err)
		assert.NotEqual(t, testUser.UpdatedAt, created[0].UpdatedAt)
		assert.NotEqual(t, testUser.CreatedAt, created[0].CreatedAt)
		assert.Equal(t, created[0].CreatedAt, created[0].UpdatedAt)
	})

	t.Run("success update", func(t *testing.T) {
		testUser.ID = uuid.New()
		created, err := runSave(repo, testUser)
		require.NoError(t, err)
		assert.NotEqual(t, testUser.UpdatedAt, created[0].UpdatedAt)
		assert.NotEqual(t, testUser.CreatedAt, created[0].CreatedAt)
		assert.Equal(t, created[0].CreatedAt, created[0].UpdatedAt)
		updated, err := runSave(repo, created[0])
		require.NoError(t, err)
		assert.NotEqual(t, created[0].UpdatedAt, updated[0].UpdatedAt)
		assert.Equal(t, created[0].CreatedAt, updated[0].CreatedAt)
	})
}

func setUpUserRepo() (*bun.DB, *repositories.User) {
	conn := testutils.ConnectTestDB()
	return conn, repositories.NewUser(conn)
}
