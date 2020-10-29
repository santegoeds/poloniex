package channel

type Channel int

const (
	AccountNotifications Channel = 1000
	TickerData           Channel = 1002
	Heartbeat            Channel = 1010
	ExchangeVolume24h    Channel = 1003
)
