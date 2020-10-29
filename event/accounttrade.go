package event

import (
	channel2 "github.com/santegoeds/poloniex/channel"
	"time"

	"github.com/santegoeds/poloniex/api/decoder"
	"github.com/santegoeds/poloniex/message"
)

var _ Event = &AccountTrade{}

type AccountTrade struct {
	TradeID       int
	Rate          float64
	Amount        float64
	FeeMultiplier float64
	FundingType   int
	OrderNumber   int64
	TotalFee      float64
	Date          string
	ClientOrderID int64
	TradeTotal    float64
}

func (at *AccountTrade) ChannelID() int {
	return int(channel2.AccountNotifications)
}

func (at *AccountTrade) Type() message.Type {
	return message.AccountTrade
}

func (at *AccountTrade) IsExchangeWallet() bool {
	return at.FundingType == 0
}

func (at *AccountTrade) IsBorrowedFunds() bool {
	return at.FundingType == 1
}

func (at *AccountTrade) IsMarginFunds() bool {
	return at.FundingType == 2
}

func (at *AccountTrade) IsLendingFunds() bool {
	return at.FundingType == 3
}

func (at *AccountTrade) DateTime() time.Time {
	ts, _ := time.Parse(timeFormat, at.Date)
	return ts
}

func (at *AccountTrade) Unmarshal(msg message.Message) error {
	type decF64 = decoder.Float64
	type decI64 = decoder.Int64

	return decoder.Unmarshal(
		msg.Data,
		&at.TradeID,
		&decF64{Value: &at.Rate},
		&decF64{Value: &at.Amount},
		&decF64{Value: &at.FeeMultiplier},
		&at.FundingType,
		&at.OrderNumber,
		&decF64{Value: &at.TotalFee},
		&at.Date,
		&decI64{Value: &at.ClientOrderID},
		&decF64{Value: &at.TradeTotal},
	)
}
