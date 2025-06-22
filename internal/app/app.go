package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	"github.com/uptrace/bun"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/api"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/config"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/handlers"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/repositories"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/services"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/pkg/database"
)

// App is the main Application interface.
type App struct {
	config           config.Application
	dbConn           *bun.DB
	logger           log.Logger
	router           *mux.Router
	fakeDependencies bool
	handlers         map[string]http.Handler
	routes           []*api.API
}

// NewApp creates a new application.
func NewApp(config config.Application, logger log.Logger, fakeDependencies bool) *App {
	return &App{
		config:           config,
		dbConn:           &bun.DB{},
		logger:           logger,
		fakeDependencies: fakeDependencies,
		router:           mux.NewRouter(),
	}
}

// Wire connects layers together(repo, services, handlers).
func (a *App) Wire() (*mux.Router, error) {

	err := a.dependenciesConnections()
	if err != nil {
		return nil, err
	}

	// Repositories
	userRepo := repositories.NewUser(a.dbConn)

	// Services
	userService := services.NewUser(
		userRepo,
	)

	// Handlers
	a.handlers = map[string]http.Handler{
		"users.create":  handlers.NewUserHandler(userService, a.logger).CreateUser(),
		"users.getByID": handlers.NewUserHandler(userService, a.logger).GetUserByID(),
	}

	router := mux.NewRouter()
	router = a.addHealthRoute(router)
	a.router = a.addBusinessRoutes(router)
	return a.router, nil
}

// addHealthRoute adds public health route.
func (a *App) addHealthRoute(router *mux.Router) *mux.Router {
	router.Handle("/health", healthHandler("golang-microservice-template"))
	return router
}

// addHealthRoute adds protected business routes on top.
func (a *App) addBusinessRoutes(router *mux.Router) *mux.Router {
	router.PathPrefix("/").Handler(newBusinessRouter(a.generateRoutes()))
	return router
}

// newBusinessRouter provisions router with a business logic handlers.
func newBusinessRouter(routes []*api.API) *mux.Router {
	router := mux.NewRouter().SkipClean(true)

	// not found route
	router.NotFoundHandler = notFoundHandler()

	// business routes
	for _, route := range routes {
		router.Handle(route.Path, route).Methods(route.Method)
	}
	return router
}

// initiate dependencies connections here.
func (a *App) dependenciesConnections() error {
	if a.fakeDependencies {
		return nil
	}

	a.dbConn = database.NewDBConnection(a.config.Database, a.logger)

	if a.config.Database.MigrateUp {
		_ = a.logger.Log("msg", "start db migration")
		err := database.MigrateUp(context.Background(), a.dbConn)
		if err != nil {
			return err
		}
		_ = a.logger.Log("msg", "db migrations successfully applied")
	}

	if a.config.Database.LoadFixtures {
		_ = a.logger.Log("msg", "start loading fixtures")
		err := database.LoadFixtures(context.Background(), a.dbConn)
		if err != nil {
			return err
		}
		_ = a.logger.Log("msg", "db fixture successfully loaded")
	}
	return nil
}

// Unwire closes the database connection.
func (a *App) Unwire() {
	a.dbConn.Close()
}

// Run starts the application.
func (a *App) Run() {
	server := a.startServer(a.config, a.router, a.logger)
	a.stopServerOnSignal(server, a.logger)
}

func (a *App) startServer(cfg config.Application, router http.Handler, logger log.Logger) *http.Server {
	logger = log.With(logger, "label", "server")

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	_ = level.Info(logger).Log("msg", fmt.Sprintf("starting %s µ service", cfg.ID), "address", cfg.Address)

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			_ = level.Error(logger).Log("msg", "an error occurred after stopping the server", "err", err)
		}
	}()

	return server
}

func (a *App) stopServerOnSignal(server *http.Server, logger log.Logger) {
	logger = log.With(logger, "label", "server")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	sig := <-sigCh

	_ = level.Info(logger).Log("msg", "going to shutdown gracefully due to received signal", "signal", sig)

	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	err := server.Shutdown(ctx)
	cancelFn()

	if err != nil {
		_ = level.Error(logger).Log("msg", "failed to shut down gracefully", "err", err)
	}
}

func healthHandler(message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonMessage := struct {
			Message string `json:"message"`
		}{Message: message}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(jsonMessage)
	})
}

func notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(fmt.Errorf("resource %s does not exist", r.URL.Path))
	})
}
