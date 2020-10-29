package private

import (
	"context"
	"encoding/json"

	"github.com/santegoeds/poloniex/api/decoder"
)

type CompleteBalance struct {
	Available float64 `json:"available,string"`
	OnOrders  float64 `json:"onOrders,string"`
	BtcValue  float64 `json:"btcValue,string"`
}

type CompleteBalances map[string]CompleteBalance

type CompleteBalancesRequest struct {
	client  *Client
	account string
}

func NewCompleteBalancesRequest(client *Client) *CompleteBalancesRequest {
	return &CompleteBalancesRequest{
		client: client,
	}
}

func (r *CompleteBalancesRequest) Account(account string) *CompleteBalancesRequest {
	r.account = account
	return r
}

func (r *CompleteBalancesRequest) Do(ctx context.Context) (CompleteBalances, error) {
	params := []string{
		"command", "returnCompleteBalances",
	}
	if r.account != "" {
		params = append(params, "account", r.account)
	}
	rsp, err := r.client.do(ctx, params...)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	objData, err := decoder.DecodeObject(rsp.Body)
	if err != nil {
		return nil, err
	}

	balances := make(CompleteBalances)
	for currency, balanceData := range objData {
		balance := CompleteBalance{}
		if err = json.Unmarshal(balanceData, &balance); err != nil {
			return nil, err
		}
		balances[currency] = balance
	}
	return balances, nil
}
