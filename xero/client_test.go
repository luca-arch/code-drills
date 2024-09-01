package xero_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/luca-arch/code-drills/xero"
	"github.com/stretchr/testify/assert"
)

type mockHTTPDoer struct {
	mockError    error
	mockResponse *http.Response
}

func (m mockHTTPDoer) Do(_ *http.Request) (*http.Response, error) {
	return m.mockResponse, m.mockError
}

func mockReadCloser(t *testing.T) io.ReadCloser {
	t.Helper()

	fake := strings.NewReader("")

	return io.NopCloser(fake)
}

func TestBalanceSheetError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		mockError     error
		mockStatus    int
		expectedError error
	}{
		"HTTP 400": {
			mockError:     nil,
			mockStatus:    http.StatusBadRequest,
			expectedError: xero.ErrInvalidRequest,
		},
		"HTTP 404": {
			mockError:     nil,
			mockStatus:    http.StatusNotFound,
			expectedError: errors.New("invalid status in Xero response: 404"),
		},
		"HTTP 429": {
			mockError:     nil,
			mockStatus:    http.StatusTooManyRequests,
			expectedError: xero.ErrTooManyRequests,
		},
		"HTTP 500": {
			mockError:     nil,
			mockStatus:    http.StatusInternalServerError,
			expectedError: xero.ErrXeroDown,
		},
		"HTTP 501": {
			mockError:     nil,
			mockStatus:    http.StatusNotImplemented,
			expectedError: xero.ErrXeroDown,
		},
		"HTTP 502": {
			mockError:     nil,
			mockStatus:    http.StatusBadGateway,
			expectedError: xero.ErrXeroDown,
		},
		"HTTP 503": {
			mockError:     nil,
			mockStatus:    http.StatusServiceUnavailable,
			expectedError: xero.ErrXeroDown,
		},
		"HTTP 504": {
			mockError:     nil,
			mockStatus:    http.StatusGatewayTimeout,
			expectedError: xero.ErrXeroDown,
		},
		"Network error": {
			mockError:     errors.New("mock TCP err"),
			mockStatus:    0,
			expectedError: errors.New("error while retrieving Balance Sheet from Xero\nmock TCP err"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := xero.HTTPClient(nil).
				WithHTTPClient(&mockHTTPDoer{
					mockError: test.mockError,
					mockResponse: &http.Response{
						Body:       mockReadCloser(t),
						StatusCode: test.mockStatus,
					},
				})

			resp, err := client.BalanceSheet(context.TODO())

			assert.Nil(t, resp)
			assert.EqualError(t, err, test.expectedError.Error())
		})
	}
}

func TestBalanceSheetResponse(t *testing.T) {
	t.Parallel()

	type fields struct {
		body   []byte
		status int
	}

	type wants struct {
		err    error
		report xero.Report
	}

	tests := map[string]struct {
		fields
		wants
	}{
		"success": {
			fields{
				body:   fixture(t, "testdata/reports.json"),
				status: http.StatusOK,
			},
			wants{
				report: xero.Report{
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
				},
			},
		},
		"error - response.Body.Status not OK": {
			fields{
				body:   fixture(t, "testdata/error.json"),
				status: http.StatusOK,
			},
			wants{
				err: xero.ErrBrokenResponse,
			},
		},
		"error - response.Body.Reports invalid": {
			fields{
				body:   fixture(t, "testdata/invalid.json"),
				status: http.StatusOK,
			},
			wants{
				err: xero.ErrInvalidJSON,
			},
		},
		"error - invalid content type": {
			fields{
				body:   []byte("hello"),
				status: http.StatusOK,
			},
			wants{
				err: xero.ErrInvalidResponse,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := xero.HTTPClient(nil).
				WithHTTPClient(&mockHTTPDoer{
					mockError: nil,
					mockResponse: &http.Response{
						Body:       io.NopCloser(bytes.NewReader(test.fields.body)),
						StatusCode: test.status,
					},
				})

			resp, err := client.BalanceSheet(context.TODO())

			if test.wants.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, resp.Reports[0], test.wants.report)

				return
			}

			assert.ErrorIs(t, err, test.wants.err)
		})
	}
}

func fixture(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return data
}

func xeroDTField(t *testing.T, unixSeconds int) xero.DateTimeField {
	t.Helper()

	return xero.DateTimeField{ 
		time.Unix(int64(unixSeconds), 0),
	}
}
