package internal

import "io"

type BitStream struct {
	src    io.Reader
	buffer bufio
}
