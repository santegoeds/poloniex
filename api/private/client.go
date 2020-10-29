package private

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/santegoeds/poloniex/api/decoder"
	"github.com/santegoeds/poloniex/errors"
)

const (
	Endpoint = "https://poloniex.com/tradingApi"

	timeFormat = "2006-01-02 15:04:05"
)

type Client struct {
	HttpClient *http.Client
	Endpoint   string
	key        string
	secret     string
}

func New(key, secret string) *Client {
	return &Client{
		HttpClient: http.DefaultClient,
		Endpoint:   Endpoint,
		key:        key,
		secret:     secret,
	}
}

//TODO:
// func (c *Client) Withdraw()
// func (c *Client) TransferBalance()
// func (c *Client) FeeInfo()
// func (c *Client) AvailableAccountBalances()
// func (c *Client) TransferBalance()
// func (c *Client) MarginAccountSummary()
// func (c *Client) MarginBuy()
// func (c *Client) MarginSell()
// func (c *Client) GetMarginPosition()
// func (c *Client) CloseMarginPosition()
// func (c *Client) CreateLoanOffer()
// func (c *Client) CancelLoanOffer()
// func (c *Client) OpenLoanOffers()
// func (c *Client) ActiveLoans()
// func (c *Client) LendingHistory()
// func (c *Client) ToggleAutoRenew()

func (c *Client) do(ctx context.Context, keyAndValue ...string) (*http.Response, error) {
	if len(keyAndValue)%2 != 0 {
		return nil, fmt.Errorf("uneven number of key/value arguments: %w", errors.ErrBadRequest)
	}

	formValues := make(url.Values)
	for idx := 0; idx < len(keyAndValue); idx += 2 {
		formValues.Set(keyAndValue[idx], keyAndValue[idx+1])
	}
	formValues.Set("nonce", strconv.FormatInt(time.Now().UnixNano(), 10))
	formData := formValues.Encode()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.Endpoint,
		strings.NewReader(formData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request for %s: %w", c.Endpoint, err)
	}

	hasher := hmac.New(sha512.New, []byte(c.secret))
	if _, err := hasher.Write([]byte(formData)); err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Key", c.key)
	req.Header.Set("Sign", hex.EncodeToString(hasher.Sum(nil)))

	rsp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// The caller is responsible for closing `rsp.Body`
	if rsp.StatusCode < 400 {
		return rsp, nil
	}
	defer rsp.Body.Close()

	if _, err = decoder.DecodeObject(rsp.Body); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("status code %d without message: %w", rsp.StatusCode, errors.ErrBadRequest)
}
