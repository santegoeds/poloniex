package public

import (
	"context"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/santegoeds/poloniex/api/decoder"
)

type Trade struct {
	GlobalTradeID int64   `json:"globalTradeID"`
	TradeID       int64   `json:"tradeID"`
	Date          string  `json:"date"`
	Side          string  `json:"type"`
	Rate          float64 `json:"rate,string"`
	Amount        float64 `json:"amount,string"`
	Total         float64 `json:"total,string"`
}

func (t *Trade) IsBuy() bool {
	return t.Side == "buy"
}

func (t *Trade) IsSell() bool {
	return t.Side == "sell"
}

func (t *Trade) DateTime() time.Time {
	ts, _ := time.Parse(timeFormat, t.Date)
	return ts
}

type TradeHistoryRequest struct {
	client       *Client
	currencyPair string
	startTime    int64
	endTime      int64
}

func NewTradeHistoryRequest(client *Client, currencyPair string) *TradeHistoryRequest {
	return &TradeHistoryRequest{
		client:       client,
		currencyPair: currencyPair,
	}
}

func (r *TradeHistoryRequest) StartTime(startTime int64) *TradeHistoryRequest {
	r.startTime = startTime
	return r
}

func (r *TradeHistoryRequest) EndTime(endTime int64) *TradeHistoryRequest {
	r.endTime = endTime
	return r
}

func (r *TradeHistoryRequest) Do(ctx context.Context) ([]Trade, error) {
	params := []string{
		"command", "returnTradeHistory",
		"currencyPair", r.currencyPair,
	}
	if r.startTime > 0 {
		params = append(params, "start", strconv.FormatInt(r.startTime, 10))
	}
	if r.endTime > 0 {
		params = append(params, "end", strconv.FormatInt(r.endTime, 10))
	}

	rsp, err := r.client.do(ctx, params...)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	trades := make([]Trade, 0, 200)
	if err = decoder.DecodeMessage(data, &trades); err != nil {
		return nil, err
	}
	return trades, nil

}
