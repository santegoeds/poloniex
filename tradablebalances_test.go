package poloniex_test

import (
	"context"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/santegoeds/poloniex"
	"github.com/stretchr/testify/require"
)

func TestTradableAccountBalances(t *testing.T) {
	ctx := context.TODO()
	cli := poloniex.New(os.Getenv("KEY"), os.Getenv("SECRET"))

	balances, err := cli.TradableAccountBalances().
		Do(ctx)
	require.NoError(t, err)

	t.Log(pretty.Sprint(balances))
}
