package pickjson

import (
	"bytes"
	"io"
	"net/http"
)

type responseStealer struct {
	rw  http.ResponseWriter
	buf *bytes.Buffer
}

func newResponseStealer(rw http.ResponseWriter, buf *bytes.Buffer) *responseStealer {
	return &responseStealer{
		rw:  rw,
		buf: buf,
	}
}

var _ = (*responseStealer)(nil)

func (s *responseStealer) WriteHeader(code int) {
	s.rw.Header().Del("Content-Length")
	s.rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	s.rw.WriteHeader(code)
}

func (s *responseStealer) Header() http.Header {
	return s.rw.Header()
}

func (s *responseStealer) Write(b []byte) (int, error) {
	return s.buf.Write(b)
}

func (s *responseStealer) Steal() ([]byte, error) {
	return io.ReadAll(s.buf)
}
