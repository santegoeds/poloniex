package private

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/santegoeds/poloniex/errors"
)

type MoveOrder struct {
	Amount        float64 `json:"amount,string"`
	ClientOrderID int64   `json:"clientOrderId,string"`
	Message       string  `json:"message"`
}

type MoveOrderRequest struct {
	client        *Client
	orderNumber   int64
	rate          float64
	amount        float64
	clientOrderID int64
}

func NewMoveOrderRequest(client *Client, orderNumber int64, rate float64) *MoveOrderRequest {
	return &MoveOrderRequest{
		client:      client,
		orderNumber: orderNumber,
		rate:        rate,
	}
}

func (r *MoveOrderRequest) Amount(amount float64) *MoveOrderRequest {
	r.amount = amount
	return r
}

func (r *MoveOrderRequest) ClientOrderID(clientOrderID int64) *MoveOrderRequest {
	r.clientOrderID = clientOrderID
	return r
}

func (r *MoveOrderRequest) Do(ctx context.Context) (*MoveOrder, error) {
	params := []string{
		"command", "moveOrder",
		"orderNumber", strconv.FormatInt(r.orderNumber, 10),
		"rate", strconv.FormatFloat(r.rate, 'f', -1, 64),
	}
	if r.amount > 0.0 {
		params = append(params, "amount", strconv.FormatFloat(r.amount, 'f', -1, 64))
	}
	if r.clientOrderID != 0 {
		params = append(params, "clientOrderId", strconv.FormatInt(r.clientOrderID, 10))
	}

	rsp, err := r.client.do(ctx, params...)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	type moveOrder struct {
		MoveOrder
		Success int    `json:"success,string"`
		Error   string `json:"error"`
	}

	mo := moveOrder{}
	dec := json.NewDecoder(rsp.Body)
	if err = dec.Decode(&mo); err != nil {
		return nil, err
	}

	if mo.Error != "" {
		return nil, fmt.Errorf("%w: "+mo.Error, errors.ErrBadRequest)
	}
	if mo.Success != 1 {
		return nil, fmt.Errorf("%w: success flag indicates failure", errors.ErrBadRequest)
	}

	return &mo.MoveOrder, nil
}
