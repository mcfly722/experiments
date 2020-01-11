package apiserver

import (
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mcfly722/experiments/go-rest/internal/app/store"
	"github.com/sirupsen/logrus"
)

// APIServer ...
type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

// New ...
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start ...
func (s *APIServer) Start() error {

	if err := s.configureLogger(); err != nil {
		return err
	}

	s.logger.Info("starting API server")

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureRouter() {
	s.router.Use(s.logRequest)
	s.router.HandleFunc("/hello", s.handleHello())
}

func (s *APIServer) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Infof("started: %s %s", r.Method, r.RequestURI)
		start := time.Now()
		next.ServeHTTP(w, r)
		s.logger.Infof("completed in: %vms", time.Now().Sub(start).Milliseconds())
	})
}

func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}
