package bencoding

import (
	"bytes"
	"errors"
	"io"
	"math/big"
	"strings"
)

var (
	ErrBigInt       = errors.New("bigint conversion failed")
	ErrInvalidInput = errors.New("input doesn't contain valid bencoded data")
	ErrKeysOrdering = errors.New("keys must appear in sorted order (sorted as raw strings, not alphanumerics)")
)

type BTypeId int

const (
	BIntId BTypeId = iota
	BStrId
	BListId
	BDictId
)

type BInt = *big.Int
type BStr = string
type BList = []any
type BDict = map[string]any

func AssertType(item any, t BTypeId) bool {
	switch item.(type) {
	case BInt:
		return t == BIntId
	case BStr:
		return t == BStrId
	case BList:
		return t == BListId
	case BDict:
		return t == BDictId
	default:
		return false
	}
}

func consumeToken(b *BufStream, s byte) error {
	c, err := b.PeekNext()
	if err != nil {
		return err
	}

	if c != s {
		return ErrInvalidInput
	}

	_, err = b.GetNext() // `err` guaranteed to be `nil` in current impl
	return err
}

func ReadInt(b *BufStream) (BInt, error) {
	var c byte
	var err error

	if err = consumeToken(b, 'i'); err != nil {
		return nil, err
	}

	var sb strings.Builder
	var prev byte // previously written byte

	for {
		c, err = b.PeekNext()
		if err != nil {
			return nil, err
		}

		if (c < '0' || c > '9') && c != '-' {
			break
		}

		if (prev == '0' && sb.Len() == 1) || (c == '0' && prev == '-') {
			return nil, ErrInvalidInput
		}

		sb.WriteByte(c)
		prev = c

		_, _ = b.GetNext()
	}

	if err := consumeToken(b, 'e'); err != nil {
		return nil, err
	}

	v, ok := new(big.Int).SetString(sb.String(), 10)
	if !ok {
		return v, ErrBigInt
	}

	return v, nil
}

func ReadStr(b *BufStream) (BStr, error) {
	length := 0
	for {
		c, err := b.GetNext()
		if err != nil {
			return "", err
		}

		if c == ':' {
			break
		}

		if c < '0' || c > '9' {
			return "", ErrInvalidInput
		}

		length = length*10 + int(c-'0')
	}

	var sb strings.Builder
	for i := 0; i < length; i++ {
		c, err := b.GetNext()
		if err != nil {
			return "", err
		}
		sb.WriteByte(c)
	}

	return sb.String(), nil
}

func ReadList(b *BufStream) (BList, error) {
	if err := consumeToken(b, 'l'); err != nil {
		return nil, err
	}

	var list BList
	for {
		c, err := b.PeekNext()
		if err != nil {
			return nil, err
		}

		if c == 'e' {
			_, _ = b.GetNext()
			break
		}

		v, err := ReadAny(b)
		if err != nil {
			return nil, err
		}
		list = append(list, v)
	}

	//fmt.Println("Read List: ", list)
	return list, nil
}

func ReadDict(b *BufStream) (BDict, error) {
	if err := consumeToken(b, 'd'); err != nil {
		return nil, err
	}

	m := make(BDict)

	var prev []byte // previous key
	for {
		c, err := b.PeekNext()
		if err != nil {
			return nil, err
		}

		if c == 'e' {
			_, _ = b.GetNext()
			break
		}

		k, err := ReadStr(b)
		if err != nil {
			return nil, err
		}

		if bytes.Compare(prev, []byte(k)) == 1 {
			return nil, ErrKeysOrdering
		}
		prev = []byte(k)

		v, err := ReadAny(b)
		if err != nil {
			return nil, err
		}

		m[k] = v
	}

	//fmt.Println("Read dict: ", m)
	return m, nil
}

func ReadAny(b *BufStream) (any, error) {
	c, err := b.PeekNext()
	if err != nil {
		return nil, err
	}

	var v any
	switch {
	case c >= '0' && c <= '9':
		v, err = ReadStr(b)
	case c == 'i':
		v, err = ReadInt(b)
	case c == 'l':
		v, err = ReadList(b)
	case c == 'd':
		v, err = ReadDict(b)
	default:
		return nil, ErrInvalidInput
	}

	if err == io.EOF {
		err = ErrInvalidInput
	}
	return v, err
}
