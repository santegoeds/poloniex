package poloniex_test

import (
	"context"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/santegoeds/poloniex"
)

func TestTickers(t *testing.T) {
	ctx := context.TODO()
	client := poloniex.New(os.Getenv("KEY"), os.Getenv("SECRET"))

	tickers, err := client.Tickers().
		Do(ctx)
	require.NoError(t, err)

	t.Log(pretty.Sprint(tickers))

	var count struct {
		LowestAsk     int
		HighestBid    int
		PercentChange int
		BaseVolume    int
		QuoteVolume   int
		High24h       int
		Low24h        int
	}

	require.NotEmpty(t, tickers)
	for currencyPair, ticker := range tickers {
		require.NotEmpty(t, currencyPair)
		require.NotZero(t, ticker.ID)
		require.GreaterOrEqual(t, ticker.Last, 0.0)
		require.GreaterOrEqual(t, ticker.LowestAsk, 0.0)
		require.GreaterOrEqual(t, ticker.HighestBid, 0.0)
		require.GreaterOrEqual(t, ticker.BaseVolume, 0.0)
		require.GreaterOrEqual(t, ticker.QuoteVolume, 0.0)
		require.GreaterOrEqual(t, ticker.High24h, 0.0)
		require.GreaterOrEqual(t, ticker.Low24h, 0.0)
		assert.Contains(t, "01", ticker.Frozen)

		if ticker.LowestAsk == 0.0 {
			count.LowestAsk++
		}
		if ticker.HighestBid == 0.0 {
			count.HighestBid++
		}
		if ticker.PercentChange == 0.0 {
			count.PercentChange++
		}
		if ticker.BaseVolume == 0.0 {
			count.BaseVolume++
		}
		if ticker.QuoteVolume == 0.0 {
			count.QuoteVolume++
		}
		if ticker.High24h == 0.0 {
			count.High24h++
		}
		if ticker.Low24h == 0.0 {
			count.Low24h++
		}
	}

	// Some of these might be zero; but if they are all zero then there's likely a bug.
	assert.Greater(t, len(tickers), count.LowestAsk)
	assert.Greater(t, len(tickers), count.HighestBid)
	assert.Greater(t, len(tickers), count.PercentChange)
	assert.Greater(t, len(tickers), count.BaseVolume)
	assert.Greater(t, len(tickers), count.QuoteVolume)
	assert.Greater(t, len(tickers), count.High24h)
	assert.Greater(t, len(tickers), count.Low24h)
}
