package xero_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/luca-arch/code-drills/xero"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	timePtr := func(timeStr string) *time.Time {
		ptr, err := time.Parse(time.RFC3339, timeStr)
		assert.NoError(t, err)

		return &ptr
	}

	tests := map[string]struct {
		arg    string
		err    error
		result *time.Time
	}{
		"empty string": {
			arg:    "",
			err:    nil,
			result: nil,
		},
		"non-date string": {
			arg:    "not a date",
			err:    nil,
			result: nil,
		},
		".NET datetime": {
			arg:    `\/Date(1724536800000)\/`,
			err:    nil,
			result: timePtr("2024-08-24T22:00:00+00:00"),
		},
		"Zero datetime": {
			arg:    `\/Date(0)\/`,
			err:    errors.New("invalid zero timestamp"),
			result: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dt := new(xero.DateTimeField)

			err := dt.UnmarshalJSON([]byte(test.arg))

			if test.err != nil {
				assert.EqualError(t, test.err, err.Error())
				assert.Empty(t, dt)
			} else {
				assert.NoError(t, err)
			}

			if test.result != nil {
				assert.EqualValues(t, test.result.UTC(), dt.UTC())
			}
		})
	}
}

func TestResponseOK(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		arg string
		ok  bool
	}{
		"response OK": {
			arg: `{"status": "OK"}`,
			ok:  true,
		},
		"wrong case": {
			arg: `{"status": "ok"}`,
			ok:  false,
		},
		"response empty": {
			arg: `{}`,
			ok:  false,
		},
		"response NOT ok": {
			arg: `{"status": "error"}`,
			ok:  false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var response xero.Response

			err := json.Unmarshal([]byte(test.arg), &response)
			assert.NoError(t, err)

			assert.Equal(t, test.ok, response.OK())
		})
	}
}
