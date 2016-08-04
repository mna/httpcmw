package requestid

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/httpcmw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockStatusDoer(code int) httpcmw.Doer {
	return httpcmw.DoerFunc(func(req *http.Request) (*http.Response, error) {
		res := &http.Response{Request: req, StatusCode: code}
		return res, nil
	})
}

func TestRequestID(t *testing.T) {
	cases := []struct {
		rid    *RequestID
		preset string
	}{
		{&RequestID{}, ""},
		{&RequestID{Header: "XYZ"}, ""},
		{&RequestID{Header: "XYZ"}, "abc"},
		{&RequestID{Header: "XYZ", Len: 12}, ""},
		{&RequestID{Header: "XYZ", Len: 12}, "abc"},
		{&RequestID{Header: "XYZ", Len: 12, ForceSet: true}, ""},
		{&RequestID{Header: "XYZ", Len: 12, ForceSet: true}, "abc"},
	}
	for i, c := range cases {
		key := "X-Request-Id"
		if c.rid.Header != "" {
			key = c.rid.Header
		}

		d := httpcmw.Wrap(mockStatusDoer(204), c.rid)
		r, _ := http.NewRequest("", "/", nil)
		if c.preset != "" {
			r.Header.Set(key, c.preset)
		}
		res, err := d.Do(r)
		require.NoError(t, err, "%d: Do error", i)

		assert.Equal(t, 204, res.StatusCode, "%d: status", i)
		got := r.Header.Get(key)
		t.Logf("%d: got request ID %q", i, got)

		if c.preset != "" && !c.rid.ForceSet {
			assert.Equal(t, c.preset, got, "%d: id", i)
			continue
		}
		wantLen := c.rid.Len
		if wantLen == 0 {
			wantLen = 8
		}
		assert.Equal(t, wantLen, len(got), "%d: length", i)
		assert.NotEqual(t, c.preset, got, "%d: not preset value", i)
	}
}
