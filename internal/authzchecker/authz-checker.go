// Package authzchecker contains the main logic of the authz-checker service
package authzchecker

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Stoakes/authz-checker/internal/config"

	"github.com/Stoakes/go-pkg/log"
	"github.com/gorilla/mux"
)

// AuthzChecker runs authz server
type AuthzChecker interface {
	Start(ctx context.Context) error
}

// authzChecker implements AuthzChecker interface
type authzChecker struct {
	port int
	// db        database.DB
	appConfig config.AppConfig
	router    *mux.Router
}

// Parameters gathers every parameter required to create an authz-checker
type Parameters struct {
	Port int
	// DB        database.DB
	AppConfig config.AppConfig
}

// New creates a new authz-checker
func New(params Parameters) (AuthzChecker, error) {
	return newAuthzChecker(params)
}

// newAuthzChecker enables tests to have access to concrete implementation of AuthzChecker interface
func newAuthzChecker(params Parameters) (*authzChecker, error) {
	if params.AppConfig.Port == params.AppConfig.ToolsPort {
		return nil, fmt.Errorf("app server port and tools server port must be different")
	}
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello world")
	})
	a := &authzChecker{
		port: params.Port,
		// db:        params.DB,
		appConfig: params.AppConfig,
		router:    r,
	}
	return a, nil
}

// Start starts the HTTP server
func (s *authzChecker) Start(ctx context.Context) error {
	log.For(ctx).Infof("Starting authz-checker HTTP server on  port %d", s.port)
	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(s.port),
		Handler:           s.router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return srv.ListenAndServe()
}
