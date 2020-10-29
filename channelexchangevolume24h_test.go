package poloniex_test

import (
	"context"
	"github.com/santegoeds/poloniex/event"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"

	"github.com/santegoeds/poloniex"
	"github.com/santegoeds/poloniex/channel"
	"github.com/santegoeds/poloniex/message"
)

func TestExchangeVolume24h(t *testing.T) {
	ctx := context.TODO()
	cli := poloniex.New(os.Getenv("KEY"), os.Getenv("SECRET"))
	defer cli.Close()

	subscription, err := cli.Subscribe().
		Channel(channel.ExchangeVolume24h).
		Do(ctx)
	require.NoError(t, err)

	msgC, errC := subscription.Channels()

	select {
	case err := <-errC:
		t.Fatal(err)

	case msg := <-msgC:
		require.Equal(t, msg.ChannelID, int(channel.ExchangeVolume24h))
		require.Equal(t, msg.Type, message.ExchangeVolume24h)

		evt := event.ExchangeVolume24h{}
		err := evt.Unmarshal(msg)
		require.NoError(t, err)

		t.Log(pretty.Sprint(evt))

		require.NotZero(t, evt.Date)
		require.True(t, evt.UsersOnline > 0)
		require.NotEmpty(t, evt.Volume24h)
	}
}
