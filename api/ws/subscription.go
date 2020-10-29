package ws

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/santegoeds/poloniex/message"
)

var (
	DefaultTimeoutDuration = time.Millisecond * 1500
	ErrTimeout             = errors.New("heartbeat timeout")
)

type Subscription struct {
	TimeoutDuration time.Duration
	client          *Client
	ctx             context.Context
	lck             sync.Mutex
	wg              sync.WaitGroup
	channel         interface{}
	msgC            chan message.Message
	errC            chan error
	conn            *websocket.Conn
	cancel          context.CancelFunc
	ticker          *time.Ticker
}

func NewSubscription(client *Client, channel interface{}) *Subscription {
	return &Subscription{
		TimeoutDuration: DefaultTimeoutDuration,
		client:          client,
		channel:         channel,
		msgC:            make(chan message.Message, 1),
		errC:            make(chan error, 1),
	}
}

func (s *Subscription) MsgC() <-chan message.Message {
	return s.msgC
}

func (s *Subscription) ErrC() <-chan error {
	return s.errC
}

func (s *Subscription) Channels() (<-chan message.Message, <-chan error) {
	return s.msgC, s.errC
}

func (s *Subscription) Do(ctx context.Context) (*Subscription, error) {
	s.lck.Lock()
	defer s.lck.Unlock()

	s.ctx, s.cancel = context.WithCancel(ctx)

	s.ticker = time.NewTicker(s.TimeoutDuration)
	err := s.dailWebSocket()
	if err != nil {
		return nil, err
	}

	s.wg.Add(1)
	go func() {
		defer s.cancel()
		defer s.wg.Done()
		s.receiveAndReconnect()
	}()

	return s, nil
}

func (s *Subscription) Close() error {
	if s.cancel != nil {
		s.cancel()
	}

	s.wg.Wait()
	return nil
}

func (s *Subscription) HandleMessage(msg message.Message) {
	s.sendMessage(msg)
}

func (s *Subscription) HandleHeartbeat() {
	s.ticker.Reset(s.TimeoutDuration)
}

func (s *Subscription) dailWebSocket() error {
	s.ticker.Reset(s.TimeoutDuration)
	conn, err := s.client.Dial(s.ctx)
	if err != nil {
		return err
	}
	s.conn = conn

	return s.subscribe()
}

func (s *Subscription) subscribe() error {
	msg := map[string]interface{}{
		"command": "subscribe",
		"channel": s.channel,
	}
	return s.conn.WriteJSON(msg)
}

func (s *Subscription) receiveAndReconnect() {
	for !s.isDone() {
		err := s.receiveMessages()
		if err != nil {
			log.Printf("failed to receive messages from websocket connection: %s", err)
			_ = s.dailWebSocket()
		}
	}

	if s.conn != nil {
		_ = s.conn.Close()
		s.conn = nil
	}
	s.ticker.Stop()

	close(s.errC)
	close(s.msgC)
}

func (s *Subscription) receiveMessages() error {
	decoder := NewDecoder(s)
	buffer := make([]json.RawMessage, 0, 3)
	for {
		select {
		case <-s.ctx.Done():
			return nil

		case <-s.ticker.C:
			return ErrTimeout

		default:
			if err := s.conn.ReadJSON(&buffer); err != nil {
				return err
			}
			if err := decoder.Decode(buffer); err != nil {
				s.sendError(err)
			}
		}
	}
}

func (s *Subscription) redialWebSocket() {
	s.lck.Lock()
	defer s.lck.Unlock()

	if s.conn != nil {
		_ = s.conn.Close()
		s.conn = nil
	}
	_ = s.dailWebSocket()
}

func (s *Subscription) isDone() bool {
	select {
	case <-s.ctx.Done():
		return true
	default:
	}
	return false
}

func (s *Subscription) sendError(err error) {
	select {
	case s.errC <- err:
	case <-s.ctx.Done():
	}
}

func (s *Subscription) sendMessage(msg message.Message) {
	s.ticker.Reset(s.TimeoutDuration)
	select {
	case s.msgC <- msg:
	case <-s.ctx.Done():
	}
}
