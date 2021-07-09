package event

import (
	"time"

	"github.com/santegoeds/poloniex/api/decoder"
	"github.com/santegoeds/poloniex/message"
)

var _ Event = &BookTrade{}

type BookTrade struct {
	CurrencyPairID int
	TradeID        string
	TradeType      int
	Price          float64
	Size           float64
	Time           time.Time
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

func (bt *BookTrade) Unmarshal(msg message.Message) error {
	type decF64 = decoder.Float64
	type decEpochMs = decoder.EpochMs

	bt.CurrencyPairID = msg.ChannelID

	if err := decoder.Unmarshal(
		msg.Data,
		&bt.SequenceNr,
		&bt.TradeID,
		&bt.TradeType,
		&decF64{Value: &bt.Size},
		&decF64{Value: &bt.Price},
		&decEpochMs{Value: &bt.Time},
	); err != nil {
		return err
	}
	return nil
}
