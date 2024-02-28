package server

import "io"

// Decoder interface
type Decoder interface {
	Decode(r io.Reader, v interface{}) error
}
