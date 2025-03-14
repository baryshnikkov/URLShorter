package res

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSON(w http.ResponseWriter, payload any, code int) {
	const op = "res.JSON"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Printf("OP: %s; DIS: Error encoding response; ERROR: %s", op, err)
	}
}
