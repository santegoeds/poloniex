package private

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/santegoeds/poloniex/errors"
)

type OrderSide string

const (
	Buy  OrderSide = "buy"
	Sell OrderSide = "sell"
)

type Trade struct {
	Amount  float64 `json:"amount,string"`
	Date    string  `json:"date"`
	Rate    float64 `json:"rate,string"`
	Total   float64 `json:"total,string"`
	TradeID string  `json:"tradeId"`
	Type    string  `json:"type"`
}

func (t *Trade) DateTime() time.Time {
	ts, _ := time.Parse(timeFormat, t.Date)
	return ts
}

type Order struct {
	OrderNumber      int64   `json:"orderNumber"`
	ResultingTrades  []Trade `json:"resultingTrades"`
	Fee              float64 `json:"fee,string"`
	ClientOrderID    int64   `json:"clientOrderId,string"`
	CurrencyPair     string  `json:"currencyPair"`
	TokenFee         float64 `json:"tokenFee,string"`
	TokenFeeCurrency string  `json:"tokenFeeCurrency"`
}

type OrderRequest struct {
	client            *Client
	side              string
	currencyPair      string
	rate              float64
	amount            float64
	fillOrKill        bool
	immediateOrCancel bool
	postOnly          bool
	clientOrderID     int64
}

func NewOrderRequest(client *Client, side OrderSide, currencyPair string) *OrderRequest {
	return &OrderRequest{
		client:       client,
		side:         string(side),
		currencyPair: currencyPair,
	}
}

func (r *OrderRequest) Rate(rate float64) *OrderRequest {
	r.rate = rate
	return r
}

func (r *OrderRequest) Amount(amount float64) *OrderRequest {
	r.amount = amount
	return r
}

func (r *OrderRequest) FillOrKill(flag bool) *OrderRequest {
	r.fillOrKill = flag
	return r
}

func (r *OrderRequest) ImmediateOrCancel(flag bool) *OrderRequest {
	r.immediateOrCancel = flag
	return r
}

func (r *OrderRequest) PostOnly(flag bool) *OrderRequest {
	r.postOnly = flag
	return r
}

func (r *OrderRequest) ClientOrderID(clientOrderID int64) *OrderRequest {
	r.clientOrderID = clientOrderID
	return r
}

func (r *OrderRequest) Do(ctx context.Context) (*Order, error) {
	params := []string{
		"command", r.side,
		"currencyPair", r.currencyPair,
		"rate", strconv.FormatFloat(r.rate, 'f', -1, 64),
		"amount", strconv.FormatFloat(r.amount, 'f', -1, 64),
	}
	if r.fillOrKill {
		params = append(params, "fillOrKill", "1")
	}
	if r.immediateOrCancel {
		params = append(params, "immediateOrCancel", "1")
	}
	if r.postOnly {
		params = append(params, "postOnly", "1")
	}
	if r.clientOrderID > 0 {
		params = append(params, "clientOrderId", strconv.FormatInt(r.clientOrderID, 10))
	}

	rsp, err := r.client.do(ctx, params...)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	type wrappedOrder struct {
		Order
		ResultingTrades []Trade
		Error           string `json:"error"`
	}

	var wOrder wrappedOrder
	dec := json.NewDecoder(rsp.Body)
	if err := dec.Decode(&wOrder); err != nil {
		return nil, err
	}
	if wOrder.Error != "" {
		return nil, fmt.Errorf("%w: "+wOrder.Error, errors.ErrBadRequest)
	}

	order := &wOrder.Order
	order.ResultingTrades = wOrder.ResultingTrades

	return order, nil
}
