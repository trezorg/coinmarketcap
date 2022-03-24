package coinmarketcap

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func coinMarketConvertionRequest(pcr priceConversionRequest) (*http.Request, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = priceConvertURL
	query := u.Query()
	query.Add("amount", pcr.amount.String())
	query.Add("symbol", pcr.fromCurrency.String())
	query.Add("convert", pcr.toCurrency.String())
	u.RawQuery = query.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set(authHeader, token)
	return req, nil
}

func isRetryableError(err error) bool {
	cme, ok := err.(coinMarketCapError)
	if !ok {
		return false
	}
	if cme.errorCode == retryErrorCode {
		return true
	}
	return false
}

func readBody(resp *http.Response) ([]byte, error) {
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func makeHttpRequest(ctx context.Context, req *http.Request) ([]byte, error) {
	req = req.WithContext(ctx)
	client := prepareClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func retryWait(max, min int) {
	time.Sleep(time.Second * time.Duration(rand.Intn(max-min+1)+min))
}

func makeConversionRequest(ctx context.Context, req *http.Request) (priceConversionCoinMarketCapResponse, error) {
	item := &priceConversionCoinMarketCapResponse{}
	body, err := makeHttpRequest(ctx, req)
	if err != nil {
		return priceConversionCoinMarketCapResponse{}, err
	}
	err = json.Unmarshal(body, item)
	if err != nil {
		return *item, err
	}
	if err = item.error(); err != nil {
		if isRetryableError(err) {
			log.Println(err)
			retryWait(3, 1)
			return makeConversionRequest(ctx, req)
		}
		return *item, err
	}
	return *item, nil
}

func prepareClient() *http.Client {
	netTransport := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: netTransport,
	}
	return client
}

type coinMarketCapError struct {
	errorCode    int
	errorMessage string
}

func (e coinMarketCapError) Error() string {
	return fmt.Sprintf("Message: %s. Code: %d", e.errorMessage, e.errorCode)
}

type price struct {
	Price       float64 `json:"price"`
	LastUpdated string  `json:"last_updated"`
}

type status struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Elapsed      int    `json:"elapsed"`
	CreditCount  int    `json:"credit_count"`
	Timestamp    string `json:"timestamp"`
}

type priceConversionData struct {
	Symbol      CurrencySymbol   `json:"symbol"`
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Amount      float64          `json:"amount"`
	LastUpdated string           `json:"last_updated"`
	Quote       map[string]price `json:"quote"`
}

type priceConversionCoinMarketCapResponse struct {
	Data   map[string]priceConversionData `json:"data"`
	Status status                         `json:"status"`
}

func (r priceConversionCoinMarketCapResponse) error() error {
	if r.Status.ErrorCode != 0 {
		return coinMarketCapError{errorMessage: r.Status.ErrorMessage, errorCode: r.Status.ErrorCode}
	}
	return nil
}

func (r priceConversionCoinMarketCapResponse) prices() Prices {
	var result Prices
	for _, value := range r.Data {
		for currency, value := range value.Quote {
			result = append(result, Price{Price: Amount(value.Price), Currency: newCurrency(CurrencySymbol(currency))})
		}
	}
	return result
}
