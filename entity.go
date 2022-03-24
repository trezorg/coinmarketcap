package coinmarketcap

import (
	"fmt"
	"strconv"
	"strings"
)

type Amount float64

func (a Amount) String() string {
	return fmt.Sprintf("%f", a)
}

func NewAmount(value string) (Amount, error) {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return Amount(0), err
	}
	return Amount(f), nil
}

type Currency struct {
	symbol CurrencySymbol
}

func (c Currency) String() string {
	return string(c.symbol)
}

func newCurrency(symbol CurrencySymbol) Currency {
	return Currency{symbol: symbol}
}

type Currencies []Currency

func (cs Currencies) String() string {
	var sb strings.Builder
	for idx, c := range cs {
		sb.WriteString(c.String())
		if idx < len(cs)-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

type priceConversionRequest struct {
	amount       Amount
	fromCurrency Currency
	toCurrency   Currencies
}

func NewPriceConversionRequest(amount Amount, fromCurrencySymbol CurrencySymbol, toCurrencySymbol ...CurrencySymbol) (priceConversionRequest, error) {
	req := priceConversionRequest{}
	var toCurrency Currencies
	for _, symbol := range toCurrencySymbol {
		if !symbol.IsValid() {
			return req, fmt.Errorf("Non valid currency symbol: %s. Valid currency symbols: %s", symbol, availableSymbolsString())
		}
		toCurrency = append(toCurrency, newCurrency(symbol))
	}
	if !fromCurrencySymbol.IsValid() {
		return req, fmt.Errorf("Non valid currency symbol: %s. Valid currency symbols: %s", fromCurrencySymbol, availableSymbolsString())
	}
	req.amount = amount
	req.fromCurrency = newCurrency(fromCurrencySymbol)
	req.toCurrency = toCurrency
	return req, nil
}

type Price struct {
	Price    Amount
	Currency Currency
}

func (p Price) String() string {
	return fmt.Sprintf("%s: %s", p.Currency, p.Price)
}

type Prices []Price

func (p Prices) ForCurrency(symbol CurrencySymbol) *Price {
	for _, p := range p {
		if p.Currency.symbol == symbol {
			return &p
		}
	}
	return nil
}

func (p Prices) Price() *Price {
	for _, p := range p {
		return &p
	}
	return nil
}
