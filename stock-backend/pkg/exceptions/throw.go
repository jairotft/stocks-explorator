package exceptions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func isTrue(value string) bool {
	options := []string{"TRUE", "SI", "1"}
	for _, option := range options {
		if strings.ToUpper(value) == option {
			return true
		}
	}
	return false
}

func Throw(w http.ResponseWriter, exception AppException, status int, err error) {
	if exception.Code == 0 {
		exception.Code = status
	}

	payload, err := json.Marshal(exception)
	w.WriteHeader(status)
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"detail\": \"%s\"}", exception.Detail), status)
	}
	http.Error(w, string(payload), status)
}
