package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"GRPC-KV-Store-System/api-service/internal/handler"
)

func setupRouter() *mux.Router {
	mockClient := NewMockClient()
	h := handler.StartHandler(mockClient)

	router := mux.NewRouter()
	router.HandleFunc("/health", h.HealthHandler).Methods("GET")
	router.HandleFunc("/kv", h.SetHandler).Methods("POST")
	router.HandleFunc("/kv/{key}", h.GetHandler).Methods("GET")
	router.HandleFunc("/kv/{key}", h.DeleteHandler).Methods("DELETE")

	return router
}

func TestHealthEndpoint(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)

	if response["status"] != "healthy" {
		t.Errorf("Expected healthy status, got %s", response["status"])
	}

	t.Log("Health check passed")
}

func TestSetAndGetFlow(t *testing.T) {
	router := setupRouter()

	// Test Set
	t.Run("Set Key-Value", func(t *testing.T) {
		body := []byte(`{"key":"username","value":"alice"}`)
		req := httptest.NewRequest("POST", "/kv", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", rr.Code)
		}

		t.Log("Set successful")
	})

	// Test Get
	t.Run("Get Key-Value", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/kv/username", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		var response map[string]string
		json.NewDecoder(rr.Body).Decode(&response)

		if response["value"] != "alice" {
			t.Errorf("Expected 'alice', got '%s'", response["value"])
		}

		t.Log("Get successful")
	})
}

func TestGetNonExistentKey(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest("GET", "/kv/nonexistent", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)

	if response["error"] != "key not found" {
		t.Errorf("Expected 'key not found', got '%s'", response["error"])
	}

	t.Log("Not found handled correctly")
}

func TestDeleteFlow(t *testing.T) {
	router := setupRouter()

	// Set a key first
	body := []byte(`{"key":"temp","value":"data"}`)
	req := httptest.NewRequest("POST", "/kv", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Delete the key
	t.Run("Delete Key", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/kv/temp", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		t.Log("Delete successful")
	})

	// Verify it's deleted
	t.Run("Verify Deletion", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/kv/temp", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status 404 after delete, got %d", rr.Code)
		}

		t.Log("Key successfully deleted")
	})
}

func TestDeleteNonExistentKey(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest("DELETE", "/kv/nonexistent", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}

	t.Log("Delete non-existent key handled correctly")
}

func TestInvalidJSON(t *testing.T) {
	router := setupRouter()

	body := []byte(`{invalid json}`)
	req := httptest.NewRequest("POST", "/kv", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}

	t.Log("Invalid JSON rejected")
}

func TestEmptyKey(t *testing.T) {
	router := setupRouter()

	body := []byte(`{"key":"","value":"test"}`)
	req := httptest.NewRequest("POST", "/kv", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}

	t.Log("Empty key rejected")
}

func TestMultipleOperations(t *testing.T) {
	router := setupRouter()

	keys := []struct {
		key   string
		value string
	}{
		{"user:1", "alice"},
		{"user:2", "bob"},
		{"user:3", "charlie"},
	}

	// Set multiple keys
	for _, kv := range keys {
		body := []byte(`{"key":"` + kv.key + `","value":"` + kv.value + `"}`)
		req := httptest.NewRequest("POST", "/kv", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Failed to set %s", kv.key)
		}
	}

	// Verify all keys
	for _, kv := range keys {
		req := httptest.NewRequest("GET", "/kv/"+kv.key, nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Failed to get %s", kv.key)
		}

		var response map[string]string
		json.NewDecoder(rr.Body).Decode(&response)

		if response["value"] != kv.value {
			t.Errorf("Expected %s, got %s", kv.value, response["value"])
		}
	}

	t.Log("Multiple operations successful")
}
