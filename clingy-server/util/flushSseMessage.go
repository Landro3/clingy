package util;

import (
	"encoding/json"
	"net/http"
	"fmt"
)

func FlushSSEMessage(w http.ResponseWriter, data interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}

	_, err = fmt.Fprintf(w, "data: %s\n\n", string(jsonBytes))
	if err != nil {
		return fmt.Errorf("error writing to response: %v", err)
	}

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	return nil
}
