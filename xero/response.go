package xero

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

var (
	// Error returned for 400 status code.
	ErrInvalidRequest = errors.New("invalid parameter")

	// See https://developer.xero.com/documentation/guides/oauth2/limits/#api-rate-limits
	ErrTooManyRequests = errors.New("request hit the rate limit")

	// Error returned for any 5xx status code.
	ErrXeroDown = errors.New("xero API is not reachable")

	// XeroDateFormat matches .NET JSON date format in a string.
	XeroDateFormat = regexp.MustCompile(`Date\((?P<Value>\d+)\)`)
)

// Response contains fields common to all Xero API's responses.
type Response struct {
	Status string `description:"Actual HTTP response status" json:"status"`
}

// OK returns whether the response was correctly returned by Xero API.
func (r Response) OK() bool {
	return r.Status == "OK"
}

// DateTimeField is a type that implements json.Unmarshaler for handling Microsoft .NET JSON date format as utilised by Xero API.
// See https://developer.xero.com/documentation/api/accounting/requests-and-responses#json-responses-and-date-formats
type DateTimeField struct {
	time.Time
}

// UnmarshalJSON satisfies json.Unmarshaler interface.
func (dt *DateTimeField) UnmarshalJSON(data []byte) error {
	xeroDate := string(data)
	if len(xeroDate) == 0 {
		return nil
	}

	matches := XeroDateFormat.FindStringSubmatch(xeroDate)
	if len(matches) != 2 {
		return nil
	}

	unixTS, err := strconv.ParseInt(matches[1], 10, 64)

	switch {
	case err != nil:
		return errors.Join(errors.New("could not unmarshal .NET timestamp"), err)
	case unixTS <= 0:
		return errors.New("invalid zero timestamp")
	default:
		*dt = DateTimeField{time.Unix(unixTS/1000, 0)}

		return nil
	}
}
