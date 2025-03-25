package req

import (
	"URLShorter/pkg/res"
	"bytes"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func JSON[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	const op = "req.JSON"

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	payload, err := decode[T](r)
	if err != nil {
		logAndRes(w, r, op, "Error decoding request", err, http.StatusBadRequest, body)
		return nil, err
	}

	err = validate(payload)
	if err != nil {
		logAndRes(w, r, op, "Error validating request", err, http.StatusBadRequest, body)
		return nil, err
	}

	return payload, nil
}

func logAndRes(w http.ResponseWriter, r *http.Request, op, msg string, err error, status int, body []byte) {
	zap.L().Error(msg,
		zap.String("op", op),
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
		zap.ByteString("payload", body),
		zap.Error(err),
	)

	res.JSON(w, msg, status)
}
