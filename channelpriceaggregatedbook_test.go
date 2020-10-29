package poloniex_test

import (
	"context"
	"math"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/santegoeds/poloniex"
	"github.com/santegoeds/poloniex/event"
	"github.com/santegoeds/poloniex/message"

	"github.com/stretchr/testify/require"
)

func TestPriceAggregatedBook(t *testing.T) {
	const (
		mskBook       = 0x1
		mskBookUpdate = 0x10
		mskBookTrade  = 0x100
		mskAll        = (mskBook | mskBookUpdate | mskBookTrade)
	)

	ctx := context.TODO()
	cli := poloniex.New(os.Getenv("KEY"), os.Getenv("SECRET"))
	defer cli.Close()

	currencyPairID, currencyPair := 121, "USDT_BTC"

	subscription, err := cli.Subscribe().
		CurrencyPair(currencyPair).
		Do(ctx)
	require.NoError(t, err)

	msgC, errC := subscription.Channels()

	for testMask := 0; testMask != mskAll; {
		select {
		case err = <-errC:
			t.Fatal(err)

		case msg := <-msgC:
			require.Equal(t, msg.ChannelID, currencyPairID)

			switch msg.Type {
			case message.Book:
				if (testMask & mskBook) > 0 {
					continue
				}
				testMask |= mskBook
				ok := t.Run("validate book event", func(t *testing.T) {
					validateBookEvent(t, msg, currencyPair)
				})
				require.True(t, ok)

			case message.BookUpdate:
				if (testMask & mskBookUpdate) > 0 {
					continue
				}
				testMask |= mskBookUpdate
				ok := t.Run("validate book update event", func(t *testing.T) {
					validateBookUpdateEvent(t, msg)
				})
				require.True(t, ok)

			case message.BookTrade:
				if (testMask & mskBookTrade) > 0 {
					continue
				}
				testMask |= mskBookTrade
				ok := t.Run("validate book trade event", func(t *testing.T) {
					validateBookTradeEvent(t, msg)
				})
				require.True(t, ok)

			default:
				t.Fatalf("unexpected event type %d", msg.Type)
			}
		}
	}
}

func validateBookEvent(t *testing.T, msg message.Message, currencyPair string) {
	evt := event.Book{}
	err := evt.Unmarshal(msg)
	require.NoError(t, err)

	t.Log(pretty.Sprint(evt))

	require.Equal(t, msg.ChannelID, evt.CurrencyPairID)
	require.Equal(t, currencyPair, evt.CurrencyPair)
	require.NotEmpty(t, evt.Asks)
	require.NotEmpty(t, evt.Bids)
	require.Greater(t, evt.SequenceNr, int64(0))

	maxBidPrice := math.Inf(-1)
	for _, bid := range evt.Bids {
		require.Greater(t, bid.Price, 0.0)
		require.Greater(t, bid.Size, 0.0)
		if bid.Price > maxBidPrice {
			maxBidPrice = bid.Price
		}
	}
	for _, ask := range evt.Asks {
		require.Greater(t, ask.Price, 0.0)
		require.Greater(t, ask.Size, 0.0)
		require.Greater(t, ask.Price, maxBidPrice)
	}
}

func validateBookUpdateEvent(t *testing.T, msg message.Message) {
	evt := event.BookUpdate{}
	err := evt.Unmarshal(msg)
	require.NoError(t, err)

	t.Log(pretty.Sprint(evt))

	require.Equal(t, msg.ChannelID, evt.CurrencyPairID)
	require.Contains(t, []int{0, 1}, evt.OrderType)
	require.Greater(t, evt.Order.Price, 0.0)
	require.GreaterOrEqual(t, evt.Order.Size, 0.0)
	require.Greater(t, evt.SequenceNr, int64(0))
}

func validateBookTradeEvent(t *testing.T, msg message.Message) {
	evt := event.BookTrade{}
	err := evt.Unmarshal(msg)
	require.NoError(t, err)

	t.Log(pretty.Sprint(evt))

	require.Equal(t, msg.ChannelID, evt.CurrencyPairID)
	require.Greater(t, evt.Price, 0.0)
	require.Greater(t, evt.Size, 0.0)
	require.NotZero(t, evt.Time)
	require.NotZero(t, evt.TradeID)
	require.Contains(t, []int{0, 1}, evt.TradeType)
	require.Greater(t, evt.SequenceNr, int64(0))
}
