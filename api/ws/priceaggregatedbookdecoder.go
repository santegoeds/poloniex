package ws

import (
	"encoding/json"
	"fmt"
	"github.com/santegoeds/poloniex/errors"
	"github.com/santegoeds/poloniex/message"
)

const (
	priceAggregatedBookTypeBook       = "i"
	priceAggregatedBookTypeBookUpdate = "o"
	priceAggregatedBookTypeTrade      = "t"
)

type PriceAggregatedBookDecoder struct {
	Subscription *Subscription
}

func (d *PriceAggregatedBookDecoder) Decode(channelID int, data []json.RawMessage) error {
	sequenceData, payloadData := data[1], data[2]

	msgsData := make([]json.RawMessage, 0)
	if err := json.Unmarshal(payloadData, &msgsData); err != nil {
		return fmt.Errorf("%w: failed to unmarshal payload array for priceaggregatedbook", err)
	}

	parts := make([]json.RawMessage, 0, 2)
	for _, msgData := range msgsData {
		if err := json.Unmarshal(msgData, &parts); err != nil {
			return fmt.Errorf("%w: failed to unmarshal payload array for priceaggregatedbook", err)
		}

		msgType, err := d.DecodeEventType(parts[0])
		if err != nil {
			return err
		}

		d.Subscription.HandleMessage(message.Message{
			ChannelID: channelID,
			Type:      msgType,
			Data:      append([]json.RawMessage{sequenceData}, parts[1:]...),
		})
	}

	return nil
}

func (d *PriceAggregatedBookDecoder) DecodeEventType(data json.RawMessage) (message.Type, error) {
	var eventType string
	if err := json.Unmarshal(data, &eventType); err != nil {
		return 0, err
	}

	switch eventType {
	case priceAggregatedBookTypeBook:
		return message.Book, nil

	case priceAggregatedBookTypeBookUpdate:
		return message.BookUpdate, nil

	case priceAggregatedBookTypeTrade:
		return message.BookTrade, nil
	}

	return 0, fmt.Errorf(
		"%w: unexpected event type %s for priceaggregatedbook",
		errors.ErrBadResponse, eventType,
	)
}
