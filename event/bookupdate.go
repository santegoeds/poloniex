package event

import (
	"github.com/santegoeds/poloniex/api/decoder"
	"github.com/santegoeds/poloniex/message"
)

var _ Event = &BookUpdate{}

type BookUpdate struct {
	CurrencyPairID int
	OrderType      int
	Order          Order
	SequenceNr     int64
}

func (bu *BookUpdate) ChannelID() int {
	return bu.CurrencyPairID
}

func (bu *BookUpdate) Type() message.Type {
	return message.BookUpdate
}

func (bu *BookUpdate) Unmarshal(msg message.Message) error {
	type decF64 = decoder.Float64
	type decEpochMs = decoder.EpochMs

	bu.CurrencyPairID = msg.ChannelID

	return decoder.Unmarshal(
		msg.Data,
		&bu.SequenceNr,
		&bu.OrderType,
		&decF64{Value: &bu.Order.Price},
		&decF64{Value: &bu.Order.Size},
		&decEpochMs{Value: &bu.Order.Time},
	)
}
