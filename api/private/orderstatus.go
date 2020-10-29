package private

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/kr/pretty"

	"github.com/santegoeds/poloniex/errors"
)

type OrderStatus struct {
	OrderNumber    int64
	Status         string  `json:"status"`
	Rate           float64 `json:"rate,string"`
	Amount         float64 `json:"amount,string"`
	CurrencyPair   string  `json:"currencyPair"`
	Date           string  `json:"date"`
	Total          float64 `json:"total,string"`
	Type           string  `json:"type"`
	StartingAmount float64 `json:"startingAmount,string"`
}

func (o *OrderStatus) IsOpen() bool {
	return o.Status == "Open"
}

func (o *OrderStatus) IsPartiallyFilled() bool {
	return o.Status == "Partially filled"
}

func (o *OrderStatus) DateTime() time.Time {
	ts, _ := time.Parse(timeFormat, o.Date)
	return ts
}

type OrderStatusRequest struct {
	client      *Client
	orderNumber int64
}

func NewOrderStatusRequest(client *Client, orderNumber int64) *OrderStatusRequest {
	return &OrderStatusRequest{
		client:      client,
		orderNumber: orderNumber,
	}
}

func (r *OrderStatusRequest) Do(ctx context.Context) (*OrderStatus, error) {
	orderNumberAsStr := strconv.FormatInt(r.orderNumber, 64)
	rsp, err := r.client.do(ctx,
		"command", "returnOrderStatus",
		"orderNumber", orderNumberAsStr,
	)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	type responseObject struct {
		Result  map[string]OrderStatus `json:"result"`
		Success int                    `json:"success"`
	}

	var rspObj responseObject
	dec := json.NewDecoder(rsp.Body)
	if err := dec.Decode(&rspObj); err != nil {
		return nil, err
	}

	if rspObj.Success != 1 {
		return nil, fmt.Errorf("success flag indicates failure: %w", errors.ErrBadRequest)
	}

	orderStatus, ok := rspObj.Result[orderNumberAsStr]
	if !ok {
		return nil, fmt.Errorf("response does not have a status: %s", pretty.Sprint(rspObj))
	}

	return &orderStatus, nil
}
