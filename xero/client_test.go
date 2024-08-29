package xero_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

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

func mockReadCloser() io.ReadCloser {
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
						Body:       mockReadCloser(),
						StatusCode: test.mockStatus,
					},
				})

			resp, err := client.BalanceSheet(context.TODO())

			assert.Nil(t, resp)
			assert.EqualError(t, err, test.expectedError.Error())
		})
	}
}

// TODO: extend this test.
func TestBalanceSheetResponse(t *testing.T) {
	t.Parallel()

	mock, err := os.ReadFile("testdata/001.json")
	assert.NoError(t, err)

	body := string(mock)

	client := xero.HTTPClient(nil).
		WithHTTPClient(&mockHTTPDoer{
			mockError: nil,
			mockResponse: &http.Response{
				Body:       io.NopCloser(strings.NewReader(body)),
				StatusCode: http.StatusOK,
			},
		})

	resp, err := client.BalanceSheet(context.TODO())

	assert.NoError(t, err)

	assert.Len(t, resp.Reports, 1)

	report := resp.Reports[0]

	assert.Equal(t, "1234", report.ReportID)
	assert.Equal(t, "Test Sheet", report.ReportName)
	assert.Equal(t, "BalanceSheet", report.ReportType)
	assert.Equal(t, []string{"Title 01", "Title 02"}, report.ReportTitles)
	assert.Equal(t, "25 August 2024", report.ReportDate)
	assert.Equal(t, int64(1724595191), report.UpdatedDateUTC.Unix())

	assert.Len(t, report.Rows, 3)
}
