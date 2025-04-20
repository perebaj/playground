package gracefulshutdown

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: setupRoutes(),
	}

	// serverErrors is a channel to receive errors from the server.
	// It is buffered to avoid blocking the goroutine that starts the server.
	// This allows us to handle server errors without blocking the main goroutine.
	serverErrors := make(chan error, 1)

	// Starting the server in a separate goroutine
	go func() {
		log.Printf("Starting server on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	select {
	case err := <-serverErrors:
		log.Fatal("Server error:", err)
	case sig := <-shutdown:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)

			if err := server.Close(); err != nil {
				log.Printf("Server close failed: %v", err)
			}
		}
		log.Printf("Received signal: %s. Shutting down server...", sig)
	}
}

func setupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/shutdown", shutdownHandler)
	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	// This is where you would implement the logic to gracefully shut down your server.
	// For example, you might want to close database connections, stop background jobs, etc.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Shutting down..."))
}
