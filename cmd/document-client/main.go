package main

import (
	"context"
	"github.com/rinatkh/homework_go_18/internal/automation"
	documentv1 "github.com/rinatkh/homework_go_18/pkg/proto/document/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	conn, err := grpc.NewClient(os.Getenv("API_GRPC_ADDR"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	id, err := strconv.ParseInt(os.Getenv("DOCUMENT_ID"), 10, 64)
	if err != nil {
		log.Fatal("positive DOCUMENT_ID required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	value, err := automation.Get(ctx, documentv1.NewDocumentServiceClient(conn), id)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("document %d status=%s", value.Id, value.Status)
}
