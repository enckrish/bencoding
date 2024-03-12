package bencoding

import (
	"bytes"
	"sort"
	"strconv"
	"strings"
)

func EncodeInt(i BInt) string {
	return "i" + i.String() + "e"
}

func EncodeStr(s BStr) string {
	l := len(s)
	return strconv.Itoa(l) + ":" + s
}

func EncodeList(a BList) string {
	e := strings.Builder{}
	e.WriteByte('l')
	for _, v := range a {
		e.WriteString(EncodeAny(v))
	}
	e.WriteByte('e')

	return e.String()
}

func EncodeDict(d BDict) string {
	// Sorting keys before encoding
	keys := make([][]byte, 0, len(d))
	for k := range d {
		keys = append(keys, []byte(k))
	}

	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i], keys[j]) == -1
	})
	e := strings.Builder{}
	e.WriteByte('d')
	for _, k := range keys {
		key := string(k)
		e.WriteString(EncodeStr(key))
		e.WriteString(EncodeAny(d[key]))
	}

	e.WriteByte('e')
	return e.String()
}

func EncodeAny(v any) string {
	switch v.(type) {
	case BInt:
		return EncodeInt(v.(BInt))
	case BStr:
		return EncodeStr(v.(BStr))
	case BList:
		return EncodeList(v.(BList))
	case BDict:
		return EncodeDict(v.(BDict))
	default:
		return ""
	}
}
