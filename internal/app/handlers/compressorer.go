package handlers

import (
	"io"

	"github.com/Schalure/urlalias/internal/app/compressor"
)

type Compressorer interface {
	Write(p []byte) (n int, err error)
	Close() error
	Type() string
}

func NewCompressorer(w io.Writer, compressType CompressType) Compressorer {

	switch compressType {
	case TypeGZIP:
		return compressor.NewCompressorGZIP(w)
	}
	return nil
}
