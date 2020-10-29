package event

import (
	"github.com/santegoeds/poloniex/api/decoder"
	channel2 "github.com/santegoeds/poloniex/channel"
	"github.com/santegoeds/poloniex/message"
)

var _ Event = &TickerData{}

type TickerData struct {
	CurrencyPairID         int
	LastTradePrice         float64
	LowestAsk              float64
	HighestBid             float64
	PercentChange24H       float64
	BaseCurrencyVolume24H  float64
	QuoteCurrencyVolume24H float64
	IsFrozen               int
	HighestTradePrice24H   float64
	LowestTradePrice24H    float64
}

func (td TickerData) ChannelID() int {
	return int(channel2.TickerData)
}

func (td TickerData) Type() message.Type {
	return message.TickerData
}

func (td *TickerData) Unmarshal(msg message.Message) error {
	type decF64 = decoder.Float64

	return decoder.Unmarshal(
		msg.Data,
		&td.CurrencyPairID,
		&decF64{Value: &td.LastTradePrice},
		&decF64{Value: &td.LowestAsk},
		&decF64{Value: &td.HighestBid},
		&decF64{Value: &td.PercentChange24H},
		&decF64{Value: &td.BaseCurrencyVolume24H},
		&decF64{Value: &td.QuoteCurrencyVolume24H},
		&td.IsFrozen,
		&decF64{Value: &td.HighestTradePrice24H},
		&decF64{Value: &td.LowestTradePrice24H},
	)
}
