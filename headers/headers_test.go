package headers

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/httpcmw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ok(req *http.Request) (*http.Response, error) {
	return &http.Response{Request: req, StatusCode: 204}, nil
}

func TestHeadersPrefix(t *testing.T) {
	head := make(Headers)
	head.Add("A", "a")
	head.Add("A", "aa")
	head.Add("+B", "b")
	head.Add("-C", "")
	head.Add("+D", "d")

	r, _ := http.NewRequest("", "/", nil)
	d := head.Wrap(httpcmw.DoerFunc(ok))
	r.Header.Set("C", "c")
	r.Header.Set("A", "z")
	r.Header.Set("B", "x")

	res, err := d.Do(r)
	require.NoError(t, err)
	assert.Equal(t, 204, res.StatusCode, "status")
	assert.Equal(t, http.Header{
		"A": {"a", "aa"},
		"B": {"x", "b"},
		"D": {"d"},
	}, r.Header, "headers")
}

func TestHeaders(t *testing.T) {
	head := make(Headers)
	head.Add("A", "a")
	head.Add("B", "b")
	head.Add("C", "c")

	d := httpcmw.Wrap(httpcmw.DoerFunc(func(r *http.Request) (*http.Response, error) {
		r.Header.Add("D", "d")
		r.Header.Add("A", "z")
		r.Header.Set("B", "x")
		return &http.Response{Request: r, StatusCode: 204}, nil
	}), head)

	r, _ := http.NewRequest("", "/", nil)
	res, err := d.Do(r)
	require.NoError(t, err)

	assert.Equal(t, 204, res.StatusCode, "status")
	assert.Equal(t, 3, len(head), "length of middleware")
	assert.Equal(t, 4, len(r.Header), "length of request")

	assert.Equal(t, map[string][]string{
		"A": {"a"},
		"B": {"b"},
		"C": {"c"},
	}, map[string][]string(head), "middleware content")

	assert.Equal(t, map[string][]string{
		"A": {"a", "z"},
		"B": {"x"},
		"C": {"c"},
		"D": {"d"},
	}, map[string][]string(r.Header), "request content")
}
