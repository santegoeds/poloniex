package ws

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/santegoeds/poloniex/errors"
	"sync"
)

const (
	Endpoint = "wss://api2.poloniex.com"
)

var (
	DefaultDialer = websocket.DefaultDialer
)

type Client struct {
	Endpoint      string
	Dialer        *websocket.Dialer
	key           string
	secret        string
	lck           sync.Mutex
	subscriptions []*Subscription
}

func New(key, secret string) *Client {
	return &Client{
		Endpoint:      Endpoint,
		Dialer:        DefaultDialer,
		key:           key,
		secret:        secret,
		subscriptions: make([]*Subscription, 0),
	}
}

func (c *Client) Subscribe(channel interface{}) *Subscription {
	c.lck.Lock()
	defer c.lck.Unlock()

	subscription := NewSubscription(c, channel)
	c.subscriptions = append(c.subscriptions, subscription)
	return subscription
}

func (c *Client) Dial(ctx context.Context) (*websocket.Conn, error) {
	conn, rsp, err := c.Dialer.DialContext(ctx, c.Endpoint, nil)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode/100 > 2 {
		_ = conn.Close()
		return nil, fmt.Errorf("%w: invalid status code %s", errors.ErrBadResponse, rsp.Status)
	}
	return conn, nil
}

func (c *Client) Close() {
	c.lck.Lock()
	defer c.lck.Unlock()

	for _, sub := range c.subscriptions {
		_ = sub.Close()
	}
	c.subscriptions = c.subscriptions[:0]
	return
}
