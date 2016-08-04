// Package logrequest implements a middleware that logs requests.
package logrequest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/httpcmw"
)

var allFields = []string{
	"body_bytes_received",
	"duration",
	"end",
	"error",
	"host",
	"method",
	"path",
	"proto",
	"query",
	"request_id",
	"start",
	"status",
	"user_agent",
}

// LogRequest holds the configuration for the LogRequest middleware.
type LogRequest struct {
	// Logger is the logger to use to log the requests.
	Logger httpcmw.Logger

	// RequestIDHeader is the name of the header that contains the request
	// ID. Defaults to X-Request-Id.
	RequestIDHeader string

	// TimeFormat is the format to use to format timestamps, as supported by
	// the time package.
	TimeFormat string

	// DurationFormat is the format to use to log the duration of the request.
	// The value to format is the number of seconds in float64. Defaults to
	// %.3f (milliseconds).
	DurationFormat string

	// Fields is the list of field names to log. Defaults to all supported
	// fields. The supported fields are:
	//
	//     body_bytes_received: bytes received in the response body
	//     duration: duration of the request
	//     end: date and time of the end of the request (UTC)
	//     error: error returned from the request, if any
	//     host: host (and possibly port) of the request
	//     method: method of the request (e.g. GET)
	//     path: path section of the request URL
	//     proto: protocol and version (e.g. HTTP/1.1)
	//     query: query string section of the request URL
	//     request_id: request id
	//     start: date and time of the start of the request (UTC)
	//     status: status code of the response, if any
	//     user_agent: value of the User-Agent request header
	//
	Fields []string
}

// Wrap returns a Doer that records the start time, calls the Doer d,
// records the end time and duration, and logs the request's fields as
// configured by the LogRequest.
func (lr *LogRequest) Wrap(d httpcmw.Doer) httpcmw.Doer {
	log := lr.Logger
	if log == nil {
		return d
	}
	dfmt := lr.DurationFormat
	if dfmt == "" {
		dfmt = "%.3f"
	}
	tf := lr.TimeFormat
	if tf == "" {
		tf = time.ANSIC
	}
	hd := lr.RequestIDHeader
	if hd == "" {
		hd = "X-Request-Id"
	}
	fields := lr.Fields
	if len(fields) == 0 {
		fields = allFields
	}

	return httpcmw.DoerFunc(func(r *http.Request) (*http.Response, error) {
		start := time.Now().UTC()
		res, err := d.Do(r)
		end := time.Now().UTC()

		rid := r.Header.Get(hd)
		vals := map[string]string{
			"start":      start.Format(tf),
			"end":        end.Format(tf),
			"duration":   fmt.Sprintf(dfmt, end.Sub(start).Seconds()),
			"host":       r.Host,
			"method":     r.Method,
			"request_id": rid,
			"path":       r.URL.Path,
			"user_agent": r.UserAgent(),
			"query":      r.URL.RawQuery,
		}
		if res != nil {
			vals["proto"] = res.Proto
			vals["status"] = strconv.Itoa(res.StatusCode)
			vals["body_bytes_received"] = strconv.FormatInt(res.ContentLength, 10)
			if rid == "" {
				vals["request_id"] = res.Header.Get(hd)
			}
		}
		if err != nil {
			vals["error"] = err.Error()
		}

		args := make([]interface{}, 0, len(fields)*2)
		for _, f := range fields {
			args = append(args, f, vals[f])
		}
		log.Log(args...)

		return res, err
	})
}
