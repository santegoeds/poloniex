package poloniex_test

import (
	"context"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"

	"github.com/santegoeds/poloniex"
)

func TestCompleteBalances(t *testing.T) {
	ctx := context.TODO()
	cli := poloniex.New(os.Getenv("KEY"), os.Getenv("SECRET"))
	defer cli.Close()

	t.Run("Balances should return complete balances", func(t *testing.T) {
		balances, err := cli.CompleteBalances().
			Do(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, balances)

		t.Log(pretty.Sprint(balances))

		hasPositiveAvailable, hasPositiveBtcValue := false, false
		for _, balance := range balances {
			require.GreaterOrEqual(t, balance.Available, 0.0)
			require.GreaterOrEqual(t, balance.OnOrders, 0.0)
			require.GreaterOrEqual(t, balance.BtcValue, 0.0)
			if balance.Available > 0.0 {
				hasPositiveAvailable = true
			}
			if balance.BtcValue > 0.0 {
				hasPositiveBtcValue = true
			}
		}
		require.Truef(t, hasPositiveAvailable, "No asset with an available value")
		require.Truef(t, hasPositiveBtcValue, "No asset with a btc value")
	})
}
