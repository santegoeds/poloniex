package event

import (
	"github.com/santegoeds/poloniex/api/decoder"
	channel2 "github.com/santegoeds/poloniex/channel"
	"github.com/santegoeds/poloniex/message"
)

var _ Event = &OrderKilled{}

type OrderKilled struct {
	OrderNumber   int64
	ClientOrderID int64
}

func (ok OrderKilled) ChannelID() int {
	return int(channel2.AccountNotifications)
}

func (ok OrderKilled) Type() message.Type {
	return message.AccountOrderKilled
}

func (ok *OrderKilled) Unmarshal(msg message.Message) error {
	return decoder.Unmarshal(
		msg.Data,
		&ok.OrderNumber,
		&decoder.Int64{Value: &ok.ClientOrderID},
	)
}
