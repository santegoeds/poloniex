package event

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/santegoeds/poloniex/api/decoder"

	"github.com/santegoeds/poloniex/errors"
	"github.com/santegoeds/poloniex/message"
)

var _ Event = &Book{}

type Order struct {
	Price float64
	Size  float64
	Time  time.Time
}

type Book struct {
	CurrencyPairID int
	CurrencyPair   string
	Asks           []Order
	Bids           []Order
	SequenceNr     int64
	Time           time.Time
}

func (b *Book) ChannelID() int {
	return b.CurrencyPairID
}

func (b *Book) Type() message.Type {
	return message.Book
}

func (b *Book) Unmarshal(msg message.Message) error {
	var err error
	b.Asks = make([]Order, 0)
	b.Bids = make([]Order, 0)

	if err = json.Unmarshal(msg.Data[0], &b.SequenceNr); err != nil {
		return fmt.Errorf("%w: failed to decode Orderbook message for priceaggregatedbook", err)
	}

	var s struct {
		CurrencyPair string              `json:"currencyPair"`
		Book         []map[string]string `json:"orderBook"`
	}
	if err = json.Unmarshal(msg.Data[1], &s); err != nil {
		return fmt.Errorf("%w: failed to decode Orderbook message for priceaggregatedbook", err)
	}

	b.CurrencyPairID = msg.ChannelID
	b.CurrencyPair = s.CurrencyPair

	if len(s.Book) != 2 {
		return fmt.Errorf(
			"%w: orderbook of length %d is not a Ask/Bid tuple for priceaggregatebook",
			errors.ErrBadResponse,
			len(s.Book),
		)
	}

	d := decoder.EpochMs{Value: &b.Time}
	if err = d.UnmarshalJSON(msg.Data[2]); err != nil {
		return fmt.Errorf(
			"%w: failed ot decode epoch_ms %s for proceaggregatedbook", err, string(msg.Data[2]),
		)
	}

	// Asks
	for price, size := range s.Book[0] {
		o := Order{Time: b.Time}
		if o.Price, err = strconv.ParseFloat(price, 64); err != nil {
			return fmt.Errorf(
				"%w failed to decode price %s for priceaggregatedbook", err, price,
			)
		}
		if o.Size, err = strconv.ParseFloat(size, 64); err != nil {
			return fmt.Errorf(
				"%w: failed to decode size %s for priceaggregatedbook", err, size,
			)
		}
		b.Asks = append(b.Asks, o)
	}
	// Bids
	for price, size := range s.Book[1] {
		o := Order{Time: b.Time}
		if o.Price, err = strconv.ParseFloat(price, 64); err != nil {
			return fmt.Errorf(
				"%w: failed to decode price %s for priceaggregatedbook", err, price,
			)
		}
		if o.Size, err = strconv.ParseFloat(size, 64); err != nil {
			return fmt.Errorf(
				"%w: failed to decode size %s for priceaggregatedbook", err, size,
			)
		}
		b.Bids = append(b.Bids, o)
	}

	return nil
}
