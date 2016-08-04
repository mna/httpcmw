// Copyright 2016 Martin Angers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpcmw

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrintfLogger(t *testing.T) {
	var buf bytes.Buffer
	stdl := log.New(&buf, "", 0)
	l := PrintfLogger(stdl.Printf)

	cases := []struct {
		in  []interface{}
		out string
	}{
		{nil, ""},
		{[]interface{}{"a"}, ""},
		{[]interface{}{"a", 1}, "a=1\n"},
		{[]interface{}{"a", 1, "b"}, "a=1\n"},
		{[]interface{}{"a", 1, "b", 2}, "a=1 b=2\n"},
		{[]interface{}{"a", 1, "b", 2, "c"}, "a=1 b=2\n"},
		{[]interface{}{"a", 1, "b", 2, "c", true}, "a=1 b=2 c=true\n"},
		{[]interface{}{"a", 1, "b", 2, "c", true, "d"}, "a=1 b=2 c=true\n"},
		{[]interface{}{"a", 1, "b", 2, "c", true, "d", "some value"}, "a=1 b=2 c=true d=\"some value\"\n"},
		{[]interface{}{"a", 1, "b", 2, "c", true, "d", time.Second}, "a=1 b=2 c=true d=\"1s\"\n"},
	}
	for i, c := range cases {
		err := l.Log(c.in...)
		if assert.NoError(t, err, "%d: Log", i) {
			assert.Equal(t, c.out, buf.String(), "%d: expected output", i)
		}
		buf.Reset()
	}
}
