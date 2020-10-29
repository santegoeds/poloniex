package private

import (
	"context"
	"encoding/json"

	"github.com/santegoeds/poloniex/api/decoder"
)

type TradableCurrencyBalances map[string]float64
type TradableCurrencyPairBalances map[string]TradableCurrencyBalances

type TradableBalancesRequest struct {
	client *Client
}

func NewTradableBalancesRequest(client *Client) *TradableBalancesRequest {
	return &TradableBalancesRequest{
		client: client,
	}
}

func (r *TradableBalancesRequest) Do(ctx context.Context) (TradableCurrencyPairBalances, error) {
	type decF64 = decoder.Float64

	rsp, err := r.client.do(ctx, "command", "returnTradableBalances")
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	currencyPairData, err := decoder.DecodeObject(rsp.Body)
	if err != nil {
		return nil, err
	}

	currencyPairBalances := make(TradableCurrencyPairBalances)
	for currencyPair, currenciesData := range currencyPairData {
		data := make(map[string]json.RawMessage)
		if err = json.Unmarshal(currenciesData, &data); err != nil {
			return nil, err
		}

		currencyBalances := make(TradableCurrencyBalances)
		for currency, balanceData := range data {
			var balance float64
			if err = json.Unmarshal(balanceData, &decF64{Value: &balance}); err != nil {
				return nil, err
			}
			currencyBalances[currency] = balance
		}

		currencyPairBalances[currencyPair] = currencyBalances
	}

	return currencyPairBalances, nil
}
