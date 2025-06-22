package repositories_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/fortytw2/dockertest"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/config"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/model"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/pkg/database"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/pkg/testutils"
)

const (
	dummyCorrectID = "00000000-0000-0000-0000-000000000000"
)

var (
	testUsers []model.User
)

func TestMain(m *testing.M) {
	loadTestData()
	os.Exit(func() int {
		// start a db container if no database is present
		err := pingDatabase()
		if err != nil {
			container := startPostgresContainer()
			defer func() {
				if !testutils.ShouldKeepTestDBContainer() {
					container.Shutdown()
				}
			}()
			fmt.Printf("started postgres container with addr %s\n", container.Addr)
		}
		return m.Run()
	}())
}

func loadTestData() {
	testutils.LoadJSON(&testUsers, "testdata", "users.json")
}

func startPostgresContainer() *dockertest.Container {
	waitFunc := func(addr string) error {
		hostAndPort := strings.Split(addr, ":")
		os.Setenv("DB_HOST", hostAndPort[0])
		os.Setenv("DB_PORT", hostAndPort[1])
		err := pingDatabase()
		return err
	}
	databaseConfig := config.Load().Database
	container, err := dockertest.RunContainer(
		"postgres:13.7",
		strconv.Itoa(databaseConfig.Port),
		waitFunc,
		"-e", fmt.Sprintf("POSTGRES_DB=%s", databaseConfig.Name),
		"-e", fmt.Sprintf("POSTGRES_PASSWORD=%s", databaseConfig.Password),
		"--rm")
	if err != nil {
		panic(err)
	}
	return container
}

func pingDatabase() error {
	// reloading config to ping database, as host and port are different now
	conn := database.ConnectDB(config.Load().Database.TestCfg())
	defer conn.Close()
	return conn.Ping()
}
