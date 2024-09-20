package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// setupTestServer initializes the server handler for testing.
func setupTestServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", ghostHandler) // Handle '/echo' explicitly
	mux.HandleFunc("/echo/", echoHandler)
	mux.HandleFunc("/", ghostHandler) // Default handler
	return logRequest(mux)
}

// Test valid statement in /echo/:statement
func TestEchoHandler_ValidStatement(t *testing.T) {
	handler := setupTestServer()
	req := httptest.NewRequest("GET", "/echo/hello", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
	}

	expectedBody := "hello hello\n"
	if string(body) != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, string(body))
	}
}

// Test /echo and /echo/ with no statement should return the ghostHandler spooky message
func TestEchoHandler_NoStatement(t *testing.T) {
	handler := setupTestServer()
	paths := []string{"/echo", "/echo/"}

	for _, path := range paths {
		req := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 Not Found, got %d for path %s", resp.StatusCode, path)
		}

		expectedBody := spookyMessage + "\n"
		if string(body) != expectedBody {
			t.Errorf("Expected body %q, got %q for path %s", expectedBody, string(body), path)
		}
	}
}

// Test unknown route should return the ghostHandler spooky message
func TestUnknownRoute(t *testing.T) {
	handler := setupTestServer()
	req := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 Not Found, got %d", resp.StatusCode)
	}

	expectedBody := spookyMessage + "\n"
	if string(body) != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, string(body))
	}
}

// Test POST or any non-GET request should return 405 Method Not Allowed
func TestMethodNotAllowed(t *testing.T) {
	handler := setupTestServer()
	req := httptest.NewRequest("POST", "/echo/hello", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 Method Not Allowed, got %d", resp.StatusCode)
	}

	expectedBody := "Method not allowed\n"
	if string(body) != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, string(body))
	}
}

// Test bad URL encoding in /echo/:statement should return 404 and call ghostHandler
func TestEchoHandler_BadURL(t *testing.T) {
	handler := setupTestServer()

	// Manually create a request with an invalid URL path
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path: "/echo/%", // Invalid URL path
		},
		Host: "localhost",
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 Not Found, got %d", resp.StatusCode)
	}

	expectedBody := spookyMessage + "\n"
	if string(body) != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, string(body))
	}
}
