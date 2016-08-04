package logrequest

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/PuerkitoBio/httpcmw"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ok(req *http.Request) (*http.Response, error) {
	return &http.Response{Request: req, StatusCode: 204}, nil
}

func TestLogRequest(t *testing.T) {
	var buf bytes.Buffer
	l := log.NewLogfmtLogger(&buf)
	lr := &LogRequest{Logger: l, DurationFormat: "%.5f", Fields: []string{"duration", "method", "path"}}
	d := httpcmw.Wrap(httpcmw.DoerFunc(ok), lr)

	r, _ := http.NewRequest("PUT", "/", nil)
	res, err := d.Do(r)
	require.NoError(t, err)

	assert.Equal(t, 204, res.StatusCode, "status")
	assert.Contains(t, buf.String(), " method=PUT path=/", "expected output")
}
