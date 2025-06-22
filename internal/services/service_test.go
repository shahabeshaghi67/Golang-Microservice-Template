//go:generate mockery --all --case snake --exported --output ../mocks
package services_test

import (
	"context"
	"os"
	"testing"

	"github.com/uptrace/bun"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/mocks"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/model"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/pkg/testutils"
)

var (
	testTx    = bun.Tx{}
	testUsers []model.User
)

// TestMain prepares data required in various repository tests.
func TestMain(m *testing.M) {
	loadTestData()
	os.Exit(m.Run())
}

func loadTestData() {
	testutils.LoadJSON(&testUsers, "testdata", "users.json")
}

type mockUserRepository struct {
	mocks.UserRepository
}

func (m *mockUserRepository) InTx(ctx context.Context, operation func(txCtx context.Context, tx bun.Tx) error) error {
	return operation(ctx, testTx)
}
