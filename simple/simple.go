package simple

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/valkey-io/valkey-go"
)

func Call(valkeyServerURL string, channel string, message string) {
	var wg sync.WaitGroup
	wg.Add(1) // Increment the WaitGroup counter

	// Create a new Redis client
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{valkeyServerURL}, // Redis server address
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Subscribe
	go Subscribe(&wg, client)

	time.Sleep(5 * time.Second)
	// Publish a message
	result := client.Do(context.Background(), client.B().Publish().Channel(channel).Message(message).Build())
	if result.Error() != nil {
		panic(result.Error())
	}

	fmt.Println("Publish result:", result.String())
	// Wait for all goroutines to finish
	wg.Wait()
}

func Subscribe(wg *sync.WaitGroup, client valkey.Client) {
	defer wg.Done() // Notify the WaitGroup that this goroutine is done
	// Subscribe to channels
	err := client.Receive(context.Background(), client.B().Subscribe().Channel("ch1", "ch2").Build(), func(msg valkey.PubSubMessage) {
		// Handle the message
		fmt.Printf("Message received from channel '%s': %s\n", msg.Channel, msg.Message)
	})
	if err != nil {
		panic(err)
	}
}
