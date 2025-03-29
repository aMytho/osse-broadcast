package redis

import (
	"context"
	"log"
	"osse-broadcast/internal/messages"
	server "osse-broadcast/internal/server"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func Connect(host string) {
	rdb = redis.NewClient(&redis.Options{
		Addr: host,
	})

	log.Println("Connected to redis on " + host)

	// Start Redis pub/sub listener
	go listenRedis()
}

func listenRedis() {
	ctx := context.Background()
	pubsub := rdb.Subscribe(ctx, "osse_database_private-scan")

	for msg := range pubsub.Channel() {
		log.Println("Received message:", msg.Payload)

		// Parse the message into a Message type
		event, err := messages.GetEventFromMessage(msg.Payload)
		if err != nil {
			println("Received message from Osse that osse-broadcast cannot parse...")
			continue
		}

		// Broadcast to connected SSE clients
		broadcastMessage(event)
	}
}

func broadcastMessage(message messages.OsseEvent) {
	server.Clients <- message
}
