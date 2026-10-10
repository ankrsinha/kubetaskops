package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ankrsinha/kubetaskops/backend/internal/traditional"
	"github.com/ankrsinha/kubetaskops/backend/pkg/store"
)

// corsMiddleware attaches standard CORS headers for local web development
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("🚀 Starting KubeTaskOps API Server...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Database Pool
	dbCfg := store.DefaultConfig()
	if host := os.Getenv("DB_HOST"); host != "" {
		dbCfg.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		dbCfg.Port = port
	}

	pool, err := store.NewPostgresPool(ctx, dbCfg)
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	defer pool.Close()
	log.Printf("✅ Connected to PostgreSQL at %s:%s/%s", dbCfg.Host, dbCfg.Port, dbCfg.DBName)

	// Build Router & Register Endpoints
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Register Traditional 3-Tier CRUD Handlers (Phase 1)
	tradRepo := traditional.NewRepository(pool)
	tradService := traditional.NewService(tradRepo)
	tradHandler := traditional.NewHandler(tradService)
	tradHandler.RegisterRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown handling
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🌐 KubeTaskOps HTTP Server listening on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-shutdownChan
	log.Println("🛑 Shutting down server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
	log.Println("👋 Server stopped.")
}
