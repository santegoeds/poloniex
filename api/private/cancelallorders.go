package private

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/santegoeds/poloniex/errors"
)

type CancelAllOrders struct {
	Message      string  `json:"message"`
	OrderNumbers []int64 `json:"orderNumbers"`
}

type CancelAllOrdersRequest struct {
	client       *Client
	currencyPair string
}

func NewCancelAllOrdersRequest(client *Client, currencyPair string) *CancelAllOrdersRequest {
	return &CancelAllOrdersRequest{
		client:       client,
		currencyPair: currencyPair,
	}
}

func (r *CancelAllOrdersRequest) Do(ctx context.Context) (*CancelAllOrders, error) {
	params := []string{
		"command", "cancelAllOrders",
	}
	if r.currencyPair != "" {
		params = append(params, "currencyPair", r.currencyPair)
	}
	rsp, err := r.client.do(ctx, params...)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	type cancelAllOrders struct {
		CancelAllOrders
		Success int    `json:"success,string"`
		Error   string `json:"error"`
	}

	co := cancelAllOrders{}
	dec := json.NewDecoder(rsp.Body)
	if err := dec.Decode(&co); err != nil {
		return nil, err
	}
	if co.Error != "" {
		return nil, fmt.Errorf("%w: "+co.Error, errors.ErrBadRequest)
	}
	if co.Success != 1 {
		return nil, fmt.Errorf("%w: success flag indicates failure", errors.ErrBadRequest)
	}
	return &co.CancelAllOrders, nil
}
