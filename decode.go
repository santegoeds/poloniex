package poloniex

import (
	"fmt"

	"github.com/santegoeds/poloniex/errors"
	"github.com/santegoeds/poloniex/event"
	"github.com/santegoeds/poloniex/message"
)

func Decode(m message.Message) (event.Event, error) {
	var e event.Event
	switch m.Type {
	case message.AccountBalanceUpdate:
		e = &event.BalanceUpdate{}
	case message.AccountLimitOrder:
		e = &event.LimitOrder{}
	case message.AccountOrderKilled:
		e = &event.OrderKilled{}
	case message.AccountOrderPending:
		e = &event.OrderPending{}
	case message.AccountOrderUpdate:
		e = &event.OrderUpdate{}
	case message.AccountMarginPositionUpdate:
		e = &event.MarginPositionUpdate{}
	case message.AccountTrade:
		e = &event.AccountTrade{}
	case message.Book:
		e = &event.Book{}
	case message.BookTrade:
		e = &event.BookTrade{}
	case message.BookUpdate:
		e = &event.BookUpdate{}
	case message.ExchangeVolume24h:
		e = &event.ExchangeVolume24h{}
	case message.TickerData:
		e = &event.TickerData{}
	default:
		return nil, fmt.Errorf("%w: unexpected message type %d", errors.ErrBadResponse, m.Type)
	}

	if err := e.Unmarshal(m); err != nil {
		return nil, err
	}
	return e, nil
}
