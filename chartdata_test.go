package poloniex_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/santegoeds/poloniex"
	"github.com/santegoeds/poloniex/period"

	"github.com/stretchr/testify/require"
)

func TestChartData(t *testing.T) {
	currencyPair := "USDT_BTC"
	periods := []period.Period{
		period.M5,
		period.M15,
		period.M30,
		period.H2,
		period.H4,
		period.D,
	}

	ctx := context.TODO()
	cli := poloniex.New(os.Getenv("KEY"), os.Getenv("SECRET"))
	defer cli.Close()

	for _, period := range periods {
		t.Run(fmt.Sprintf("chart data for period %d", period), func(t *testing.T) {
			data, err := cli.ChartData(currencyPair, period).
				BarCount(100).
				Do(ctx)
			require.NoError(t, err)
			require.Len(t, data, 100)

			for _, bar := range data {
				require.Greater(t, bar.Low, 0.0)
				require.Greater(t, bar.Open, 0.0)
				require.Greater(t, bar.Close, 0.0)
				require.Greater(t, bar.Low, 0.0)
				require.Greater(t, bar.WeightedAverage, 0.0)
				require.GreaterOrEqual(t, bar.Volume, 0.0)
				require.GreaterOrEqual(t, bar.QuoteVolume, 0.0)
				require.GreaterOrEqual(t, bar.High, bar.Low)
				require.GreaterOrEqual(t, bar.High, bar.Open)
				require.GreaterOrEqual(t, bar.High, bar.Close)
				require.LessOrEqual(t, bar.Low, bar.Open)
				require.LessOrEqual(t, bar.Low, bar.Close)
			}

			t.Log(pretty.Sprint(data))
		})
	}
}
