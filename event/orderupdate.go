package event

import (
	"github.com/santegoeds/poloniex/api/decoder"
	channel2 "github.com/santegoeds/poloniex/channel"
	"github.com/santegoeds/poloniex/message"
)

const (
	orderTypeFilled    = "f"
	orderTypeSelfTrade = "s"
	orderTypeCancelled = "c"
)

var _ Event = &OrderUpdate{}

type OrderUpdate struct {
	OrderNumber   int64
	NewAmount     float64
	OrderType     string
	ClientOrderID int64
}

func (ou OrderUpdate) ChannelID() int {
	return int(channel2.AccountNotifications)
}

func (ou OrderUpdate) Type() message.Type {
	return message.AccountOrderUpdate
}

func (ou OrderUpdate) IsFill() bool {
	return ou.OrderType == orderTypeFilled
}

func (ou OrderUpdate) IsSelfTrade() bool {
	return ou.OrderType == orderTypeSelfTrade
}

func (ou OrderUpdate) IsCancelled() bool {
	return ou.OrderType == orderTypeCancelled
}

func (ou *OrderUpdate) Unmarshal(msg message.Message) error {
	type decF64 = decoder.Float64
	type decI64 = decoder.Int64

	return decoder.Unmarshal(
		msg.Data,
		&ou.OrderNumber,
		&decF64{Value: &ou.NewAmount},
		&ou.OrderType,
		&decI64{Value: &ou.ClientOrderID},
	)
}
