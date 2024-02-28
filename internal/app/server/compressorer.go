package server

import (
	"io"

	"github.com/Schalure/urlalias/internal/app/compressor"
)

//	Compressorer type
type Compressorer interface {
	Write(p []byte) (n int, err error)
	Close() error
	Type() string
}

//	Compressorer constructor
func NewCompressorer(w io.Writer, compressType CompressType) Compressorer {

	switch compressType {
	case TypeGZIP:
		return compressor.NewCompressorGZIP(w)
	}
	return nil
}
