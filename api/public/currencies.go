package public

import (
	"context"

	"github.com/santegoeds/poloniex/api/decoder"
)

// Currency holds information about a currency or token.
type Currency struct {
	// Unique currency identifier.
	ID int `json:"id"`
	// Name of the currency.
	Name string `json:"name"`
	// Network fee needed to withdraw the currency.
	TxFee float64 `json:"txFee,string"`
	// Minimum number of confirmation blocks before a deposit is credited to an account.
	MinConf int `json:"minConf"`
	// Deposit account address (Optional)
	DepositAddress string `json:"depositAddress"`
	// Flag that indicates whether deposits and withdrawals are disabled.
	Disabled int `json:"disabled"`
	// Flag that indicates whether the currency has been delisted.
	Delisted int `json:"delisted"`
	// Flag that indicates whether the currency is unavailable for trading.
	Frozen int `json:"frozen"`
}

// IsDisabled indicates whether deposits and withdrawals have been disabled.
func (c *Currency) IsDisabled() bool {
	return c.Disabled == 1
}

// IsDelisted indicates whether a currency has been delisted.
func (c *Currency) IsDelisted() bool {
	return c.Delisted == 1
}

// IsFrozen indicates whether a currency is unavailable for trading.
func (c *Currency) IsFrozen() bool {
	return c.Frozen == 1
}

// HasDepositAddress indicates whether a currency has a deposit address.
func (c *Currency) HasDepositAddress() bool {
	return c.DepositAddress != ""
}

type CurrenciesRequest struct {
	client *Client
}

func NewCurrenciesRequest(client *Client) *CurrenciesRequest {
	return &CurrenciesRequest{
		client: client,
	}
}

func (r *CurrenciesRequest) Do(ctx context.Context) (map[string]Currency, error) {
	rsp, err := r.client.do(ctx, "command", "returnCurrencies")
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	objData, err := decoder.DecodeObject(rsp.Body)
	if err != nil {
		return nil, err
	}
	currencies := make(map[string]Currency)
	for curName, curData := range objData {
		cur := Currency{}
		if err := decoder.DecodeMessage(curData, &cur); err != nil {
			return nil, err
		}
		currencies[curName] = cur
	}

	return currencies, nil
}
