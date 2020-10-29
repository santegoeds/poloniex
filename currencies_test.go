package poloniex_test

import (
	"context"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"

	"github.com/santegoeds/poloniex"
)

func TestCurrencies(t *testing.T) {
	ctx := context.TODO()
	cli := poloniex.New(os.Getenv("KEY"), os.Getenv("SECRET"))
	defer cli.Close()

	currencies, err := cli.Currencies().
		Do(ctx)
	require.NoError(t, err)

	t.Log(pretty.Sprint(currencies))

	require.NotEmpty(t, currencies)
	for curName, cur := range currencies {
		require.NotEmpty(t, curName)
		require.NotZero(t, cur.ID)
		require.NotEmpty(t, cur.Name)
		require.GreaterOrEqual(t, cur.TxFee, 0.0)
		require.GreaterOrEqual(t, cur.MinConf, 0)
		require.Contains(t, []int{0, 1}, cur.Disabled)
		require.Contains(t, []int{0, 1}, cur.Delisted)
		require.Contains(t, []int{0, 1}, cur.Frozen)
	}
}
