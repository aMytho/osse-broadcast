package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"osse-broadcast/internal/messages"
	"time"

	"github.com/tmaxmax/go-sse"
)

// SSE connections
var Clients = make(chan messages.OsseEvent)

func Start(host string, allowOrigin string) {
	sseHandler := createSseSetup()

	mux := http.NewServeMux()
	mux.Handle("/sse", sseHandler)

	httpServer := &http.Server{
		Addr:              host,
		Handler:           cors(mux, allowOrigin),
		ReadHeaderTimeout: time.Second * 10,
	}

	httpServer.RegisterOnShutdown(func() {
		e := &sse.Message{Type: sse.Type("close")}
		// Adding data is necessary because spec-compliant clients
		// do not dispatch events without data.
		e.AppendData("bye")
		// Broadcast a close message so clients can gracefully disconnect.
		_ = sseHandler.Publish(e)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		// We use a context with a timeout so the program doesn't wait indefinitely
		// for connections to terminate. There may be misbehaving connections
		// which may hang for an unknown timespan, so we just stop waiting on Shutdown
		// after a certain duration.
		_ = sseHandler.Shutdown(ctx)
	})

	// Listen for redis messages
	go func() {
		for event := range Clients {
			eventJson, err := messages.GetJsonOfEvent(event)
			if err != nil {
				continue
			}

			message := &sse.Message{}
			message.AppendData(eventJson)

			// Set the event name (client listens for this)
			eventName, err := sse.NewType(event.GetType())
			if err != nil {
				continue
			}
			message.Type = eventName

			sseHandler.Publish(message, messages.AllTopics...)
			println("Sent Message to client")
		}
	}()

	log.Println("Osse Broadcast running on " + host)
	runServer(httpServer)
}

func createSseSetup() *sse.Server {
	return &sse.Server{
		Provider: &sse.Joe{},
		OnSession: func(s *sse.Session) (sse.Subscription, bool) {
			// Get the user ID and token
			userID := s.Req.URL.Query().Get("id")
			token := s.Req.URL.Query().Get("token")

			log.Println("User attempted to connect.")

			// Validate the userID and token
			if !validateUserToken(userID, token) {
				return sse.Subscription{}, false
			}

			return sse.Subscription{
				Client: s,
				Topics: messages.AllTopics,
			}, true
		},
	}
}

func runServer(s *http.Server) error {
	shutdownError := make(chan error)

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return <-shutdownError
}
