package res

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

func JSON(w http.ResponseWriter, payload any, code int) {
	const op = "res.JSON"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		zap.L().Error("Error encoding response",
			zap.String("op", op),
			zap.Error(err))
	}
}
