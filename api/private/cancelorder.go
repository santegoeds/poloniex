package private

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/santegoeds/poloniex/errors"
)

type CancelOrder struct {
	Amount        float64 `json:"amount,string"`
	ClientOrderID int64   `json:"clientOrderId,string"`
	Message       string  `json:"message"`
}

type CancelOrderRequest struct {
	client        *Client
	orderNumber   int64
	clientOrderID int64
}

func NewCancelOrderRequest(client *Client) *CancelOrderRequest {
	return &CancelOrderRequest{
		client: client,
	}
}

func (r *CancelOrderRequest) OrderNumber(orderNumber int64) *CancelOrderRequest {
	r.orderNumber = orderNumber
	return r
}

func (r *CancelOrderRequest) ClientOrderID(clientOrderID int64) *CancelOrderRequest {
	r.clientOrderID = clientOrderID
	return r
}

func (r *CancelOrderRequest) Do(ctx context.Context) (*CancelOrder, error) {
	params := []string{
		"command", "cancelOrder",
	}
	if r.orderNumber > 0 {
		params = append(params, "orderNumber", strconv.FormatInt(r.orderNumber, 10))
	}
	if r.clientOrderID > 0 {
		params = append(params, "clientOrderId", strconv.FormatInt(r.clientOrderID, 10))
	}
	if len(params) == 2 {
		return nil, fmt.Errorf("one of OrderNumber or ClientOrderID is required: %w", errors.ErrBadRequest)
	}

	rsp, err := r.client.do(ctx, params...)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	type cancelOrder struct {
		CancelOrder
		Success int `json:"success,string"`
	}

	cancelRsp := cancelOrder{}
	dec := json.NewDecoder(rsp.Body)
	if err = dec.Decode(&cancelRsp); err != nil {
		return nil, err
	}

	if cancelRsp.Success != 1 {
		return nil, fmt.Errorf("success flag indicates failure: %w", errors.ErrBadRequest)
	}
	return &cancelRsp.CancelOrder, nil
}
