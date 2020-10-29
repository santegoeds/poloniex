package private

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"github.com/santegoeds/poloniex/api/decoder"
)

type OpenOrder struct {
	OrderNumber    int64   `json:"orderNumber,string"`
	Type           string  `json:"type"`
	Rate           float64 `json:"rate,string"`
	StartingAmount float64 `json:"startingAmount,string"`
	Amount         float64 `json:"amount,string"`
	Total          float64 `json:"total,string"`
	Date           string  `json:"date"`
	Margin         int     `json:"margin"`
}

func (oo *OpenOrder) DateTime() time.Time {
	ts, _ := time.Parse(timeFormat, oo.Date)
	return ts
}

type OpenOrders map[string][]OpenOrder

type OpenOrdersRequest struct {
	client       *Client
	currencyPair string
}

func NewOpenOrdersRequest(client *Client, currencyPair string) *OpenOrdersRequest {
	return &OpenOrdersRequest{
		client:       client,
		currencyPair: currencyPair,
	}
}

func (r *OpenOrdersRequest) Do(ctx context.Context) (OpenOrders, error) {
	rsp, err := r.client.do(ctx, "command", "returnOpenOrders", "currencyPair", r.currencyPair)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if r.currencyPair == "all" {
		return decodeAllMarkets(rsp.Body)
	}
	return decodeOneMarket(r.currencyPair, rsp.Body)
}

func decodeOneMarket(currencyPair string, r io.Reader) (OpenOrders, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	msgs := make([]json.RawMessage, 0)
	if err = decoder.DecodeMessage(data, &msgs); err != nil {
		return nil, err
	}

	marketOrders := make([]OpenOrder, 0, len(msgs))
	for _, msg := range msgs {
		oo := OpenOrder{}
		if err = json.Unmarshal(msg, &oo); err != nil {
			return nil, err
		}
		marketOrders = append(marketOrders, oo)
	}

	openOrders := make(OpenOrders)
	openOrders[currencyPair] = marketOrders
	return openOrders, nil
}

func decodeAllMarkets(r io.Reader) (OpenOrders, error) {
	objData, err := decoder.DecodeObject(r)
	if err != nil {
		return nil, err
	}
	openOrders := make(OpenOrders)
	for currencyPair, currencyPairData := range objData {
		orderData := make([]json.RawMessage, 0)
		if err = json.Unmarshal(currencyPairData, &orderData); err != nil {
			return nil, err
		}
		for _, data := range orderData {
			oo := OpenOrder{}
			if err = json.Unmarshal(data, &oo); err != nil {
				return nil, err
			}
			openOrders[currencyPair] = append(openOrders[currencyPair], oo)
		}
	}

	return openOrders, nil
}
