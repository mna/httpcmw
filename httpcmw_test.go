// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpcmw

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	var buf bytes.Buffer
	nop := DoerFunc(func(req *http.Request) (*http.Response, error) { return nil, nil })
	fn := func(char byte) WrapperFunc {
		return func(d Doer) Doer {
			return DoerFunc(func(req *http.Request) (*http.Response, error) {
				buf.WriteByte(char)
				return d.Do(req)
			})
		}
	}

	cases := []struct {
		d   Doer
		log string
	}{
		{Wrap(nop), ""},
		{Wrap(nop, fn('a')), "a"},
		{Wrap(nop, fn('a'), fn('b')), "ab"},
		{Wrap(nop, fn('a'), fn('b'), fn('c')), "abc"},
	}
	for i, c := range cases {
		buf.Reset()
		req, _ := http.NewRequest("", "/", nil)
		c.d.Do(req)

		assert.Equal(t, c.log, buf.String(), "%d: log", i)
	}
}
