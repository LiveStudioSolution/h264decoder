package main

import (
	"github.com/LiveStudioSolution/h264decoder/internal"
	"github.com/LiveStudioSolution/h264decoder/internal/logger"
)
var log = logger.Log

func init() {
	logger.EnableStderrLog()
	log = logger.Log
}

func main() {
	iFile := "docs/videosamples/txjg.h264"
	h264Decoder, err := internal.NewH264DecoderWithFile(iFile)
	if err != nil {
		log.Printf("NewH264DecoderWithFile error:%v", err)
		return
	}
	for {
		frame, err := h264Decoder.NextFrame()
		if err != nil {
			log.Printf("NextFrame error %v", err)
			break
		}
		if frame != nil {
			log.Printf("got frame %v",frame)
		}
	}
}
