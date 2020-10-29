package event

import (
	"github.com/santegoeds/poloniex/api/decoder"
	"github.com/santegoeds/poloniex/message"
	"time"
)

var _ Event = &BookTrade{}

type BookTrade struct {
	CurrencyPairID int
	TradeID        string
	TradeType      int
	Price          float64
	Size           float64
	Time           int64
	SequenceNr     int64
}

func (bt *BookTrade) ChannelID() int {
	return bt.CurrencyPairID
}

func (bt *BookTrade) Type() message.Type {
	return message.BookTrade
}

func (bt *BookTrade) IsSell() bool {
	return bt.TradeType == 0
}

func (bt *BookTrade) IsBuy() bool {
	return bt.TradeType == 1
}

func (bt *BookTrade) DateTime() time.Time {
	return time.Unix(bt.Time, 0)
}

func (bt *BookTrade) Unmarshal(msg message.Message) error {
	type decF64 = decoder.Float64

	bt.CurrencyPairID = msg.ChannelID

	if err := decoder.Unmarshal(
		msg.Data,
		&bt.SequenceNr,
		&bt.TradeID,
		&bt.TradeType,
		&decF64{Value: &bt.Size},
		&decF64{Value: &bt.Price},
		&bt.Time,
	); err != nil {
		return err
	}
	return nil
}
