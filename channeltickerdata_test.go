package poloniex_test

import (
	"context"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/santegoeds/poloniex"
	"github.com/santegoeds/poloniex/channel"
	"github.com/santegoeds/poloniex/event"
	"github.com/santegoeds/poloniex/message"

	"github.com/stretchr/testify/require"
)

func TestTickerData(t *testing.T) {
	ctx := context.TODO()
	cli := poloniex.New(os.Getenv("KEY"), os.Getenv("SECRET"))
	defer cli.Close()

	sub, err := cli.Subscribe().
		Channel(channel.TickerData).
		Do(ctx)
	require.NoError(t, err)

	msgC, errC := sub.Channels()

	select {
	case err := <-errC:
		t.Fatal(err)

	case msg := <-msgC:
		require.Equal(t, msg.ChannelID, int(channel.TickerData))
		require.Equal(t, msg.Type, message.TickerData)

		evt, err := poloniex.Decode(msg)
		require.NoError(t, err)

		data, ok := evt.(*event.TickerData)
		require.True(t, ok)
		t.Log(pretty.Sprint(data))

		require.NotZero(t, data.CurrencyPairID)
		require.Greater(t, data.BaseCurrencyVolume24H, 0.0)
		require.Greater(t, data.HighestBid, 0.0)
		require.Greater(t, data.HighestTradePrice24H, 0.0)
		require.Greater(t, data.LastTradePrice, 0.0)
		require.Contains(t, []int{0, 1}, data.IsFrozen)
		require.Greater(t, data.LowestAsk, 0.0)
		require.Greater(t, data.LowestTradePrice24H, 0.0)
		require.NotEqual(t, data.PercentChange24H, 0.0)
		require.Greater(t, data.QuoteCurrencyVolume24H, 0.0)
	}
}
