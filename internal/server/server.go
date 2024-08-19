package server

import (
	"context"
	"net/http"

	"github.com/Ser9unin/RealEstate/internal/auth"
	"github.com/Ser9unin/RealEstate/internal/config"
	repository "github.com/Ser9unin/RealEstate/internal/storage/repo"
)

type Server struct {
	srv     *http.Server
	router  *http.ServeMux
	storage *repository.Queries
	logger  Logger
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

func NewServer(cfg config.SrvCfg, logger Logger, storage *repository.Queries) *Server {
	router := NewRouter(storage, logger)

	srv := &http.Server{
		Addr:              cfg.Host + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout, // Настраиваем тайм-аут ожидания заголовков
		ReadTimeout:       cfg.ReadTimeout,       // Настраиваем общий тайм-аут запроса
		WriteTimeout:      cfg.WriteTimeout,      // Настраиваем тайм-аут записи ответа
		IdleTimeout:       cfg.IdleTimeout,       // Настраиваем тайм-аут простоя соединения
	}

	return &Server{srv, router, storage, logger}
}

func NewRouter(storage *repository.Queries, logger Logger) *http.ServeMux {
	mux := http.NewServeMux()

	mw := func(next http.HandlerFunc) http.HandlerFunc {
		return HTTPLogger(CheckHTTPMethod(next))
	}

	a := newAPI(storage, logger)

	mux.HandleFunc("/", HTTPLogger(a.greetings))
	// mux.HandleFunc("/dummyLogin", HTTPLogger(a.dummyLogin))
	mux.HandleFunc("/register", mw(a.register))
	mux.HandleFunc("/login", mw(a.login))
	mux.HandleFunc("/house/create", auth.Moderator(mw(a.houseCreate)))
	mux.HandleFunc("/house/{id}", auth.Any(HTTPLogger(a.houseFlats)))
	mux.HandleFunc("/house/{id}/subscribe", auth.Any(mw(a.houseSubscribe)))
	mux.HandleFunc("/flat/create", auth.Any(mw(a.flatCreate)))
	mux.HandleFunc("/flat/update", auth.Moderator(mw(a.flatUpdate)))

	return mux
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
