package poloniex_test

import (
	"context"
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/santegoeds/poloniex"
	"github.com/santegoeds/poloniex/api/public"
	"github.com/stretchr/testify/require"
)

func TestOrderbook(t *testing.T) {
	ctx := context.TODO()
	cli := poloniex.New(os.Getenv("KEY"), os.Getenv("SECRET"))

	currencyPair := "USDT_BTC"
	t.Run(fmt.Sprintf("should return an orderbook for currency pair %s", currencyPair), func(t *testing.T) {
		books, err := cli.OrderBook(currencyPair).
			Do(ctx)
		require.NoError(t, err)

		t.Log(pretty.Sprint(books))
		validateBook(t, currencyPair, books)
	})

	t.Run("should return orderbooks for all currency pairs", func(t *testing.T) {
		books, err := cli.OrderBook("all").
			Do(ctx)
		require.NoError(t, err)

		t.Log(pretty.Sprint(books))

		require.Greater(t, len(books), 1)
		for cp := range books {
			validateBook(t, cp, books)
		}
	})
}

func validateBook(t *testing.T, currencyPair string, books map[string]public.Book) {
	require.Contains(t, books, currencyPair)

	b := books[currencyPair]
	require.Contains(t, []string{"0", "1"}, b.Frozen)
	require.Greater(t, b.SequenceNr, 0)

	maxBidPrice, lastBidPrice := 0.0, math.Inf(1)
	for _, o := range b.Bids {
		require.Greater(t, o.Price, 0.0)
		require.Greater(t, o.Size, 0.0)

		require.Greater(t, lastBidPrice, o.Price)
		lastBidPrice = o.Price
		if o.Price > maxBidPrice {
			maxBidPrice = o.Price
		}
	}

	lastAskPrice := math.Inf(-1)
	for _, o := range b.Asks {
		require.Greater(t, o.Price, 0.0)
		require.Greater(t, o.Size, 0.0)

		require.Greater(t, o.Price, lastAskPrice)
		require.Greater(t, o.Price, maxBidPrice)

		lastAskPrice = o.Price
	}
}
