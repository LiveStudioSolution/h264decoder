package rbr

import (
	"fmt"
	"github.com/32bitkid/bitreader"
)

var errNotImplemented = fmt.Errorf("not implemented")

func MoreRBSPData(br bitreader.BitReader) bool {
	if _, err := br.Peek1(); err != nil {
		return false
	}
		
	panic(errNotImplemented)
	return false
}

func lastSetBitOffset(br bitreader.BitReader) int {
	lSetBitOffset := 0
	currBitOffset := 0
	for isSet, err := br.Peek1(); err != nil; isSet, err = br.Peek1() {
		currBitOffset +=1
		if isSet {
			lSetBitOffset = currBitOffset
		}
	}
	return lSetBitOffset
}
