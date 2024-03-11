package bencoding

import (
	"bufio"
	"io"
)

// BufStream streaming reader with buffer with single byte peek
type BufStream struct {
	r      *bufio.Reader // reader provided for reads
	buffer []byte        // buffer for storing reads
	length int           // number of bytes from ptr to buffer end
	ptr    int           // current location of read pointer
}

func NewBufStream(r io.Reader, size int) *BufStream {
	return &BufStream{
		r:      bufio.NewReader(r),
		buffer: make([]byte, size),
	}
}

// loadIfEmpty fills the buffer if no more elements are left to read
func (b *BufStream) loadIfEmpty() error {
	if b.ptr < b.length {
		return nil
	}

	n, err := b.r.Read(b.buffer)
	if err != nil {
		return err
	}

	b.length = n
	b.ptr = 0

	return nil
}

// Clear resets the buffer state
func (b *BufStream) Clear() {
	b.length = 0
	b.ptr = 0
}

// GetNext returns item at current position, and increments ptr
func (b *BufStream) GetNext() (byte, error) {
	if err := b.loadIfEmpty(); err != nil {
		return 0, err
	}

	b.ptr++
	return b.buffer[b.ptr-1], nil
}

// PeekNext returns item at current position without incrementing ptr
func (b *BufStream) PeekNext() (byte, error) {
	if err := b.loadIfEmpty(); err != nil {
		return 0, err
	}
	return b.buffer[b.ptr], nil
}
