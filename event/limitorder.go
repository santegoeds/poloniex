package event

import (
	channel2 "github.com/santegoeds/poloniex/channel"
	"time"

	"github.com/santegoeds/poloniex/api/decoder"
	"github.com/santegoeds/poloniex/message"
)

var _ Event = &LimitOrder{}

type LimitOrder struct {
	CurrencyPairID int
	OrderNumber    int64
	OrderType      int
	Rate           float64
	Amount         float64
	Date           string
	OriginalAmount float64
	ClientOrderID  int64
}

func (lo *LimitOrder) ChannelID() int {
	return int(channel2.AccountNotifications)
}

func (lo *LimitOrder) Type() message.Type {
	return message.AccountLimitOrder
}

func (lo *LimitOrder) IsSell() bool {
	return lo.OrderType == 0
}

func (lo *LimitOrder) IsBuy() bool {
	return lo.OrderType == 1
}

func (lo *LimitOrder) DateTime() time.Time {
	ts, _ := time.Parse(timeFormat, lo.Date)
	return ts
}

func (lo *LimitOrder) Unmarshal(msg message.Message) error {
	type decF64 = decoder.Float64
	type decI64 = decoder.Int64

	return decoder.Unmarshal(
		msg.Data,
		&lo.CurrencyPairID,
		&lo.OrderNumber,
		&lo.OrderType,
		&decF64{Value: &lo.Rate},
		&decF64{Value: &lo.Amount},
		&lo.Date,
		&decF64{Value: &lo.OriginalAmount},
		&decI64{Value: &lo.ClientOrderID},
	)
}
