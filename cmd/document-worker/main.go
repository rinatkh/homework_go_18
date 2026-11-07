package main

import (
	"context"
	goredis "github.com/redis/go-redis/v9"
	"github.com/rinatkh/homework_go_18/internal/messaging"
	redisrepo "github.com/rinatkh/homework_go_18/internal/repository/redis"
	"github.com/rinatkh/homework_go_18/internal/worker"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	options, err := goredis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	client := goredis.NewClient(options)
	defer client.Close()
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	pub := messaging.NewPublisher(brokers, "document.events")
	defer pub.Writer.Close()
	processor := &worker.Processor{Publisher: pub, Receipts: &redisrepo.ReceiptStore{Client: client}}
	if err := messaging.Consume(ctx, brokers, "document.events", "document-workers", "document.created", processor.Handle); err != nil {
		log.Fatal(err)
	}
}
