package public

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/santegoeds/poloniex/api/decoder"
	"github.com/santegoeds/poloniex/errors"
)

const (
	Endpoint = "https://poloniex.com/public"

	timeFormat = "2006-01-02 15:04:05"
)

type Client struct {
	HttpClient *http.Client
	Endpoint   string
}

func New() *Client {
	return &Client{
		HttpClient: http.DefaultClient,
		Endpoint:   Endpoint,
	}
}

func (c *Client) do(ctx context.Context, keyAndValue ...string) (*http.Response, error) {
	if len(keyAndValue)%2 != 0 {
		return nil, fmt.Errorf("%w: uneven number of key/value arguments", errors.ErrBadRequest)
	}

	URL, err := url.ParseRequestURI(c.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid endpoint %s", errors.ErrBadRequest, c.Endpoint)
	}
	query := make(url.Values)
	for idx := 0; idx < len(keyAndValue); idx += 2 {
		query.Set(keyAndValue[idx], keyAndValue[idx+1])
	}
	URL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create a request for %s", err, URL.String())
	}
	req.Header.Add("Accept", "application/json")

	rsp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	// Caller is responsible for closing `rsp.Body`
	if rsp.StatusCode < 400 {
		return rsp, nil
	}
	defer rsp.Body.Close()

	if _, err = decoder.DecodeObject(rsp.Body); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("%w: status code %d without message", errors.ErrBadRequest, rsp.StatusCode)
}
