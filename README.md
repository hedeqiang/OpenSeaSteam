# OpenSea Go Client

The OpenSea Go client is a package that provides a convenient way to connect to and interact with the OpenSea API using the Go programming language.

## Installation

To install the OpenSea Go client package, use the `go get` command:

```
go get github.com/hedeqiang/opensea-stream-sdk-go
```

## Usage

Import the `opensea` package into your Go code:
```
import "github.com/hedeqiang/opensea-stream-sdk-go"
```

### Creating a Client

To create a new OpenSea client, use the `NewClient` function:

```go
options := opensea.NewClientOptions()
client := opensea.NewClient(options)
```

You can customize the client options by setting additional properties on the `ClientOptions` struct. The `Token` field is required and should be set to your OpenSea API token.

### Connecting to OpenSea

To establish a connection to OpenSea, use the `Connect` method:

```go
err := client.Connect()
if err != nil {
    // handle error
}
```

### Subscribing to Collections

You can subscribe to collections on OpenSea using the `Subscribe` method:

```go
err := client.Subscribe([]string{"collection1", "collection2"})
if err != nil {
    // handle error
}
```
This will subscribe to updates for the specified collections.

### Handling Events

To handle specific events, you can register event handlers using the `OnItems` method:

```go
client.OnItemListed(func(message opensea.Message) {
    // handle item listed event
})

client.OnItemSold(func(message opensea.Message) {
    // handle item sold event
})
```

You can register multiple event handlers for different event types.

### Unsubscribing from Collections

If you no longer want to receive updates for certain collections, you can unsubscribe using the UnSubscribe method:

```go
err := client.UnSubscribe([]string{"collection1", "collection2"})
if err != nil {
    // handle error
}
```

This will unsubscribe from updates for the specified collections.

### Closing the Connection

To close the connection to OpenSea, use the `Close` method:

```go
err := client.Close()
if err != nil {
    // handle error
}
```

## Example

Here's a complete example that demonstrates how to use the OpenSea Go client:

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hedeqiang/opensea-stream-sdk-go"
)

func main() {
	options := opensea.NewClientOptions().WithToken("<TOKEN>")
	client := opensea.NewClient(options)

	client.OnItemListed(func(message opensea.Message) {
		fmt.Println("Item Listed:", message.Payload)
	})

	client.OnItemSold(func(message opensea.Message) {
		fmt.Println("Item Sold:", message.Payload)
	})

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = client.Subscribe([]string{"collection1", "collection2"})
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Minute)

	err = client.UnSubscribe([]string{"collection1", "collection2"})
	if err != nil {
		log.Fatal(err)
	}

	err = client.Close()
	if err != nil {
		log.Fatal(err)
	}
}
```

## License

This package is licensed under the MIT License. See the LICENSE file for more information.