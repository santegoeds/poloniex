package poloniex_test

import (
	"context"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"

	"github.com/santegoeds/poloniex"
	"github.com/santegoeds/poloniex/test"
)

func TestOrders(t *testing.T) {
	t.Skipf("Work out a testing strategy")

	ctx := context.TODO()
	clientOrderID := int64(9999)
	currencyPair := "ETH_USDT"
	cli := poloniex.New(os.Getenv("KEY"), os.Getenv("SECRET"))

	require.True(t, t.Run("should create a sell order", func(t *testing.T) {
		to := getTestOrder(test.WithT(ctx, t), cli, currencyPair)
		order, err := cli.Buy(currencyPair, to.Rate, to.Amount).
			ClientOrderID(clientOrderID).
			Do(context.TODO())
		require.NoError(t, err)

		t.Log(pretty.Sprint(order))

		defer func() {
			canceledOrder, err := cli.CancelOrder().
				OrderNumber(order.OrderNumber).
				Do(context.TODO())
			require.NoError(t, err)

			t.Log(pretty.Sprint(canceledOrder))
		}()
	}))

	// require.True(t, t.Run("should move an order", func(t *testing.T) {

	// }))

	// require.True(t, t.Run("should cancel an order", func(t *testing.T) {

	// }))
}

type TestOrder struct {
	Rate   float64
	Amount float64
}

func getTestOrder(ctx context.Context, client *poloniex.Client, currencyPair string) TestOrder {
	t := test.T(ctx)

	books, err := client.OrderBook(currencyPair).
		Do(ctx)
	require.NoError(t, err)

	book := books[currencyPair]
	worstBid := book.Bids[len(book.Bids)-1]

	quote := currencyPair[:3]

	minOrderSize := map[string]float64{
		"BTC": 0.0001,
		"ETH": 0.0001,
		"USD": 1.0,
		"TRX": 100.0,
		"BNB": 0.06,
	}

	return TestOrder{
		Rate:   worstBid.Price,
		Amount: minOrderSize[quote],
	}
}
