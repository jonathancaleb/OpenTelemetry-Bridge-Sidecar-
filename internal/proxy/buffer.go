package proxy

import (
	"io"
	"sync/atomic"
)

// CountingReader wraps an io.Reader to count bytes read without breaking streaming.
// It implements the io.Reader and io.ReadCloser interfaces.
type CountingReader struct {
	reader    io.Reader
	bytesRead int64
}

// NewCountingReader creates a new CountingReader that wraps the given reader.
// The wrapped reader passes all data through unchanged while tracking the total bytes read.
func NewCountingReader(r io.Reader) *CountingReader {
	return &CountingReader{
		reader:    r,
		bytesRead: 0,
	}
}

// Read implements io.Reader. It reads from the underlying reader and counts the bytes.
func (cr *CountingReader) Read(p []byte) (n int, err error) {
	n, err = cr.reader.Read(p)
	if n > 0 {
		atomic.AddInt64(&cr.bytesRead, int64(n))
	}
	return n, err
}

// Close implements io.Closer. If the underlying reader is a Closer, it closes it.
func (cr *CountingReader) Close() error {
	if closer, ok := cr.reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// BytesRead returns the total number of bytes read so far.
// This is safe to call concurrently with Read operations.
func (cr *CountingReader) BytesRead() int64 {
	return atomic.LoadInt64(&cr.bytesRead)
}

// Reset resets the byte counter to zero.
func (cr *CountingReader) Reset() {
	atomic.StoreInt64(&cr.bytesRead, 0)
}
