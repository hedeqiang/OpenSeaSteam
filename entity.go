package opensea

import "time"

type Collection struct {
	Slug string `json:"slug"`
}
type Chain struct {
	Name string `json:"name"`
}
type Metadata struct {
	AnimationURL interface{} `json:"animation_url"`
	ImageURL     string      `json:"image_url"`
	MetadataURL  string      `json:"metadata_url"`
	Name         string      `json:"name"`
}
type Item struct {
	Chain     Chain    `json:"chain"`
	Metadata  Metadata `json:"metadata"`
	NftID     string   `json:"nft_id"`
	Permalink string   `json:"permalink"`
}
type Maker struct {
	Address string `json:"address"`
}
type PaymentToken struct {
	Address  string `json:"address"`
	Decimals int    `json:"decimals"`
	EthPrice string `json:"eth_price"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	UsdPrice string `json:"usd_price"`
}
type Consideration struct {
	EndAmount            string `json:"endAmount"`
	IdentifierOrCriteria string `json:"identifierOrCriteria"`
	ItemType             int    `json:"itemType"`
	Recipient            string `json:"recipient"`
	StartAmount          string `json:"startAmount"`
	Token                string `json:"token"`
}
type Offer struct {
	EndAmount            string `json:"endAmount"`
	IdentifierOrCriteria string `json:"identifierOrCriteria"`
	ItemType             int    `json:"itemType"`
	StartAmount          string `json:"startAmount"`
	Token                string `json:"token"`
}
type Parameters struct {
	ConduitKey    string          `json:"conduitKey"`
	Consideration []Consideration `json:"consideration"`
	//Counter                         int64           `json:"counter"`  //此参数类型不保证，会返回 int 和 string 两种类型，导致无法正确接收
	EndTime                         string  `json:"endTime"`
	Offer                           []Offer `json:"offer"`
	Offerer                         string  `json:"offerer"`
	OrderType                       int     `json:"orderType"`
	Salt                            string  `json:"salt"`
	StartTime                       string  `json:"startTime"`
	TotalOriginalConsiderationItems int     `json:"totalOriginalConsiderationItems"`
	Zone                            string  `json:"zone"`
	ZoneHash                        string  `json:"zoneHash"`
}

type ProtocolData struct {
	Parameters Parameters `json:"parameters"`
	Signature  string     `json:"signature"`
}

type Transaction struct {
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp"`
}

type Taker struct {
	Address string `json:"address"`
}

type FromAccount struct {
	Address string `json:"address"`
}

type ToAccount struct {
	Address string `json:"address"`
}

type ItemCancelledEventPayload struct {
	BasePrice       string       `json:"base_price"`
	Collection      Collection   `json:"collection"`
	EventTimestamp  time.Time    `json:"event_timestamp"`
	ExpirationDate  time.Time    `json:"expiration_date"`
	IsPrivate       bool         `json:"is_private"`
	Item            Item         `json:"item"`
	ListingDate     time.Time    `json:"listing_date"`
	ListingType     interface{}  `json:"listing_type"`
	Maker           Maker        `json:"maker"`
	OrderHash       string       `json:"order_hash"`
	PaymentToken    PaymentToken `json:"payment_token"`
	ProtocolAddress string       `json:"protocol_address"`
	ProtocolData    ProtocolData `json:"protocol_data"`
	Quantity        int          `json:"quantity"`
	Taker           interface{}  `json:"taker"`
	Transaction     Transaction  `json:"transaction"`
}

type ItemListedEventPayload struct {
	BasePrice       string       `json:"base_price"`
	Collection      Collection   `json:"collection"`
	EventTimestamp  time.Time    `json:"event_timestamp"`
	ExpirationDate  time.Time    `json:"expiration_date"`
	IsPrivate       bool         `json:"is_private"`
	Item            Item         `json:"item"`
	ListingDate     time.Time    `json:"listing_date"`
	ListingType     string       `json:"listing_type"`
	Maker           Maker        `json:"maker"`
	OrderHash       string       `json:"order_hash"`
	PaymentToken    PaymentToken `json:"payment_token"`
	ProtocolAddress string       `json:"protocol_address"`
	ProtocolData    ProtocolData `json:"protocol_data"`
	Quantity        int64        `json:"quantity"`
	Taker           interface{}  `json:"taker"`
}

type ItemSoldEventPayload struct {
	ClosingDate     time.Time    `json:"closing_date"`
	Collection      Collection   `json:"collection"`
	EventTimestamp  time.Time    `json:"event_timestamp"`
	IsPrivate       bool         `json:"is_private"`
	Item            Item         `json:"item"`
	ListingType     string       `json:"listing_type"`
	Maker           Maker        `json:"maker"`
	OrderHash       string       `json:"order_hash"`
	PaymentToken    PaymentToken `json:"payment_token"`
	ProtocolAddress string       `json:"protocol_address"`
	ProtocolData    ProtocolData `json:"protocol_data"`
	Quantity        int64        `json:"quantity"`
	SalePrice       string       `json:"sale_price"`
	Taker           Taker        `json:"taker"`
	Transaction     Transaction  `json:"transaction"`
}

type ItemTransferredEventPayload struct {
	Collection     Collection  `json:"collection"`
	EventTimestamp time.Time   `json:"event_timestamp"`
	FromAccount    FromAccount `json:"from_account"`
	Item           Item        `json:"item"`
	Quantity       int64       `json:"quantity"`
	ToAccount      ToAccount   `json:"to_account"`
	Transaction    Transaction `json:"transaction"`
}

type BaseStreamMessage struct {
	EventType string `json:"event_type"`
	SentAt    string `json:"sent_at"`
}

type ItemCancelledEvent struct {
	BaseStreamMessage
	Payload ItemCancelledEventPayload `json:"payload"`
}

type ItemListedEvent struct {
	BaseStreamMessage
	Payload ItemListedEventPayload `json:"payload"`
}

type ItemSoldEvent struct {
	BaseStreamMessage
	Payload ItemSoldEventPayload `json:"payload"`
}

type ItemTransferredEvent struct {
	BaseStreamMessage
	Payload ItemTransferredEventPayload `json:"payload"`
}

type EventPayload struct {
	Event   string `json:"event"`
	Payload struct {
		Response struct {
		} `json:"response"`
		Status string `json:"status"`
	} `json:"payload"`
	Ref   string `json:"ref"`
	Topic string `json:"topic"`
}

type Listings struct {
	Listings []struct {
		OrderHash string `json:"order_hash"`
		Chain     string `json:"chain"`
		Type      string `json:"type"`
		Price     struct {
			Current struct {
				Value    string `json:"value"`
				Decimals int    `json:"decimals"`
				Currency string `json:"currency"`
			} `json:"current"`
		} `json:"price"`
		ProtocolData struct {
			Parameters struct {
				Offerer string `json:"offerer"`
				Offer   []struct {
					ItemType             int    `json:"itemType"`
					Token                string `json:"token"`
					IdentifierOrCriteria string `json:"identifierOrCriteria"`
					StartAmount          string `json:"startAmount"`
					EndAmount            string `json:"endAmount"`
				} `json:"offer"`
				Consideration []struct {
					ItemType             int    `json:"itemType"`
					Token                string `json:"token"`
					IdentifierOrCriteria string `json:"identifierOrCriteria"`
					StartAmount          string `json:"startAmount"`
					EndAmount            string `json:"endAmount"`
					Recipient            string `json:"recipient"`
				} `json:"consideration"`
				StartTime                       string `json:"startTime"`
				EndTime                         string `json:"endTime"`
				OrderType                       int    `json:"orderType"`
				Zone                            string `json:"zone"`
				ZoneHash                        string `json:"zoneHash"`
				Salt                            string `json:"salt"`
				ConduitKey                      string `json:"conduitKey"`
				TotalOriginalConsiderationItems int    `json:"totalOriginalConsiderationItems"`
				Counter                         int    `json:"counter"`
			} `json:"parameters"`
			Signature string `json:"signature"`
		} `json:"protocol_data"`
		ProtocolAddress string `json:"protocol_address"`
	} `json:"listings"`
	Next string `json:"next"`
}
