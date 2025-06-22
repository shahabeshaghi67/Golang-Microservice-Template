//go:generate mockery --all --case snake --exported --output ../mocks
package handlers_test

import (
	"os"
	"testing"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/model"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/pkg/testutils"
)

var (
	testUser     []model.User
	testUserJSON string
)

func TestMain(m *testing.M) {
	loadTestData()
	os.Exit(m.Run())
}

func loadTestData() {
	testUserJSON = testutils.LoadJSON(&testUser, "testdata", "users.json")
}
