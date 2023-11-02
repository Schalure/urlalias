package compressor

import (
	"compress/gzip"
	"io"
)

// ========================================================================================================================
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

func (c *DecompressorGZIP) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *DecompressorGZIP) Close() error {
	return c.zr.Close()
}

func (c *DecompressorGZIP) Type() string {
	return "gzip"
}

// ========================================================================================================================
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

func (c *CompressorGZIP) Write(p []byte) (n int, err error) {
	return c.zw.Write(p)
}

func (c *CompressorGZIP) Close() error {
	return c.zw.Close()
}

func (c *CompressorGZIP) Type() string {
	return "gzip"
}
