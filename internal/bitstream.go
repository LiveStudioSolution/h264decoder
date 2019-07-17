package internal

import (
	"bufio"
	"bytes"
	"io"
)

var annexBSpliter1 = []byte{0, 0, 1}
var annexBSpliter2 = []byte{0, 0, 0, 1}

const (
	maxNaluRbspSize = 512 * 1024
)

func init() {

}

// BitStream H264 bit stream, contains nalu
type BitStream struct {
	src     io.Reader
	scanner *bufio.Scanner
}

// NewBitStream return a new BitStream read from src
func NewBitStream(src io.Reader) *BitStream {
	bs := BitStream{
		src: src,
	}
	bs.scanner = bufio.NewScanner(bs.src)
	bs.scanner.Buffer(make([]byte, maxNaluRbspSize), maxNaluRbspSize)
	bs.scanner.Split(ScanNalu)
	return &bs
}

// NextNalu read next nalu from src stream
func (bs *BitStream) NextNalu() (*Nalu, error) {
	scan := bs.scanner.Scan()
	if scan == false {
		return nil, bs.scanner.Err()
	}
	nl := NewNalu()
	err := nl.Load(bs.scanner.Bytes())
	if err != nil {
		return nil, err
	}
	return nl, nil
}

// ScanNalu  split func for bufio  to split nalu in bit stream
func ScanNalu(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading splite bytes.
	//if len(data) < 3 {
	//	if atEOF {
	//		return len(data), data, nil
	//	} else {
	//		return 0, nil, nil
	//	}
	//}
	//log.Printf("data = %v",data)
	start := 0
	if bytes.HasPrefix(data, annexBSpliter1) {
		start = len(annexBSpliter1)
	} else if bytes.HasPrefix(data, annexBSpliter2) {
		start = len(annexBSpliter2)
	}

	// Scan until next spliter, marking end of nalu.
	for i := start; i+3 < len(data); i++ {
		if bytes.HasPrefix(data[i:], annexBSpliter1) {
			return i + len(annexBSpliter1), data[start:i], nil
		}
		if bytes.HasPrefix(data[i:], annexBSpliter2) {
			return i + len(annexBSpliter2), data[start:i], nil
		}
	}
	// If we're at EOF, last nalu . Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}
