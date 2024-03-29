package server

import (
	"context"
	"fmt"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (h *HttpServer) routes(ctx context.Context) error {
	// General middleware
	h.router.Use(cors.Handler(cors.Options{
		AllowOriginFunc:  allowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))
	h.router.Use(middleware.Logger)

	// Credentials from database
	creds, err := h.db.getCreds(ctx)
	if err != nil {
		return fmt.Errorf("server: routes: failed to get credentials: %w", err)
	}

	// Public routes
	h.router.Group(func(r chi.Router) {
		// Get
		r.Get("/public/test", h.handleGetStaticTest())
		r.Get("/isAuthorised", h.handleGetIsAuthorised(creds))

		// Post
		r.Post("/public/test/dynamic", h.handlePostDynamicTest(ctx))
	})

	// Private routes
	h.router.Group(func(r chi.Router) {
		r.Use(h.basicAuth("GamblingFPGAs-API", creds))

		// Get
		r.Get("/test", h.handleGetStaticTest())
		r.Get("/poker/openGameStatus", h.handlePokerGetGameOpenStatus())
		r.Get("/poker/activeGameStatus", h.handlePokerGetGameActiveStatus())
		r.Get("/poker/activeGameStatus/showdown", h.handlePokerGetGameShowdownData())
		r.Get("/poker/fpgaData", h.handlePokerGetFPGAData())

		// Post
		r.Post("/test/dynamic", h.handlePostDynamicTest(ctx))
		r.Post("/poker/openGame", h.handlePokerOpenGame())
		r.Post("/poker/joinGame", h.handlePokerJoinGame())
		r.Post("/poker/startGame", h.handlePokerStartGame())
		r.Post("/poker/startNewGameSamePlayers", h.handlePokerStarNewGameSamePlayers())
		r.Post("/poker/terminateGame", h.handlePokerTerminateGame())
		r.Post("/poker/fpgaData", h.handlePokerReceiveFPGAData())
	})

	// Not found
	h.router.NotFound(h.handleNotFound())

	return nil
}
