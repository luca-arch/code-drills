package xero

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

const DefaultBaseURL = "https://api.xero.com" // Default Xero API domain.

var (
	ErrBrokenResponse  = errors.New("xero response with error")                       // Xero response status (in the body) not OK .
	ErrHTTPFailure     = errors.New("internal error")                                 // HTTP transport error.
	ErrInvalidJSON     = errors.New("error while unmarshalling Xero reports")         // Xero API returned malformed JSON.
	ErrInvalidResponse = errors.New("error while unmarshalling Xero response")        // Xero API returned non-JSON response.
	ErrRequestFailure  = errors.New("error while retrieving Balance Sheet from Xero") // HTTP request error.
)

// HTTPDoer defines an interface to make HTTP requests.
type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

// client defines a concrete type to invoke the Xero API.
type client struct {
	base   string
	client HTTPDoer
	logger *slog.Logger
}

// HTTPClient returns a new Xero HTTP client with default configuration.
func HTTPClient(logger *slog.Logger) *client {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	logger.Debug("initialising new Xero HTTP client")

	return &client{
		base:   DefaultBaseURL,
		client: http.DefaultClient,
		logger: logger,
	}
}

// BalanceSheet invokes the Reports BalanceSheet endpoint and returns a list of reports.
// See https://developer.xero.com/documentation/api/accounting/reports#balance-sheet
func (c *client) BalanceSheet(ctx context.Context) (*ReportResponse, error) {
	var (
		r  Response
		rr ReportResponse
	)

	c.logger.Debug("Outgoing HTTP request", "endpoint", " /api.xro/2.0/Reports/BalanceSheet")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+"/api.xro/2.0/Reports/BalanceSheet", nil)
	if err != nil {
		return nil, errors.Join(ErrHTTPFailure, err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Join(ErrRequestFailure, err)
	}

	c.logger.Debug("HTTP request finished", "status", resp.StatusCode)

	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusOK:
		break
	case resp.StatusCode == http.StatusBadRequest:
		return nil, ErrInvalidRequest
	case resp.StatusCode == http.StatusTooManyRequests:
		return nil, ErrTooManyRequests
	case resp.StatusCode >= http.StatusInternalServerError:
		return nil, ErrXeroDown
	default:
		return nil, errors.New("invalid status in Xero response: " + strconv.Itoa(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err == nil {
		err = json.Unmarshal(body, &r)
	}

	switch {
	case err != nil:
		return nil, errors.Join(ErrInvalidResponse, err)
	case !r.OK():
		return nil, ErrBrokenResponse
	}

	if err = json.Unmarshal(body, &rr); err != nil {
		return nil, errors.Join(ErrInvalidJSON, err)
	}

	return &rr, nil
}

// WithBaseURL sets the client's base URL.
func (c *client) WithBaseURL(base string) *client {
	c.base = base

	return c
}

// WithHTTPClient sets the client's HTTP doer.
func (c *client) WithHTTPClient(client HTTPDoer) *client {
	c.client = client

	return c
}
