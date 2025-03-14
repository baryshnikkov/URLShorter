package req

import (
	"encoding/json"
	"net/http"
)

func decode[T any](r *http.Request) (*T, error) {
	var payload T
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
