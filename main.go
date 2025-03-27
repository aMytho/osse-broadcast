package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"osse-broadcast/internal/config"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func main() {
	config := config.GetOsseConfig()

	// Connect to Redis/Valkey
	rdb = redis.NewClient(&redis.Options{
		Addr: config.RedisHost,
		// Addr: "localhost:6379",
	})

	// Start Redis pub/sub listener
	go listenRedis()

	// Start HTTP server with SSE route
	http.HandleFunc("/sse", sseHandler)
	log.Println("Osee Broadcast running on " + config.HttpHost)
	log.Fatal(http.ListenAndServe(config.HttpHost, nil))
}

func listenRedis() {
	ctx := context.Background()
	pubsub := rdb.Subscribe(ctx, "example")

	for msg := range pubsub.Channel() {
		log.Println("Received message:", msg.Payload)
		// Broadcast to connected SSE clients
		broadcastMessage(msg.Payload)
	}
}

// SSE connections
var clients = make(map[chan string]struct{})

func sseHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create a new channel for this client
	messageChan := make(chan string)
	clients[messageChan] = struct{}{}
	defer delete(clients, messageChan)

	for msg := range messageChan {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		flusher.Flush()
	}
}

func broadcastMessage(message string) {
	for client := range clients {
		client <- message
	}
}
