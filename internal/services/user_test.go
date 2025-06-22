package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/model"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/services"
)

func setupUserService() (
	*services.User,
	*mockUserRepository,
) {
	mockUserRepo := &mockUserRepository{}
	sso := services.NewUser(mockUserRepo)
	return sso, mockUserRepo
}

func TestUser_Create(t *testing.T) {
	testUser := testUsers[0]
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		svc, userRepo := setupUserService()
		userRepo.On("Save", ctx, testTx, testUser).Return([]model.User{testUser}, nil)

		user, err := svc.Create(ctx, testUser)
		require.NoError(t, err)
		assert.Equal(t, testUser, user)
		userRepo.AssertExpectations(t)
	})

	t.Run("save error", func(t *testing.T) {
		svc, userRepo := setupUserService()
		userRepo.On("Save", ctx, testTx, testUser).Return([]model.User{}, assert.AnError)

		_, err := svc.Create(ctx, testUser)
		require.Error(t, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("no users saved", func(t *testing.T) {
		svc, userRepo := setupUserService()
		userRepo.On("Save", ctx, testTx, testUser).Return([]model.User{}, nil)

		_, err := svc.Create(ctx, testUser)
		require.Error(t, err)
		assert.EqualError(t, err, "no users saved")
		userRepo.AssertExpectations(t)
	})
}

func TestUser_GetByID(t *testing.T) {
	testUser := testUsers[0]
	ctx := context.Background()
	idStr := testUser.ID.String()

	t.Run("success", func(t *testing.T) {
		svc, userRepo := setupUserService()
		userRepo.On("GetByID", ctx, testTx, testUser.ID).Return(testUser, nil)

		user, err := svc.GetByID(ctx, idStr)
		require.NoError(t, err)
		assert.Equal(t, testUser, user)
		userRepo.AssertExpectations(t)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		svc, userRepo := setupUserService()
		_, err := svc.GetByID(ctx, "not-a-uuid")
		require.Error(t, err)
		assert.EqualError(t, err, "invalid user ID format")
		userRepo.AssertExpectations(t)
	})

	t.Run("get error", func(t *testing.T) {
		svc, userRepo := setupUserService()
		userRepo.On("GetByID", ctx, testTx, testUser.ID).Return(model.User{}, assert.AnError)

		_, err := svc.GetByID(ctx, idStr)
		require.Error(t, err)
		assert.EqualError(t, err, "assert.AnError general error for testing")
		userRepo.AssertExpectations(t)
	})
}
