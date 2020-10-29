package event

import (
	"github.com/santegoeds/poloniex/api/decoder"
	channel2 "github.com/santegoeds/poloniex/channel"
	"github.com/santegoeds/poloniex/message"
)

type MarginPositionUpdate struct {
	OrderNumber   int64
	CurrencyID    int
	Amount        float64
	ClientOrderID int64
}

func (mpu *MarginPositionUpdate) ChannelID() int {
	return int(channel2.AccountNotifications)
}

func (mpu *MarginPositionUpdate) Type() message.Type {
	return message.AccountMarginPositionUpdate
}

func (mpu *MarginPositionUpdate) Unmarshal(msg message.Message) error {
	type decF64 = decoder.Float64
	type decI64 = decoder.Int64

	return decoder.Unmarshal(
		msg.Data,
		&mpu.OrderNumber,
		&mpu.CurrencyID,
		&decF64{Value: &mpu.Amount},
		&decI64{Value: &mpu.ClientOrderID},
	)
}
