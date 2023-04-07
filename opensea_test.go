package opensea

import (
	"fmt"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	ops := &ClientOptions{
		Token:             "<token>",
		NetWork:           TESTNet,
		HeartbeatInterval: 30 * time.Second,
		RetryInterval:     3 * time.Second,
		MaxRetries:        3,
	}
	client := NewClient(ops)
	err := client.Connect()
	if err != nil {
		fmt.Println("Error connecting client:", err)
		return
	}

	client.WithSlugs([]string{"*"})
	client.Subscribe(ops.Slugs)

	client.OnItemListed(func(message Message) {
		fmt.Println(string(message.Payload))
	})

	client.OnItemCancelled(func(message Message) {
		fmt.Println(string(message.Payload))
	})

	client.OnItemSold(func(message Message) {
		fmt.Println(string(message.Payload))
	})

	client.OnItemTransferred(func(message Message) {
		fmt.Println(string(message.Payload))
	})

	select {}
}
