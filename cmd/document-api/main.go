package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rinatkh/homework_go_18/internal/document"
	"github.com/rinatkh/homework_go_18/internal/grpcapi"
	"github.com/rinatkh/homework_go_18/internal/messaging"
	"github.com/rinatkh/homework_go_18/internal/repository/postgres"
	documentv1 "github.com/rinatkh/homework_go_18/pkg/proto/document/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	repo := postgres.New(pool)
	brokers := strings.Split(env("KAFKA_BROKERS", "localhost:19094"), ",")
	pub := messaging.NewPublisher(brokers, "document.events")
	defer pub.Writer.Close()
	uc := document.NewUseCase(repo, pub)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	mux.HandleFunc("POST /documents", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Title string `json:"title"`
		}
		if json.NewDecoder(http.MaxBytesReader(w, r.Body, 65536)).Decode(&input) != nil {
			http.Error(w, "invalid JSON", 400)
			return
		}
		value, err := uc.Create(r.Context(), input.Title)
		if errors.Is(err, document.ErrEmptyTitle) {
			http.Error(w, "empty title", 400)
			return
		}
		if errors.Is(err, document.ErrNotImplemented) {
			http.Error(w, "TODO-02: saved Document but publish not implemented", 501)
			return
		}
		if err != nil {
			http.Error(w, "create failed", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		_ = json.NewEncoder(w).Encode(struct {
			ID     int64  `json:"id"`
			Title  string `json:"title"`
			Status string `json:"status"`
		}{value.ID, value.Title, value.Status})
	})
	server := &http.Server{Addr: env("HTTP_ADDR", "127.0.0.1:8082"), Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	listener, err := net.Listen("tcp", env("GRPC_ADDR", "127.0.0.1:50052"))
	if err != nil {
		log.Fatal(err)
	}
	rpc := grpc.NewServer()
	documentv1.RegisterDocumentServiceServer(rpc, &grpcapi.Server{UseCase: uc})
	reflection.Register(rpc)
	failed := make(chan error, 3)
	go func() { failed <- server.ListenAndServe() }()
	go func() { failed <- rpc.Serve(listener) }()
	// TODO-03 consumes document.processed in owner service and invokes repo.MarkReady.
	go func() {
		err := messaging.Consume(ctx, brokers, "document.events", "document-api", "document.processed", func(ctx context.Context, e document.Event) error { return repo.MarkReady(ctx, e.DocumentID) })
		if errors.Is(err, document.ErrNotImplemented) {
			log.Print("TODO-03: result consumer not implemented")
			return
		}
		failed <- err
	}()
	log.Print("Document starter ready; TODO routes are intentionally incomplete")
	select {
	case <-ctx.Done():
	case err := <-failed:
		log.Print(err)
	}
	cancel()
	stop, done := context.WithTimeout(context.Background(), 3*time.Second)
	defer done()
	_ = server.Shutdown(stop)
	rpc.Stop()
}
func env(k, f string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return f
}
