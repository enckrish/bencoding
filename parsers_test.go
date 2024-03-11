package bencoding

import (
	"os"
	"testing"
)

func BenchmarkReadAny(b *testing.B) {
	f, _ := os.Open("mock/data/debian-12.5.0-amd64-netinst.iso.torrent")

	bs := NewBufStream(f, 1<<20)

	for i := 0; i < b.N; i++ {
		// Reset seek to minimize deviations not due to parsers
		_, _ = f.Seek(0, 0)
		bs.Clear()

		_, err := ReadAny(bs)
		if err != nil {
			panic(err)
		}
	}
}
