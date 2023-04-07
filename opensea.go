package opensea

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

const (
	defaultEndpoint      = "https://api.opensea.io"
	defaultRetryInterval = 5 * time.Second
	defaultHeartbeat     = 30 * time.Second
	defaultMaxRetries    = 3
	defaultNetwork       = TESTNet
)

type NetWork string

const (
	MainNet NetWork = "Mainnet"
	TESTNet NetWork = "Testnet"
)

type payload struct {
	Topic   string      `json:"topic"`
	Event   string      `json:"event"`
	Ref     string      `json:"ref"`
	Payload interface{} `json:"payload"`
}

type Message struct {
	Topic   string          `json:"topic"`
	Event   string          `json:"event"`
	Ref     string          `json:"ref"`
	Payload json.RawMessage `json:"payload"`
}

type EventType string

const (
	ItemMetadataUpdated EventType = "item_metadata_updated"
	ItemListed          EventType = "item_listed"
	ItemSold            EventType = "item_sold"
	ItemTransferred     EventType = "item_transferred"
	ItemReceivedOffer   EventType = "item_received_offer"
	ItemReceivedBid     EventType = "item_received_bid"
	ItemCancelled       EventType = "item_cancelled"
	CollectionOffer     EventType = "collection_offer"
	TraitOffer          EventType = "trait_offer"
)

type EventHandleFunc func(message Message)

type Client struct {
	options         *ClientOptions
	messageHandlers map[EventType]EventHandleFunc
	mu              sync.Mutex
	conn            *websocket.Conn
}

type ClientOptions struct {
	Token             string
	Slugs             []string
	NetWork           NetWork
	HeartbeatInterval time.Duration
	RetryInterval     time.Duration
	MaxRetries        int
	reconnecting      bool
}

func NewClientOptions() *ClientOptions {
	return &ClientOptions{
		Token:             "<TOKEN>",
		NetWork:           TESTNet,
		HeartbeatInterval: defaultHeartbeat,
		RetryInterval:     defaultRetryInterval,
		MaxRetries:        defaultMaxRetries,
	}
}

var client *Client

// GetClient 获取 client
func GetClient() *Client {
	//fmt.Println(fmt.Sprintf("%p", &client))
	return client
}

func NewClient(options *ClientOptions) *Client {
	c := &Client{
		messageHandlers: make(map[EventType]EventHandleFunc),
		mu:              sync.Mutex{},
	}
	if options == nil {
		options = NewClientOptions()
	}
	c.WithOptions(options)

	client = c
	return client
}

func (c *Client) WithOptions(options *ClientOptions) *Client {
	c.options = options
	return c
}

func (c *Client) WithSlugs(slugs []string) *Client {
	c.options.Slugs = slugs
	return c
}

func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var err error
	var retryCount int

	ticker := time.NewTicker(c.options.RetryInterval)
	defer ticker.Stop()
	url := ""
	if c.options.NetWork == MainNet {
		url = "wss://stream.openseabeta.com/socket/websocket?token=" + c.options.Token
	} else {
		url = "wss://testnets-stream.openseabeta.com/socket/websocket?token=" + c.options.Token
	}

	// 尝试连接
	for {
		if c.conn != nil {
			c.Close()
		}
		c.conn, _, err = websocket.DefaultDialer.Dial(url, nil)
		if err == nil {
			// 连接成功
			break
		}

		// 连接失败，等待一段时间后重试
		retryCount++
		if retryCount >= c.options.MaxRetries {
			log.Println(fmt.Sprintf("connect failed after %d retries: %v ,url:%s", c.options.MaxRetries, err, url))
			return fmt.Errorf("connect failed after %d retries: %v", c.options.MaxRetries, err)
		}

		<-ticker.C
	}

	// 连接成功，启动心跳检测和读取消息协程
	go c.heartbeat()
	go c.readMessages()

	return nil
}

func (c *Client) Conn() *websocket.Conn {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn
}

func (c *Client) Close() error {
	if c.conn == nil {
		log.Println("websocket is not connected")
		return nil
	}
	if err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		log.Println(fmt.Sprintf("error closing websocket: %s", err))
	}
	return c.conn.Close()
}

func (c *Client) heartbeat() {
	ticker := time.NewTicker(c.options.HeartbeatInterval)
	for {
		select {
		case <-ticker.C:
			for _, slug := range c.options.Slugs {
				// 发送心跳消息
				//logx.Logger.Info("send heartbeat...")
				heartbeat := payload{
					Topic: "phoenix",
					Event: "heartbeat",
					Ref:   slug,
				}
				if c.conn == nil {
					log.Println("链接已关闭，无法发送心跳")
					return
				}
				c.mu.Lock()
				if err := c.conn.WriteJSON(heartbeat); err != nil {
					c.mu.Unlock()
					log.Printf("error sending heartbeat message: %s", err)
					return
				}
				c.mu.Unlock()
			}

		}
	}
}

func (c *Client) readMessages() {
	for {
		// Set a deadline for reading a message from the websocket connection
		if c.conn == nil {
			log.Println("websocket is not connected, cannot read messages")
			return
		}
		//c.conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// Check if the error is due to a closed connection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(fmt.Sprintf("Websocket connection closed unexpectedly: %s", err))
			}
			c.mu.Lock()
			fmt.Println(err.Error())
			c.conn = nil
			c.mu.Unlock()
			return
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println(fmt.Sprintf("error decoding message: %s", err))
			continue
		}
		if msg.Event == "item_received_bid" || msg.Event == "trait_offer" || msg.Event == "collection_offer" {
			continue
		}
		//if msg.Event == "item_transferred" {
		//	c.conn = nil
		//}
		var msgEvent EventPayload
		json.Unmarshal(message, &msgEvent)
		//fmt.Println(string(message))
		if msg.Event == "phx_reply" && msgEvent.Payload.Status != "ok" {
			// 重新订阅
			c.Subscribe([]string{msgEvent.Ref})
		}

		if msg.Event == "phx_close" {
			//split := strings.Split(msgEvent.Topic, ":")
			//c.Subscribe([]string{split[1]})
			log.Println(fmt.Sprintf("%s 订阅失败 closed", msgEvent.Topic))
		}

		c.mu.Lock()
		handler, ok := c.messageHandlers[EventType(msg.Event)]
		c.mu.Unlock()

		if ok {
			handler(msg)
		}
	}
}

// Subscribe to collection
func (c *Client) Subscribe(collections []string) error {
	if collections == nil {
		return fmt.Errorf("collections is nil")
	}
	for _, collection := range collections {
		topic := fmt.Sprintf("collection:%s", collection)
		subMsg := payload{
			Topic: topic,
			Event: "phx_join",
			Ref:   collection,
		}
		// Lock the mutex before writing to the WebSocket connection
		c.mu.Lock()
		if err := c.conn.WriteJSON(subMsg); err != nil {
			c.mu.Unlock()
			log.Println(fmt.Sprintf("Error sending subscription message: %s", err))
			return err
		}
		c.mu.Unlock()
		log.Println(fmt.Sprintf("Successfully joined channel %s", topic))
	}

	return nil
}

// UnSubscribe  collection
func (c *Client) UnSubscribe(collections []string) error {
	if collections == nil {
		return fmt.Errorf("collections is nil")
	}
	for _, collection := range collections {
		topic := fmt.Sprintf("collection:%s", collection)
		subMsg := payload{
			Topic: topic,
			Event: "phx_leave",
			Ref:   collection,
		}
		c.mu.Lock()
		if err := c.conn.WriteJSON(subMsg); err != nil {
			log.Println(fmt.Sprintf("Error sending unsubscription message: %s", err))
			return err
		}
		c.mu.Unlock()
		log.Println(fmt.Sprintf("Successfully unsubscribe channel %s", topic))
	}

	return nil
}

func (c *Client) OnItemCancelled(handleFunc EventHandleFunc) {
	c.OnItems([]EventType{ItemCancelled}, handleFunc)
}

func (c *Client) OnItemListed(handleFunc EventHandleFunc) {
	c.OnItems([]EventType{ItemListed}, handleFunc)

}

func (c *Client) OnItemSold(handleFunc EventHandleFunc) {
	c.OnItems([]EventType{ItemSold}, handleFunc)
}

func (c *Client) OnItemTransferred(handleFunc EventHandleFunc) {
	c.OnItems([]EventType{ItemTransferred}, handleFunc)

}

func (c *Client) OnItemMetadataUpdated(handleFunc EventHandleFunc) {
	c.OnItems([]EventType{ItemMetadataUpdated}, handleFunc)
}

func (c *Client) OnItemReceivedOffer(handleFunc EventHandleFunc) {
	c.OnItems([]EventType{ItemReceivedOffer}, handleFunc)
}

func (c *Client) OnItemReceivedBid(handleFunc EventHandleFunc) {
	c.OnItems([]EventType{ItemReceivedBid}, handleFunc)
}

func (c *Client) OnCollectionOffer(handleFunc EventHandleFunc) {
	c.OnItems([]EventType{CollectionOffer}, handleFunc)
}

func (c *Client) OnTraitOffer(handleFunc EventHandleFunc) {
	c.OnItems([]EventType{TraitOffer}, handleFunc)
}

func (c *Client) OnItems(eventTypes []EventType, handleFunc EventHandleFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, eventType := range eventTypes {
		c.messageHandlers[eventType] = handleFunc
	}
}
