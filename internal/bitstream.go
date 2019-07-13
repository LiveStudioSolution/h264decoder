package internal

import (
	"bufio"
	"io"
)

var (
	annexBSpliter1 = []byte{0, 0, 1}
	annexBSpliter2 = []byte{0, 0, 0, 1}
)

func init() {

}

// BitStream H264 bit stream, contains nalu
type BitStream struct {
	src     io.Reader
	scanner bufio.Scanner
}
