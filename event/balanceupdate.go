package event

import (
	"github.com/santegoeds/poloniex/api/decoder"
	channel2 "github.com/santegoeds/poloniex/channel"
	"github.com/santegoeds/poloniex/message"
)

var _ Event = &BalanceUpdate{}

type BalanceUpdate struct {
	CurrencyID int
	Wallet     string
	Amount     float64
}

func (bu *BalanceUpdate) ChannelID() int {
	return int(channel2.AccountNotifications)
}

func (bu *BalanceUpdate) Type() message.Type {
	return message.AccountBalanceUpdate
}

func (bu *BalanceUpdate) IsExchange() bool {
	return bu.Wallet == "e"
}

func (bu *BalanceUpdate) IsMargin() bool {
	return bu.Wallet == "m"
}

func (bu *BalanceUpdate) IsLending() bool {
	return bu.Wallet == "l"
}

func (bu *BalanceUpdate) Unmarshal(msg message.Message) error {
	return decoder.Unmarshal(
		msg.Data,
		&bu.CurrencyID,
		&bu.Wallet,
		&decoder.Float64{Value: &bu.Amount},
	)
}
