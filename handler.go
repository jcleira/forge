package forge

import (
	"encoding/json"
	"net/http"
)

func createClientHandler(w http.ResponseWriter, r *http.Request) {
	var payload map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process the request and generate the response
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"id":         "client-id",
			"type":       "customers",
			"attributes": payload["data"].(map[string]interface{})["attributes"],
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
