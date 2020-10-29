package message

import (
	"encoding/json"
)

type Type int

const (
	_ Type = iota
	// ChannelAccountNotifications
	AccountOrderPending
	AccountBalanceUpdate
	AccountLimitOrder
	AccountOrderUpdate
	AccountMarginPositionUpdate
	AccountTrade
	AccountOrderKilled
	// ChannelTickerData
	TickerData
	// ChannelExchangeVolume24h
	ExchangeVolume24h
	// PriceAggregatedBook
	Book
	BookUpdate
	BookTrade
)

type Message struct {
	ChannelID int
	Type      Type
	Data      []json.RawMessage
}
