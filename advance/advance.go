package advance

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/valkey-io/valkey-go"
)

// PubSubService handles publishing and subscribing to Redis channels
type PubSubService struct {
	client valkey.Client
}

// NewPubSubService creates a new PubSubService with a Redis client
func NewPubSubService(redisAddress string) (*PubSubService, error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{redisAddress}, // Redis server address
	})
	if err != nil {
		return nil, err
	}

	return &PubSubService{client: client}, nil
}

// Subscribe listens for messages on the specified channels
func (ps *PubSubService) Subscribe(wg *sync.WaitGroup, channels ...string) {
	defer wg.Done() // Notify the WaitGroup that this goroutine is done

	err := ps.client.Receive(context.Background(), ps.client.B().Subscribe().Channel(channels...).Build(), func(msg valkey.PubSubMessage) {
		// Handle the message
		fmt.Printf("Message received from channel '%s': %s\n", msg.Channel, msg.Message)
	})

	if err != nil {
		panic(err)
	}
}

// Publish sends a message to the specified channel
func (ps *PubSubService) Publish(channel, message string) error {
	result := ps.client.Do(context.Background(), ps.client.B().Publish().Channel(channel).Message(message).Build())
	if result.Error() != nil {
		return result.Error()
	}
	fmt.Println("Publish result:", result.String())
	return nil
}

// Close shuts down the PubSubService by closing the Redis client
func (ps *PubSubService) Close() {
	ps.client.Close()
}

func Call(valkeyServerURL string) {
	// Initialize the PubSubService
	pubSubService, err := NewPubSubService(valkeyServerURL)
	if err != nil {
		panic(err)
	}
	defer pubSubService.Close()

	var wg sync.WaitGroup
	wg.Add(1) // Increment the WaitGroup counter

	// Start the subscriber in a separate goroutine
	go pubSubService.Subscribe(&wg, "ch1", "ch2")

	// Example HTTP server setup
	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		// Example of publishing a message through an HTTP request
		message := r.URL.Query().Get("message")
		if message == "" {
			http.Error(w, "Message is required", http.StatusBadRequest)
			return
		}

		err := pubSubService.Publish("ch1", message)
		if err != nil {
			http.Error(w, "Failed to publish message", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Message published: %s\n", message)
	})

	// Start the HTTP server
	go func() {
		fmt.Println("HTTP server started on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()
}
