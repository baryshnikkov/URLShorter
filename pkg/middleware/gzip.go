package middleware

import (
	"URLShorter/pkg/res"
	"bytes"
	"compress/gzip"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type wrapperWriterGzip struct {
	http.ResponseWriter
	body             *bytes.Buffer
	isAcceptEncoding bool
	statusCode       int
}

func (w *wrapperWriterGzip) Write(b []byte) (int, error) {
	w.body.Write(b)

	return len(b), nil
}

func (w *wrapperWriterGzip) WriteHeader(statusCode int) {
	if w.isAcceptEncoding {
		w.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	}
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAcceptEncoding := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

		rec := &wrapperWriterGzip{
			ResponseWriter:   w,
			body:             &bytes.Buffer{},
			isAcceptEncoding: isAcceptEncoding,
			statusCode:       http.StatusOK,
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			body, err := decompress(r.Body)
			if err != nil {
				zap.L().Error("decompress gzip", zap.Error(err))
				res.JSON(w, "Failed to decompress gzip request", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(body))
			r.ContentLength = int64(len(body))
			r.Header.Del("Content-Encoding")
		}

		next.ServeHTTP(rec, r)

		if isAcceptEncoding {
			newBody, err := compress(rec.body.Bytes())
			if err != nil {
				zap.L().Error("compress gzip", zap.Error(err))
				res.JSON(w, "Failed to compress gzip request", http.StatusInternalServerError)
				return
			}

			w.Write(newBody)
		} else {
			w.Write(rec.body.Bytes())
		}
	})
}

func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}
	err = gz.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decompress(reader io.ReadCloser) ([]byte, error) {
	defer reader.Close()
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	var result bytes.Buffer
	_, err = io.Copy(&result, gzipReader)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}
