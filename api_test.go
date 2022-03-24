package coinmarketcap

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	response = []byte(`
	{
		"status": {
			"timestamp": "2022-03-23T19:45:25.786Z",
			"error_code": 0,
			"error_message": null,
			"elapsed": 0,
			"credit_count": 2,
			"notice": null
		},
		"data": {
			"BTC": {
				"symbol": "BTC",
				"id": "vp11kfj8iao",
				"name": "nxa27ggshk",
				"amount": 101,
				"last_updated": "2022-03-23T19:45:25.786Z",
				"quote": {
					"USD": {
						"price": 9093,
						"last_updated": "2022-03-23T19:45:25.786Z"
					},
					"ETH": {
						"price": 9093,
						"last_updated": "2022-03-23T19:45:25.786Z"
					}
				}
			}
		}
	}
	`)
)

func TestDeserializeConversionResponse(t *testing.T) {
	item := &priceConversionCoinMarketCapResponse{}
	err := json.Unmarshal(response, item)
	require.NoError(t, err)
	assert.Equal(t, len(item.Data), 1)
	assert.Equal(t, len(item.Data["BTC"].Quote), 2)
	assert.Equal(t, item.Status.ErrorMessage, "")
	assert.Equal(t, item.Status.ErrorCode, 0)
}

func TestConversionResponseGetPrice(t *testing.T) {
	item := &priceConversionCoinMarketCapResponse{}
	err := json.Unmarshal(response, item)
	require.NoError(t, err)
	prices := item.prices()
	price := prices.ForCurrency(USD)
	assert.NotNil(t, price)
	assert.Equal(t, float64((*price).Price), float64(9093))
}
