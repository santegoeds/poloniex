package poloniex_test

import (
	"context"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"

	"github.com/santegoeds/poloniex"
)

func TestBalances(t *testing.T) {
	ctx := context.TODO()
	cli := poloniex.New(os.Getenv("KEY"), os.Getenv("SECRET"))
	defer cli.Close()

	balances, err := cli.Balances().
		Do(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, balances)

	t.Log(pretty.Sprint(balances))

	hasPositiveBalance := false
	for _, balance := range balances {
		require.GreaterOrEqual(t, balance, 0.0)
		if balance > 0.0 {
			hasPositiveBalance = true
		}
	}
	require.Truef(t, hasPositiveBalance, "No asset with a positive balance")
}
