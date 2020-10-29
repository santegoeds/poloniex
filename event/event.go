package event

import (
	"github.com/santegoeds/poloniex/message"
)

type Event interface {
	ChannelID() int
	Type() message.Type
	Unmarshal(msg message.Message) error
}

const (
	timeFormat = "2006-01-02 15:04:05"
)
