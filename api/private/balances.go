package private

import (
	"context"
	"encoding/json"

	"github.com/santegoeds/poloniex/api/decoder"
)

type Balances map[string]float64

type BalancesRequest struct {
	client *Client
}

func NewBalancesRequest(client *Client) *BalancesRequest {
	return &BalancesRequest{
		client: client,
	}
}

func (r *BalancesRequest) Do(ctx context.Context) (Balances, error) {
	type decF64 = decoder.Float64

	rsp, err := r.client.do(ctx, "command", "returnBalances")
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	objData, err := decoder.DecodeObject(rsp.Body)
	if err != nil {
		return nil, err
	}

	balances := make(Balances)
	for currency, balanceData := range objData {
		var balance float64
		if err = json.Unmarshal(balanceData, &decF64{Value: &balance}); err != nil {
			return nil, err
		}
		balances[currency] = balance
	}
	return balances, nil
}
