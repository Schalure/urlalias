package handlers

import "io"

type Decoder interface {
	Decode(r io.Reader, v interface{}) error
}
