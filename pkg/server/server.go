package server

import (
	"context"
	"golang-jwt-example/pkg/handler"
	"golang-jwt-example/pkg/middleware"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Config struct {
	Log *zap.Logger
}

type Server struct {
	Mux        *http.ServeMux
	server     *http.Server
	handler    *handler.Handler
	middleware *middleware.Middleware
	log        *zap.Logger
}

func NewServer(handler *handler.Handler, middleware *middleware.Middleware, cfg *Config) *Server {
	s := &Server{
		Mux:        http.NewServeMux(),
		handler:    handler,
		middleware: middleware,
	}

	if cfg != nil {
		if log := cfg.Log; log != nil {
			s.log = log
		}
	}

	s.registerHandler()
	return s
}

func (s *Server) Serve(listener net.Listener) error {
	server := &http.Server{
		Handler: s.Mux,
	}
	s.server = server
	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) GracefulShutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) registerHandler() {
	r := mux.NewRouter()

	st := r.PathPrefix("/user").Subrouter()
	st.Handle("/all", s.handler.User.ListHandler()).Methods(http.MethodGet)
	st.Handle("/login", s.handler.User.LoginHandler()).Methods(http.MethodPost)
	// st.Handle("/token", s.handler.User.GetToken()).Methods(http.MethodPost)

	// r.Handle("/version", s.handler.General.Version())
	// r.NotFoundHandler = s.handler.General.NotFound()

	st.Handle("/auth", s.middleware.Auth(s.handler.User.GetAuthUser())).Methods(http.MethodGet)
	s.Mux.Handle("/", r)
}
