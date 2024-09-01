package web_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/luca-arch/code-drills/web"
	"github.com/luca-arch/code-drills/xero"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	err error
	res *xero.ReportResponse
}

func (m *mockClient) BalanceSheet(context.Context) (*xero.ReportResponse, error) {
	return m.res, m.err
}

func TestBalance(t *testing.T) {
	t.Parallel()

	nopLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	time.Local = time.UTC

	type fields struct {
		mockClient *mockClient
	}

	type wants struct {
		body   string
		status int
	}

	tests := map[string]struct {
		fields
		wants
	}{
		"success": {
			fields{
				mockClient: &mockClient{
					res: xeroStubReports(t),
				},
			},
			wants{
				body:   fixture(t, "testdata/get-balance.json"),
				status: http.StatusOK,
			},
		},
		"error - invalid GET parameters": {
			fields{
				mockClient: &mockClient{
					err: xero.ErrInvalidRequest,
				},
			},
			wants{
				body:   "invalid parameter\n",
				status: http.StatusBadRequest,
			},
		},
		"error - rate limit exceeded": {
			fields{
				mockClient: &mockClient{
					err: xero.ErrTooManyRequests,
				},
			},
			wants{
				body:   "enhance your calm!\n",
				status: http.StatusTooManyRequests,
			},
		},
		"error - Xero API not reachable": {
			fields{
				mockClient: &mockClient{
					err: xero.ErrXeroDown,
				},
			},
			wants{
				body:   "Xero API not available at the moment\n",
				status: http.StatusGatewayTimeout,
			},
		},
		"error - other errors": {
			fields{
				mockClient: &mockClient{
					err: errors.New("something wrong"),
				},
			},
			wants{
				body:   "something wrong\n",
				status: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			server := web.HTTPServer(nopLogger, test.fields.mockClient)
			testServer := httptest.NewServer(server.Mux())

			t.Cleanup(testServer.Close)

			//nolint:noctx // Ok when testing
			res, err := http.Get(testServer.URL + "/balance")
			assert.NoError(t, err)

			body, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			res.Body.Close()

			bodyStr := string(body)

			assert.Equal(t, test.wants.status, res.StatusCode)
			assert.Equal(t, test.wants.body, bodyStr, "Actual: "+bodyStr)
		})
	}
}

func fixture(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}

func xeroDTField(t *testing.T, unixSeconds int) xero.DateTimeField {
	t.Helper()

	return xero.DateTimeField{ //nolint:govet
		time.Unix(int64(unixSeconds), 0),
	}
}

func xeroStubReports(t *testing.T) *xero.ReportResponse {
	t.Helper()

	stubReport := xero.Report{
		ReportID:   "1234",
		ReportName: "Test Sheet",
		ReportType: "BalanceSheet",
		ReportTitles: []string{
			"Title 01",
			"Title 02",
		},
		ReportDate:     "25 August 2024",
		UpdatedDateUTC: xeroDTField(t, 1724595191),
		Rows: []xero.Row{
			{
				RowType: "Header",
				Cells: []xero.Cell{
					{
						Value: "",
					},
					{
						Value: "25 August 2024",
					},
					{
						Value: "26 August 2023",
					},
				},
			},
			{
				RowType: "Section",
				Title:   "Assets",
			},
			{
				RowType: "Section",
				Title:   "Bank",
				Rows: []xero.Row{
					{
						RowType: "Row",
						Cells: []xero.Cell{
							{
								Value: "My Bank Account",
								Attributes: []xero.Attributes{
									{
										ID:    "account-id",
										Value: "some value",
									},
								},
							},
							{
								Value: "126.70",
								Attributes: []xero.Attributes{
									{
										ID:    "account-id",
										Value: "other value",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return &xero.ReportResponse{
		Reports: []xero.Report{
			stubReport,
		},
	}
}
