package public

import (
	"context"
	"encoding/json"
	"io"
	"strconv"

	"github.com/santegoeds/poloniex/api/decoder"
)

type Order struct {
	Price float64
	Size  float64
}

type Book struct {
	Asks       []Order
	Bids       []Order
	Frozen     string
	SequenceNr int
}

func (b *Book) IsFrozen() bool {
	return b.Frozen == "1"
}

type OrderbookRequest struct {
	client       *Client
	currencyPair string
	depth        int64
}

func NewOrderbookRequest(client *Client, currencyPair string) *OrderbookRequest {
	return &OrderbookRequest{
		client:       client,
		currencyPair: currencyPair,
	}
}

func (r *OrderbookRequest) Depth(depth int64) *OrderbookRequest {
	r.depth = depth
	return r
}

func (r *OrderbookRequest) Do(ctx context.Context) (map[string]Book, error) {
	keyAndValues := []string{
		"command", "returnOrderBook",
		"currencyPair", r.currencyPair,
	}
	if r.depth > 0 {
		keyAndValues = append(keyAndValues, "depth", strconv.FormatInt(r.depth, 10))
	}
	rsp, err := r.client.do(ctx, keyAndValues...)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if r.currencyPair == "all" {
		return decodeAll(rsp.Body)
	}

	objData, err := decoder.DecodeObject(rsp.Body)
	if err != nil {
		return nil, err
	}
	b := Book{}
	if err = decodeCurrencyPair(objData, &b); err != nil {
		return nil, err
	}
	return map[string]Book{
		r.currencyPair: b,
	}, nil
}

func decodeAll(r io.Reader) (map[string]Book, error) {
	objData, err := decoder.DecodeObject(r)
	if err != nil {
		return nil, err
	}
	books := make(map[string]Book)
	for currencyPair, bookData := range objData {
		bookObj := make(map[string]json.RawMessage)
		if err = json.Unmarshal(bookData, &bookObj); err != nil {
			return nil, err
		}
		b := Book{}
		if err = decodeCurrencyPair(bookObj, &b); err != nil {
			return nil, err
		}
		books[currencyPair] = b
	}
	return books, nil
}

type orderTuple struct {
	Order
}

func (ot *orderTuple) UnmarshalJSON(data []byte) error {
	// Each entry in the order book consists of tuple [price, size] where the
	// price is encoded as a string and the size as a float
	tup := [2]interface{}{
		&decoder.Float64{Value: &ot.Price},
		&ot.Size,
	}
	return json.Unmarshal(data, &tup)
}

func decodeCurrencyPair(data map[string]json.RawMessage, b *Book) error {
	var asks []orderTuple
	var bids []orderTuple

	err := decoder.Unmarshal(
		[]json.RawMessage{
			data["asks"],
			data["bids"],
			data["isFrozen"],
			data["seq"],
		},
		&asks,
		&bids,
		&b.Frozen,
		&b.SequenceNr,
	)
	if err != nil {
		return err
	}

	for _, tup := range asks {
		b.Asks = append(b.Asks, tup.Order)
	}
	for _, tup := range bids {
		b.Bids = append(b.Bids, tup.Order)
	}

	return nil
}
