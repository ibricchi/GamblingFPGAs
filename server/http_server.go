package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type HttpServer struct {
	// db     *SQLiteDB
	router *chi.Mux
	logger *zap.Logger
}

func OpenHttpServer(ctx context.Context, logger *zap.Logger, router *chi.Mux) (*HttpServer, error) {
	h := &HttpServer{
		router: router,
		logger: logger,
	}

	return h, nil
}

func (h *HttpServer) Serve(port string) error {
	h.routes()

	portStr := ":" + port
	if err := http.ListenAndServe(portStr, h.router); err != nil {
		return fmt.Errorf("server: http_server: http.ListenAndServe failed: %w", err)
	}
	return nil
}

func (h *HttpServer) Close() error {
	// used for closing db, etc.
	return nil
}