package main

import (
	"fmt"
	"github.com/LiveStudioSolution/h264decoder/internal"
	"github.com/LiveStudioSolution/h264decoder/internal/logger"
	"os"
)

var log = logger.Log

func init() {
	logger.EnableStderrLog()
	log = logger.Log
}
func main() {
	ifile := "docs/videosamples/txjg.h264"
	h264Reader, err := os.Open(ifile)
	if err != nil {
		fmt.Printf("open file error %v\n", err)
		os.Exit(1)
	}
	bs := internal.NewBitStream(h264Reader)
	var nl *internal.Nalu
	for {
		nl, err = bs.NextNalu();
		if err != nil || nl == nil {
			break
		}
		log.Printf("get nalu type = %v, rbr size = %v\n", nl.Type(), nl.RbspSize())
	}
	if err != nil {
		log.Printf("get nalu err = %v\n", err)
	}
}
