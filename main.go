package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

const (
	defaultPort   = "1323"
	spookyMessage = "oooOOOoooOOOooo"
)

func main() {
	port := getServerPort()

	srv := setupServer(port)

	go func() {
		log.Printf("Starting server on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	waitForShutdown(srv)
}

// setupServer configures and returns a new HTTP server.
func setupServer(port string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/echo", ghostHandler) // Handle '/echo' explicitly
	mux.HandleFunc("/echo/", echoHandler)
	mux.HandleFunc("/", ghostHandler) // Default handler

	handler := logRequest(mux)

	return &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
}

// getServerPort returns the server port from the environment variable or a default value.
func getServerPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return defaultPort
}

// waitForShutdown handles server shutdown gracefully on receiving a termination signal.
func waitForShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}

// echoHandler handles /echo/:statement requests.
func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	statement := r.URL.Path[len("/echo/"):]
	if statement == "" {
		ghostHandler(w, r)
		return
	}

	decoded, err := url.PathUnescape(statement)
	if err != nil {
		ghostHandler(w, r)
		return
	}

	fmt.Fprintf(w, "%s %s\n", decoded, decoded)
}

// ghostHandler handles all other requests with a spooky message.
func ghostHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, spookyMessage)
}

// logRequest logs each HTTP request along with the response status code.
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)
		log.Printf(
			"%s - %s %s %s %d %dms %s",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			r.Proto,
			lrw.statusCode,
			duration.Milliseconds(),
			r.UserAgent(),
		)
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture the status code.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code for logging.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
