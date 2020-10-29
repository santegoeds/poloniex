package ws

import (
	"encoding/json"
	"fmt"
	"github.com/santegoeds/poloniex/errors"
	"github.com/santegoeds/poloniex/message"
)

const (
	accountNotificationTypeOrderPending  = "p"
	accountNotificationTypeBalanceUpdate = "b"
	accountNotificationTypeLimitOrder    = "n"
	accountNotificationTypeOrderUpdate   = "o"
	accountNotificationTypeMarginUpdate  = "m"
	accountNotificationTypeTrade         = "t"
	accountNotificationTypeOrderKilled   = "k"
)

type AccountNotificationsDecoder struct {
	Subscription *Subscription
}

func (d *AccountNotificationsDecoder) Decode(channelID int, data []json.RawMessage) error {
	payload := make([]json.RawMessage, 0)
	if err := json.Unmarshal(data[2], &payload); err != nil {
		return fmt.Errorf("%w: failed to marshal message payload", err)
	}

	for _, data := range data {
		payload := make([]json.RawMessage, 0)
		if err := json.Unmarshal(data, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal account notification: %w", err)
		}

		var typeFlag string
		if err := json.Unmarshal(payload[0], &typeFlag); err != nil {
			return fmt.Errorf("%w: failed to unmarshal account notification type flag", err)
		}
		payload = payload[1:]

		if len(payload) == 0 {
			return fmt.Errorf("%w: account notification %s has no payload", errors.ErrBadResponse, typeFlag)
		}

		switch typeFlag {
		case accountNotificationTypeOrderPending:
			d.Subscription.HandleMessage(message.Message{
				ChannelID: channelID,
				Type:      message.AccountOrderPending,
				Data:      payload,
			})

		case accountNotificationTypeBalanceUpdate:
			d.Subscription.HandleMessage(message.Message{
				ChannelID: channelID,
				Type:      message.AccountBalanceUpdate,
				Data:      payload,
			})

		case accountNotificationTypeLimitOrder:
			d.Subscription.HandleMessage(message.Message{
				ChannelID: channelID,
				Type:      message.AccountLimitOrder,
				Data:      payload,
			})

		case accountNotificationTypeOrderUpdate:
			d.Subscription.HandleMessage(message.Message{
				ChannelID: channelID,
				Type:      message.AccountOrderUpdate,
				Data:      payload,
			})

		case accountNotificationTypeMarginUpdate:
			d.Subscription.HandleMessage(message.Message{
				ChannelID: channelID,
				Type:      message.AccountMarginPositionUpdate,
				Data:      payload,
			})

		case accountNotificationTypeTrade:
			d.Subscription.HandleMessage(message.Message{
				ChannelID: channelID,
				Type:      message.AccountTrade,
				Data:      payload,
			})

		case accountNotificationTypeOrderKilled:
			d.Subscription.HandleMessage(message.Message{
				ChannelID: channelID,
				Type:      message.AccountOrderKilled,
				Data:      payload,
			})

		default:
			return fmt.Errorf("%w: unexpected typeFlag %s", errors.ErrBadResponse, typeFlag)
		}
	}
	return nil
}
