package internal

import "image"

//H264Decoder  decoder of h264 codec
type H264Decoder struct {
	bs  *BitStream
	sps *SPS
	pps *PPS
}

func (hd *H264Decoder) InitWithFile(filePath string) error {
	return nil
}

func (hd *H264Decoder) NextFrame() (image.Image, error) {

	return nil, nil
}
