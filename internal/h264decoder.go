package internal

import (
	"fmt"
	"github.com/LiveStudioSolution/h264decoder/internal/logger"
	"image"
	"os"
)

//H264Decoder  decoder of h264 codec
type H264Decoder struct {
	bs  *BitStream
	sps *SPS
	pps *PPS
}

func NewH264DecoderWithFile(filePath string) (*H264Decoder, error) {
	hd := &H264Decoder{}
	if err := hd.InitWithFile(filePath); err != nil {
		return nil, err
	}
	return hd, nil
}

func (hd *H264Decoder) InitWithFile(filePath string) error {
	iFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	hd.bs = NewBitStream(iFile)
	return nil
}

func (hd *H264Decoder) NextFrame() (image.Image, error) {
	nalu, err := hd.bs.NextNalu()
	if err != nil {
		return nil, err
	}
	switch nalu.uType {
	case NaluUnspecified:
		return nil, fmt.Errorf("NaluUnspecified")
	case NaluSlice:
		return nil, fmt.Errorf("NaluSlice")
	case NaluSliceDpa:
		return nil, fmt.Errorf("NaluSliceDpa")
	case NaluSliceDpb:
		return nil, fmt.Errorf("NaluSliceDpb")
	case NaluSliceDpc:
		return nil, fmt.Errorf("NaluSliceDpc")
	case NaluSliceIdr:
		return nil, fmt.Errorf("NaluSliceIdr")
	case NaluSei:
		return nil, fmt.Errorf("NaluSei")
	case NaluSps:
		return nil, hd.parseSps(nalu)
	case NaluPps:
		return nil, fmt.Errorf("NaluPps")
	case NaluAud:
		return nil, fmt.Errorf("NaluAud")
	case NaluEoseq:
		return nil, fmt.Errorf("NaluEoseq")
	case NaluEostream:
		return nil, fmt.Errorf("NaluEostream")
	case NaluFiller:
		return nil, fmt.Errorf("NaluFiller")
	}
	return nil, nil
}

func (hd *H264Decoder) parseSps(nalu *Nalu) error {
	var err error
	hd.sps, err = SpsParseFromRBSP(nalu.rbsp)
	if err != nil {
		return err
	}
	l := logger.Log
	l.Printf("got sps %v", hd.sps)
	return nil
}
