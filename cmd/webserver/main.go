package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/luca-arch/code-drills/web"
	"github.com/luca-arch/code-drills/xero"
)

func debugLogger() *slog.Logger {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})

	return slog.New(handler)
}

func main() {
	logger := debugLogger()

	apiClient := xero.HTTPClient(logger).
		WithBaseURL("http://mock-xero:3000")

	server := web.HTTPServer(logger, apiClient)

	//nolint:gosec // "G114: Use of net/http serve function that has no support for setting timeouts" can be ignored for this demo
	err := http.ListenAndServe(":4000", server.Mux())
	if err != nil {
		logger.Error(err.Error())

		os.Exit(1)
	}
}
