package internal

import "io"

type Nalu struct {
	src io.Reader
}

func NewNalu() *Nalu {
	nl := Nalu{}
	return &nl
}
func (this *Nalu) SetSource(src io.Reader) {
	this.src = src
}
