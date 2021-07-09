package event

import (
	"time"

	"github.com/santegoeds/poloniex/api/decoder"
	channel2 "github.com/santegoeds/poloniex/channel"
	"github.com/santegoeds/poloniex/message"
)

var _ Event = &OrderPending{}

type OrderPending struct {
	OrderNumber    int64
	CurrencyPairID int
	Rate           float64
	Amount         float64
	OrderType      string
	ClientOrderID  int64
	Time           time.Time
}

func (op OrderPending) ChannelID() int {
	return int(channel2.AccountNotifications)
}

func (op OrderPending) Type() message.Type {
	return message.AccountOrderPending
}

func (op OrderPending) IsSell() bool {
	return op.OrderType == "0"
}

func (op OrderPending) IsBuy() bool {
	return op.OrderType == "1"
}

func (op *OrderPending) Unmarshal(msg message.Message) error {
	type decF64 = decoder.Float64
	type decI64 = decoder.Int64
	type decEpochMs = decoder.EpochMs

	return decoder.Unmarshal(
		msg.Data,
		&op.OrderNumber,
		&op.CurrencyPairID,
		&decF64{Value: &op.Rate},
		&decF64{Value: &op.Amount},
		&op.OrderType,
		&decI64{Value: &op.ClientOrderID},
		&decEpochMs{Value: &op.Time},
	)
}
