package public

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/santegoeds/poloniex/api/decoder"
)

type CurrencyVolumes = map[string]float64

type Volume24h struct {
	currencyPairs map[string]CurrencyVolumes
	totalBTC      float64
	totalETH      float64
	totalUSDC     float64
}

type Volume24hRequest struct {
	client *Client
}

func NewVolume24hRequest(client *Client) *Volume24hRequest {
	return &Volume24hRequest{
		client: client,
	}
}

func (r Volume24hRequest) Do(ctx context.Context) (*Volume24h, error) {
	rsp, err := r.client.do(ctx, "command", "return24hVolume")
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	objData, err := decoder.DecodeObject(rsp.Body)
	if err != nil {
		return nil, err
	}

	volume24h := &Volume24h{
		currencyPairs: make(map[string]CurrencyVolumes),
	}
	for currencyPair, data := range objData {
		switch currencyPair {
		case "totalBTC":
			wrappedVol := decoder.Float64{&volume24h.totalBTC}
			if err = json.Unmarshal(data, &wrappedVol); err != nil {
				return nil, err
			}

		case "totalETH":
			wrappedVol := decoder.Float64{&volume24h.totalETH}
			if err = json.Unmarshal(data, &wrappedVol); err != nil {
				return nil, err
			}

		case "totalUSDC":
			wrappedVol := decoder.Float64{&volume24h.totalUSDC}
			if err = json.Unmarshal(data, &wrappedVol); err != nil {
				return nil, err
			}

		default:
			volumesAsStr := make(map[string]string)
			if err = json.Unmarshal(data, &volumesAsStr); err != nil {
				return nil, err
			}

			currencies := make(map[string]float64)
			for currency, volAsStr := range volumesAsStr {
				if currencies[currency], err = strconv.ParseFloat(volAsStr, 64); err != nil {
					return nil, err
				}
			}
			volume24h.currencyPairs[currencyPair] = currencies
		}
	}

	return volume24h, nil
}
