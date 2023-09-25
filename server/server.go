package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jaro825/api-server-demo/api"
	"github.com/jaro825/api-server-demo/cache"
	"github.com/jaro825/api-server-demo/store"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	logger     zerolog.Logger
}

func New(logger zerolog.Logger, store store.Store, cache cache.Cache, cfg *Config) (*Server, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("server configuration error: %w", err)
	}

	srv := &Server{logger: logger}

	mux := chi.NewRouter()

	// middleware
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(ZeroLogger(logger))
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))

	userAPI := api.NewUserAPI(logger, store, cache)

	// routes
	mux.Route("/api/v1", func(r chi.Router) {
		r.Mount("/user", userRouter(userAPI))
	})

	srv.httpServer = &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       time.Minute,
	}

	return srv, nil
}

func ZeroLogger(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
			start := time.Now()

			defer func() {
				switch ww.Status() {
				case http.StatusOK, http.StatusCreated:
					logger.Info().
						Int("status", ww.Status()).
						Int("bytes", ww.BytesWritten()).
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Str("query", r.URL.RawQuery).
						Str("ip", r.RemoteAddr).
						Dur("latency", time.Since(start)).
						Msg("request completed")
				default:
					logger.Warn().
						Int("status", ww.Status()).
						Int("bytes", ww.BytesWritten()).
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Str("query", r.URL.RawQuery).
						Str("ip", r.RemoteAddr).
						Dur("latency", time.Since(start)).
						Msg("request completed")
				}
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}

func userRouter(api *api.UserAPI) chi.Router { //nolint: ireturn
	r := chi.NewRouter()

	r.Post("/", api.CreateUser) // POST /api/v1/user

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", api.GetUser)       // GET /api/v1//user/c444f05f-67d6-491a-93be-f2ec09503b71
		r.Put("/", api.UpdateUser)    // PUT /api/v1/user/c444f05f-67d6-491a-93be-f2ec09503b71
		r.Delete("/", api.DeleteUser) // DELETE /api/v1/user/c444f05f-67d6-491a-93be-f2ec09503b71
	})

	return r
}

func (s *Server) Run() error {
	if err := s.httpServer.ListenAndServe(); err != nil {
		if errors.Is(http.ErrServerClosed, err) {
			s.logger.Debug().Msg("http server closed")

			return nil
		}

		s.logger.Err(err).Msg("unexpected error in http server")

		return fmt.Errorf("ListenAndServe error: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
