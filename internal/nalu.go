package internal

import (
	"bytes"
	"fmt"

	"github.com/32bitkid/bitreader"
)

// NaluType of h264 codec
type NaluType uint8

// ITU-T Rec. H.264 (05/2003) page 48
// Table 7-1 – NAL unit type codes
const (
	NaluUnspecified NaluType = 0
	// Coded slice of a non-IDR picture
	NaluSlice NaluType = 1
	// Coded slice data partition A
	NaluSliceDpa NaluType = 2
	// Coded slice data partition B
	NaluSliceDpb NaluType = 3
	// Coded slice data partition C
	NaluSliceDpc NaluType = 4
	// Coded slice of an IDR picture
	NaluSliceIdr NaluType = 5
	// Supplemental enhancement information (SEI)
	NaluSei NaluType = 6
	// Sequence parameter set
	NaluSps NaluType = 7
	// Picture parameter set
	NaluPps NaluType = 8
	// Access unit delimiter
	NaluAud NaluType = 9
	// End of sequence
	NaluEoseq NaluType = 10
	// End of stream
	NaluEostream NaluType = 11
	// Filler data
	NaluFiller NaluType = 12
)

func (nt NaluType) String() string {
	switch nt {
	case NaluUnspecified:
		return "NaluUnspecified"
	case NaluSlice:
		return "NaluSlice"
	case NaluSliceDpa:
		return "NaluSliceDpa"
	case NaluSliceDpb:
		return "NaluSliceDpb"
	case NaluSliceDpc:
		return "NaluSliceDpc"
	case NaluSliceIdr:
		return "NaluSliceIdr"
	case NaluSei:
		return "NaluSei"
	case NaluSps:
		return "NaluSps"
	case NaluPps:
		return "NaluPps"
	case NaluAud:
		return "NaluAud"
	case NaluEoseq:
		return "NaluEoseq"
	case NaluEostream:
		return "NaluEostream"
	case NaluFiller:
		return "NaluFiller"
	}
	return fmt.Sprintf("NaluUnspecified:%d", nt)
}

// Nalu  of h264 codec
type Nalu struct {
	rbsp   []byte
	refIdc uint8
	uType  NaluType
	br     bitreader.BitReader
}

//NewNalu  create new Nalu
func NewNalu() *Nalu {
	nl := Nalu{}
	return &nl
}

//Load load bit from data, add parse nalu fields
func (nl *Nalu) Load(data []byte) error {
	if data == nil || len(data) < 1 {
		return fmt.Errorf("invalid nalu data")
	}
	nl.rbsp = data[1:]
	nl.br = bitreader.NewReader(bytes.NewBuffer(data))
	return nl.parse()
}

func (nl *Nalu) parse() error {
	if len(nl.rbsp) < 1 {
		return fmt.Errorf("nalu invalid rbr size 0")
	}
	// parse type
	// forbidden_zero_bit
	t, err := nl.br.Read1()
	if err != nil {
		return err
	}
	if t == true {
		return fmt.Errorf("nalu invalid forbidden zero bit")
	}
	// nal_ref_idc
	nl.refIdc, err = nl.br.Read8(2)
	if err != nil {
		return err
	}
	nt, err := nl.br.Read8(5)
	nl.uType = NaluType(nt)
	if nl.uType > NaluFiller {
		return fmt.Errorf("invalid nalu type %v", nl.uType)
	}
	return nil
}

// Type return nalu type
func (nl *Nalu) Type() NaluType {
	return nl.uType
}

// RbspSize  return nalu rbr bytes count
func (nl *Nalu) RbspSize() int {
	return len(nl.rbsp)
}
