package main

import (
	"bencoding"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("mock/data/debian-12.5.0-amd64-netinst.iso.torrent")

	if err != nil {
		print("File open err: ")
		panic(err)
	}

	b := bencoding.NewBufStream(f, 1<<10)
	o, _ := bencoding.ReadAny(b)

	if !bencoding.AssertType(o, bencoding.BDictId) {
		panic("bencoded file doesn't represent a dictionary")
	}

	m := o.(bencoding.BDict)
	j, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println(string(j))
}
