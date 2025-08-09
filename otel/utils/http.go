package utils

import "net/http"

type CaptureHTTPResponseWriter struct {
	http.ResponseWriter

	status int
	size   int64
}

func (w *CaptureHTTPResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *CaptureHTTPResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	if err != nil {
		return n, err
	}

	w.size += int64(n)

	return n, nil
}

func (w *CaptureHTTPResponseWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}

	return w.status
}

func (w *CaptureHTTPResponseWriter) Size() int64 {
	return w.size
}
