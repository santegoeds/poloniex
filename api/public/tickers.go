package public

import (
	"context"

	"github.com/santegoeds/poloniex/api/decoder"
)

// Ticker provides a summary of the market state for a currency pair.
type Ticker struct {
	// ID of the currency pair.
	ID int `json:"id"`
	// Price of the most recent trade.
	Last float64 `json:"last,string"`
	// Lowest available ASK price.
	LowestAsk float64 `json:"lowestAsk,string"`
	// Highest available BID price.
	HighestBid float64 `json:"highestBid,string"`
	// Price change as a percentage.
	PercentChange float64 `json:"percentChange,string"`
	// Number of base units that have traded within the last 24 hours.
	BaseVolume float64 `json:"baseVolume,string"`
	// Number of quote units that have traded within the last 24 hours.
	QuoteVolume float64 `json:"quoteVolume,string"`
	// Flag that indicates if the market for the pair is available for trading.
	Frozen string `json:"isFrozen"`
	// Highest price at which the pair traded within the last 24 hours.
	High24h float64 `json:"high24hr,string"`
	// Lowest price at which the pair traded within the last 24 hours.
	Low24h float64 `json:"low24hr,string"`
}

// IsFrozen indicates if the market for a the Ticker is unavailable for trading.
func (t *Ticker) IsFrozen() bool {
	return t.Frozen == "1"
}

type TickerRequest struct {
	client *Client
}

func NewTickerRequest(client *Client) *TickerRequest {
	return &TickerRequest{
		client: client,
	}
}

func (r TickerRequest) Do(ctx context.Context) (map[string]Ticker, error) {
	rsp, err := r.client.do(ctx, "command", "returnTicker")
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	msgData, err := decoder.DecodeObject(rsp.Body)
	if err != nil {
		return nil, err
	}
	tickers := make(map[string]Ticker)
	for currencyPair, cpData := range msgData {
		t := Ticker{}
		if err = decoder.DecodeMessage(cpData, &t); err != nil {
			return nil, err
		}
		tickers[currencyPair] = t
	}
	return tickers, nil
}
