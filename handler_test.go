package forge

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateClient(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "customers",
			"attributes": map[string]interface{}{
				"name":  "Test Client",
				"email": "test@example.com",
			},
			"relationships": map[string]interface{}{
				"organization": map[string]interface{}{
					"data": map[string]interface{}{
						"type": "organizations",
						"id":   "org-id",
					},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", "/v1/direct_debit_collection_clients", bytes.NewReader(payloadBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createClientHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Add your response validation here
}
