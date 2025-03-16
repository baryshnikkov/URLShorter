package req

import (
	"URLShorter/pkg/res"
	"fmt"
	"log"
	"net/http"
)

func JSON[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	path := r.URL.Path
	method := r.Method
	op := fmt.Sprintf("req.JSON path=%s method=%s", path, method)

	payload, err := decode[T](r)
	if err != nil {
		log.Printf("OP: %s; DIS: error decoding request; ERROR: %s", op, err)
		res.JSON(w, "Error decoding request", http.StatusBadRequest)
		return nil, err
	}

	err = validate(payload)
	if err != nil {
		log.Printf("OP: %s; DIS: error validate request; ERROR: %s", op, err)
		res.JSON(w, "Error validate request", http.StatusBadRequest)
		return nil, err
	}

	return payload, nil
}
