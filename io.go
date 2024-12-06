package conduit

import (
	"fmt"
	"io"
)

// LimitedReader is an io.Reader that enforces a maximum read limit.
// If more than 'limit' bytes are read, it returns an error.
//
// This is useful for preventing large malicious payloads from consuming too much memory.
type LimitedReader struct {
	r        io.Reader
	limit    int64
	consumed int64
}

// NewLimitedReader creates a new LimitedReader that wraps the provided io.Reader 'r'
// and enforces a maximum allowed read size of 'limit' bytes.
func NewLimitedReader(r io.Reader, limit int64) *LimitedReader {
	return &LimitedReader{
		r:     r,
		limit: limit,
	}
}

// Read reads up to len(p) bytes into p from the underlying reader.
// If reading would exceed the limit, it returns an error.
func (l *LimitedReader) Read(p []byte) (n int, err error) {
	n, err = l.r.Read(p)
	l.consumed += int64(n)
	if l.consumed > l.limit {
		return n, fmt.Errorf("message size exceeds limit of %d bytes", l.limit)
	}
	return n, err
}
