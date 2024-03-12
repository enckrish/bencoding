package bencoding

import "math/big"

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
