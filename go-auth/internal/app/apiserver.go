package apiserver

import (
	"io"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	s.logger.Info("starting API server at ", s.config.BindAddr)

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
	s.router.Handle("/", http.FileServer(http.Dir("./frontend/")))

	var authorizationHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "hello" && password == "world" {
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["userId"] = 1234
			claims["timeStamp"] = time.Now().Unix()
			if s.config.JWTPrivateKey == "" {
				s.logger.Error("JWT Token error: token is not specified")
				http.Error(w, "Unauthorized. See server log for detailed error.", http.StatusUnauthorized)
			} else {
				tokenString, err := token.SignedString([]byte(s.config.JWTPrivateKey))
				if err != nil {
					s.logger.Error("JWT Token error:", err)
					http.Error(w, "Unauthorized. See server log for detailed error.", http.StatusUnauthorized)
				} else {
					w.Write([]byte(tokenString))
				}
			}
		} else {
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		}
	})

	s.router.Handle("/authorization", authorizationHandler).Methods("GET")

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
