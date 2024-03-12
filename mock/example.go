package main

import (
	"bencoding"
	"encoding/json"
	"fmt"
	"os"
)

const filename = "mock/data/debian-12.5.0-amd64-netinst.iso.torrent"

func main() {
	// Reading bencoded file
	m := Read()

	// Pretty printing
	j, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println(string(j))

	// Encoding the created dictionary
	enc := bencoding.EncodeAny(m)
	if enc == "" {
		panic("Encoding failed!")
	}
	fmt.Println(enc)
}

func Read() map[string]any {
	f, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	b := bencoding.NewBufStream(f, 1<<10)
	o, _ := bencoding.ReadAny(b)

	if !bencoding.AssertType(o, bencoding.BDictId) {
		panic("bencoded file doesn't represent a dictionary")
	}

	m := o.(bencoding.BDict)
	return m
}
