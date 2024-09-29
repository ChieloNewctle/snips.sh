package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/robherley/snips.sh/internal/config"
	"github.com/robherley/snips.sh/internal/db"
)

type Service struct {
	*http.Server
	Router *chi.Mux
}

func New(cfg *config.Config, database db.DB, assets Assets) (*Service, error) {
	router := chi.NewRouter()

	router.Use(WithRequestID)
	router.Use(WithLogger)
	router.Use(WithMetrics)
	router.Use(WithRecover)

	router.Get("/", DocHandler(cfg, assets))
	router.Get("/docs/{name}", DocHandler(cfg, assets))
	router.Get("/health", HealthHandler)
	router.Get("/f/{fileID}", FileHandler(cfg, database, assets))
	router.Get("/assets/index.js", assets.ServeJS)
	router.Get("/assets/index.css", assets.ServeCSS)
	router.Get("/meta.json", MetaHandler(cfg))

	if cfg.Debug {
		router.Mount("/_debug", middleware.Profiler())
	}

	handler := chi.NewRouter()
	handler.Mount(fmt.Sprintf("/{}", cfg.HTTP.BasePath), router)

	httpServer := &http.Server{
		Addr:    cfg.HTTP.Internal.Host,
		Handler: handler,
	}

	return &Service{httpServer, router}, nil
}
