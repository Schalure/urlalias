package compressor

import (
	"compress/gzip"
	"io"
)

// DecompressorGZIP type
type DecompressorGZIP struct {
	zr *gzip.Reader
}

// ------------------------------------------------------------
//
//	Constructor of DecompressorGZIP
//	Input:
//		r io.Reader
//	Output:
//		*DecompressorGZIP
//		error
func NewDecompressorGZIP(r io.Reader) (*DecompressorGZIP, error) {

	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &DecompressorGZIP{
		zr: zr,
	}, nil
}

// Read
func (c *DecompressorGZIP) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close
func (c *DecompressorGZIP) Close() error {
	return c.zr.Close()
}

// Type
func (c *DecompressorGZIP) Type() string {
	return "gzip"
}

// CompressorGZIP type
type CompressorGZIP struct {
	zw *gzip.Writer
}

// ------------------------------------------------------------
//
//	Constructor of CompressorGZIP
//	Input:
//		r io.Reader
//	Output:
//		*CompressorGZIP
//		error
func NewCompressorGZIP(w io.Writer) *CompressorGZIP {

	return &CompressorGZIP{
		zw: gzip.NewWriter(w),
	}
}

// Write
func (c *CompressorGZIP) Write(p []byte) (n int, err error) {
	return c.zw.Write(p)
}

// Close
func (c *CompressorGZIP) Close() error {
	return c.zw.Close()
}

// Type
func (c *CompressorGZIP) Type() string {
	return "gzip"
}
