package ws

import (
	"encoding/json"
	"fmt"

	"github.com/santegoeds/poloniex/channel"
	"github.com/santegoeds/poloniex/errors"
	"github.com/santegoeds/poloniex/message"
)

type Decoder struct {
	subscription *Subscription
	bookDecoder  *PriceAggregatedBookDecoder
}

func NewDecoder(subscription *Subscription) *Decoder {
	return &Decoder{
		subscription: subscription,
	}
}

func (d *Decoder) Decode(data []json.RawMessage) error {
	switch len(data) {
	case 0:
		return fmt.Errorf("unexpected short message: %w", errors.ErrBadResponse)

	case 2:
		// Confirmation message without payload
		return nil

	case 1, 3:
		// Heartbeat or message with payload
		break

	default:
		return fmt.Errorf("%w: unexpected long message", errors.ErrBadResponse)
	}

	var channelID int
	if err := json.Unmarshal(data[0], &channelID); err != nil {
		return fmt.Errorf("%w: failed to unmarshal channel id", err)
	}

	if channelID == int(channel.Heartbeat) {
		d.subscription.HandleHeartbeat()
		return nil
	}

	switch channel.Channel(channelID) {
	case channel.ExchangeVolume24h, channel.TickerData:
		return d.decodeSimpleMessage(channelID, data)

	case channel.AccountNotifications:
		return d.decodeAccountNotifications(channelID, data)
	}
	return d.decodePriceAggregatedBook(channelID, data)
}

func (d *Decoder) decodeSimpleMessage(channelID int, data []json.RawMessage) error {
	payload := make([]json.RawMessage, 0)
	if err := json.Unmarshal(data[2], &payload); err != nil {
		return fmt.Errorf("%w: failed to unmarshal message payload", err)
	}

	switch channel.Channel(channelID) {
	case channel.ExchangeVolume24h:
		d.subscription.HandleMessage(message.Message{
			ChannelID: channelID,
			Type:      message.ExchangeVolume24h,
			Data:      payload,
		})
		return nil

	case channel.TickerData:
		d.subscription.HandleMessage(message.Message{
			ChannelID: channelID,
			Type:      message.TickerData,
			Data:      payload,
		})
		return nil
	}
	return fmt.Errorf("%w: Unexpected channelID %d", errors.ErrBadResponse, channelID)
}

func (d *Decoder) decodeAccountNotifications(channelID int, data []json.RawMessage) error {
	accountNotificationsDecoder := AccountNotificationsDecoder{Subscription: d.subscription}
	return accountNotificationsDecoder.Decode(channelID, data)
}

func (d *Decoder) decodePriceAggregatedBook(channelID int, data []json.RawMessage) error {
	bookDecoder := PriceAggregatedBookDecoder{d.subscription}
	return bookDecoder.Decode(channelID, data)
}
