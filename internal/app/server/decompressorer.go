package server

import (
	"fmt"
	"io"

	"github.com/Schalure/urlalias/internal/app/compressor"
)

// Compress type
type CompressType string

// Compress type enumeration
const (
	TypeGZIP   CompressType = "gzip"
	TypeLZW    CompressType = "compress"
	TypeZLIB   CompressType = "deflate"
	TypeBrotli CompressType = "br"
	TypeNone   CompressType = ""
)

// PossibleCompressionTypes
var PossibleCompressionTypes = []CompressType{
	TypeGZIP,
}

// Decompressorer interface
type Decompressorer interface {
	Read(p []byte) (n int, err error)
	Close() error
	Type() string
}

// Decompressorer constructor
func NewDecompressorer(r io.Reader, compressType CompressType) (Decompressorer, error) {

	var err error
	switch compressType {
	case TypeGZIP:
		if decompressorGZIP, err := compressor.NewDecompressorGZIP(r); err == nil {
			return decompressorGZIP, nil
		}
	default:
		err = fmt.Errorf("%s compression type is not supported", compressType)
	}
	return nil, err
}
