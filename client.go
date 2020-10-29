package poloniex

import (
	"github.com/santegoeds/poloniex/period"

	"github.com/santegoeds/poloniex/api/private"
	"github.com/santegoeds/poloniex/api/public"
	"github.com/santegoeds/poloniex/api/ws"
)

type Client struct {
	public  *public.Client
	private *private.Client
	ws      *ws.Client
}

type OrderSide = private.OrderSide

const (
	Buy  = private.Buy
	Sell = private.Sell
)

func New(key, secret string) *Client {
	return &Client{
		private: private.New(key, secret),
		public:  public.New(),
		ws:      ws.New(key, secret),
	}
}

func (c *Client) Close() {
	c.ws.Close()
}

func (c *Client) Currencies() *public.CurrenciesRequest {
	return public.NewCurrenciesRequest(c.public)
}

func (c *Client) Tickers() *public.TickerRequest {
	return public.NewTickerRequest(c.public)
}

func (c *Client) ChartData(currencyPair string, period period.Period) *public.ChartDataRequest {
	return public.NewChartDataRequest(c.public, currencyPair, period)
}

func (c *Client) OrderBook(currencyPair string) *public.OrderbookRequest {
	if currencyPair == "" {
		currencyPair = "all"
	}
	return public.NewOrderbookRequest(c.public, currencyPair)
}

func (c *Client) TradeHistory(currencyPair string) *public.TradeHistoryRequest {
	return public.NewTradeHistoryRequest(c.public, currencyPair)
}

func (c *Client) Volume24h() *public.Volume24hRequest {
	return public.NewVolume24hRequest(c.public)
}

func (c *Client) Buy(currencyPair string, rate float64, amount float64) *private.OrderRequest {
	return private.NewOrderRequest(c.private, Buy, currencyPair).
		Rate(rate).
		Amount(amount)
}

func (c *Client) Sell(currencyPair string, rate float64, amount float64) *private.OrderRequest {
	return private.NewOrderRequest(c.private, Sell, currencyPair).
		Rate(rate).
		Amount(amount)
}

func (c *Client) MoveOrder(orderNumber int64, rate float64) *private.MoveOrderRequest {
	return private.NewMoveOrderRequest(c.private, orderNumber, rate)
}

func (c *Client) Balances() *private.BalancesRequest {
	return private.NewBalancesRequest(c.private)
}

func (c *Client) CompleteBalances() *private.CompleteBalancesRequest {
	return private.NewCompleteBalancesRequest(c.private)
}

func (c *Client) OpenOrders(currencyPair string) *private.OpenOrdersRequest {
	if currencyPair == "" {
		currencyPair = "all"
	}
	return private.NewOpenOrdersRequest(c.private, currencyPair)
}

func (c *Client) OrderStatus(orderNumber int64) *private.OrderStatusRequest {
	return private.NewOrderStatusRequest(c.private, orderNumber)
}

func (c *Client) CancelOrder() *private.CancelOrderRequest {
	return private.NewCancelOrderRequest(c.private)
}

func (c *Client) CancelAllOrders(currencyPair string) *private.CancelAllOrdersRequest {
	return private.NewCancelAllOrdersRequest(c.private, currencyPair)
}

func (c *Client) TradableAccountBalances() *private.TradableBalancesRequest {
	return private.NewTradableBalancesRequest(c.private)
}

func (c *Client) Subscribe() *ws.SubscriptionRequest {
	return ws.NewSubscriptionRequest(c.ws)
}
