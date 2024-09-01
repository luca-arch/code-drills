// Package web provides a server mux for serving HTTP requests.
package web

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/luca-arch/code-drills/xero"
)

// xeroclient defines an interface to make Xero API requests.
type xeroclient interface {
	BalanceSheet(context.Context) (*xero.ReportResponse, error)
}

// server defines a concrete type to serve HTTP requests.
type server struct {
	client xeroclient
	logger *slog.Logger
}

// HTTPServer returns a new HTTP server with default configuration.
func HTTPServer(logger *slog.Logger, apiClient xeroclient) *server {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	logger.Debug("initialising new HTTP server")

	return &server{
		client: apiClient,
		logger: logger,
	}
}

// Mux returns a new server mux with the following routes:
// - GET /balance.
func (s *server) Mux() http.Handler {
	mux := &http.ServeMux{}

	mux.Handle("GET /balance", s.listBalanceSheetHandler())

	return mux
}

// listBalanceSheetHandler returns an HTTP handler that serves the GET "/balance" endpoint.
func (s *server) listBalanceSheetHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debug("incoming HTTP request", "client", r.Header.Get("User-Agent"))

		rr, err := s.client.BalanceSheet(r.Context())

		switch {
		case err == nil:
			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(rr); err != nil {
				s.logger.Warn("Could not marshal into response", "err", err)
			}
		case errors.Is(err, xero.ErrInvalidRequest):
			// Fixme: this assumes an invalid GET parameter was passed from the client.
			http.Error(w, "invalid parameter", http.StatusBadRequest)
		case errors.Is(err, xero.ErrTooManyRequests):
			http.Error(w, "enhance your calm!", http.StatusTooManyRequests)
		case errors.Is(err, xero.ErrXeroDown):
			http.Error(w, "Xero API not available at the moment", http.StatusGatewayTimeout)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
