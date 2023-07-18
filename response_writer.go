package redmid

import (
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter

	code int
	size int
	data []byte
}

func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	return ResponseWriter{
		ResponseWriter: w,
	}
}

func (r *ResponseWriter) WriteHeader(code int) {
	if r.Code() == 0 {
		r.code = code
		r.ResponseWriter.WriteHeader(code)
	}
}

// Overrides `http.ResponseWriter` type.
func (r *ResponseWriter) Write(body []byte) (int, error) {
	if r.Code() == 0 {
		r.WriteHeader(http.StatusOK)
	}

	r.data = body

	var err error
	r.size, err = r.ResponseWriter.Write(body)

	return r.size, err
}

func (r *ResponseWriter) Flush() {
	if fl, ok := r.ResponseWriter.(http.Flusher); ok {
		if r.Code() == 0 {
			r.WriteHeader(http.StatusOK)
		}

		fl.Flush()
	}
}

func (r *ResponseWriter) Code() int {
	return r.code
}

func (r *ResponseWriter) Size() int {
	return r.size
}

func (r *ResponseWriter) Data() []byte {
	return r.data
}
