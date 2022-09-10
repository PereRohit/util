package response

import (
	"encoding/json"
	"net/http"

	"github.com/PereRohit/util/model"
)

func ToJson(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(&model.Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}
