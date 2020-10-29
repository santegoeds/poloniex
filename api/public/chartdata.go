package public

import (
	"context"
	"github.com/santegoeds/poloniex/period"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/santegoeds/poloniex/api/decoder"
)

type Bar struct {
	Time            int64   `json:"date"`
	High            float64 `json:"high"`
	Low             float64 `json:"low"`
	Open            float64 `json:"open"`
	Close           float64 `json:"close"`
	Volume          float64 `json:"volume"`
	QuoteVolume     float64 `json:"quoteVolume"`
	WeightedAverage float64 `json:"weightedAverage"`
}

type ChartDataRequest struct {
	client       *Client
	currencyPair string
	period       period.Period
	startTime    int64
	endTime      int64
}

func NewChartDataRequest(
	client *Client, currencyPair string, period period.Period,
) *ChartDataRequest {
	return &ChartDataRequest{
		client:       client,
		currencyPair: currencyPair,
		period:       period,
	}
}

func (r *ChartDataRequest) BarCount(count int) *ChartDataRequest {
	if r.endTime == 0 {
		r.endTime = time.Now().Unix()
	}
	// Start- and end times are inclusive, so we need to decrement by 1.
	r.startTime = r.endTime - int64(r.period)*int64(count-1)
	return r
}

func (r *ChartDataRequest) StartTime(startTime int64) *ChartDataRequest {
	r.startTime = startTime
	return r
}

func (r *ChartDataRequest) EndTime(endTime int64) *ChartDataRequest {
	r.endTime = endTime
	return r
}

func (r *ChartDataRequest) Do(ctx context.Context) ([]Bar, error) {
	rsp, err := r.client.do(
		ctx,
		"command", "returnChartData",
		"currencyPair", r.currencyPair,
		"period", strconv.FormatInt(int64(r.period), 10),
		"start", strconv.FormatInt(r.startTime, 10),
		"end", strconv.FormatInt(r.endTime, 10),
	)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	bars := make([]Bar, 0)
	if err = decoder.DecodeMessage(data, &bars); err != nil {
		return nil, err
	}
	return bars, nil
}
