package handlers

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

// ========================================================================================================================
//
//	Compress writer type
type compressWriter struct {
	w          http.ResponseWriter
	compressor Compressorer
}

// ------------------------------------------------------------
//
//	Constructor of compressWriter type
//	Input:
//		w http.ResponseWriter
//		compressType CompressType
//	Output:
//		*compressWriter
func newCompressWriter(w http.ResponseWriter, compressType CompressType) *compressWriter {

	compressor := NewCompressorer(w, compressType)
	return &compressWriter{
		w:          w,
		compressor: compressor,
	}
}

// ------------------------------------------------------------
//
//	Header method of compressWriter type
//	Output:
//		http.Header
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// ------------------------------------------------------------
//
//	Write method of compressWriter type
//	Input:
//		p []byte
//	Output:
//		n int - count of write bytes
//		error
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.compressor.Write(p)
}

// ------------------------------------------------------------
//
//	WriteHeader method of compressWriter type
//	Input:
//		statusCode int
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set(contentEncoding, c.compressor.Type())
	}
	c.w.WriteHeader(statusCode)
}

// ------------------------------------------------------------
//
//	 Close method of compressWriter type
//		Sends all data from the buffer.
//		Output:
//			error
func (c *compressWriter) Close() error {
	return c.compressor.Close()
}

// ========================================================================================================================
//
//	Compress reader type
type compressReader struct {
	r            io.ReadCloser
	decompressor Decompressorer
}

// ------------------------------------------------------------
//
//	Constructor of compressReader type
//	Input:
//		r io.ReadClose
//		compressType CompressType
//	Output:
//		*compressReader
//		error
func newCompressReader(r io.ReadCloser, compressType CompressType) (*compressReader, error) {

	decompressor, err := NewDecompressorer(r, compressType)
	if err != nil {
		return nil, err
	}
	return &compressReader{
		r:            r,
		decompressor: decompressor,
	}, nil
}

// ------------------------------------------------------------
//
//	Read method of compressReader type
//	Input:
//		p []byte
//	Output:
//		n int - count of read bytes
//		error
func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.decompressor.Read(p)
}

// ------------------------------------------------------------
//
//	Close method of compressReader type
//	Output:
//	error
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.decompressor.Close()
}

// ------------------------------------------------------------
//
//	WithCompress middleware - method of Middleware type
//	Compress and decompress data
//	Receiver:
//		m* Middleware
//	Input:
//		h http.Handler
//	Output:
//		http.Handler
func (m *Middleware) WithCompress(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if isCompressionPossible(r.Header.Get(contentType)){
			encodingMethod, err := getEncodingMethod(r)
		}


		contentEncodingType := r.Header.Get(contentEncoding)
		if contentEncodingType == "" {
			h.ServeHTTP(w, r)
		} else {
			compressor := NewDecompressorer(CompressType(contentEncodingType))
			compressReader, err := compressor.Decompress(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				m.logger.Info(err)
				return
			}
			newCompressReader()

		}
	})
}

func isCompressionPossible(contentType string) bool{
	
	for _, ct := range ContentTypeToCompress{
		if strings.Contains(contentType, ct){
			return true
		}
	} 
	return false
}

func getEncodingMethod(r *http.Request) CompressType{

	var methodType string
	qMax := -100.0

	acceptEncoding := r.Header.Values(acceptEncoding)

	if len(acceptEncoding) == 0 {
		return TypeNone
	}

	for _, method := range acceptEncoding{
		s := strings.Split(method, ";")
		
		if (len(s) == 2) {
			q, err := strconv.ParseFloat(strings.TrimPrefix(s[1], "q="), 32)
			if err != nil{
				q = 0;
			}

			if q > qMax && func () bool  {
				for _, t := range PossibleCompressionTypes{
					if CompressType(s[0]) == t{
						return true
					}
				}
						return false
			}() {
				methodType = s[0]
			}
		}
	}
	return CompressType(methodType)
}
