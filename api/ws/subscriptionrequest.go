package ws

import (
	"github.com/santegoeds/poloniex/channel"
)

type SubscriptionRequest struct {
	client *Client
}

func NewSubscriptionRequest(client *Client) *SubscriptionRequest {
	return &SubscriptionRequest{
		client: client,
	}
}

func (r *SubscriptionRequest) CurrencyPair(currencyPair string) *Subscription {
	return r.client.Subscribe(currencyPair)
}

func (r *SubscriptionRequest) CurrencyPairID(currencyPairID int) *Subscription {
	return r.client.Subscribe(currencyPairID)
}
func (r *SubscriptionRequest) Channel(channel channel.Channel) *Subscription {
	return r.client.Subscribe(channel)
}
