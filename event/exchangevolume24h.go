package event

import (
	channel2 "github.com/santegoeds/poloniex/channel"
	"strconv"
	"time"

	"github.com/santegoeds/poloniex/api/decoder"
	"github.com/santegoeds/poloniex/message"
)

const (
	minuteTimeFormat = "2006-01-02 15:04"
)

var _ Event = &ExchangeVolume24h{}

type ExchangeVolume24h struct {
	Date        string
	UsersOnline int64
	Volume24h   map[string]float64
}

func (ev *ExchangeVolume24h) ChannelID() int {
	return int(channel2.ExchangeVolume24h)
}

func (ev *ExchangeVolume24h) Type() message.Type {
	return message.ExchangeVolume24h
}

func (ev *ExchangeVolume24h) DateTime() time.Time {
	ts, _ := time.Parse(minuteTimeFormat, ev.Date)
	return ts
}

func (ev *ExchangeVolume24h) Unmarshal(msg message.Message) error {
	volume24h := make(map[string]string)
	err := decoder.Unmarshal(
		msg.Data,
		&ev.Date,
		&ev.UsersOnline,
		&volume24h,
	)
	if err != nil {
		return err
	}
	ev.Volume24h = make(map[string]float64)
	for currency, volumeStr := range volume24h {
		volume, err := strconv.ParseFloat(volumeStr, 64)
		if err != nil {
			return err
		}
		ev.Volume24h[currency] = volume
	}
	return nil
}
