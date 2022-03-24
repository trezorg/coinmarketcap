package coinmarketcap

import (
	"fmt"
	"strings"
)

type CurrencySymbol string

const (
	BTC CurrencySymbol = CurrencySymbol("BTC")
	USD CurrencySymbol = CurrencySymbol("USD")
	EUR CurrencySymbol = CurrencySymbol("EUR")
	ETH CurrencySymbol = CurrencySymbol("ETH")
	LTC CurrencySymbol = CurrencySymbol("LTC")
	GBP CurrencySymbol = CurrencySymbol("GBP")
)

func (cs CurrencySymbol) String() string {
	return string(cs)
}

func NewCurrencySymbol(s string) (CurrencySymbol, error) {
	res := CurrencySymbol(s)
	if !res.IsValid() {
		return CurrencySymbol(""), fmt.Errorf("Non valid currency symbol: %s, Valid currency symbols: %s", s, availableSymbolsString())
	}
	return res, nil
}

func (cs CurrencySymbol) IsValid() bool {
	switch cs {
	case BTC, USD, EUR, ETH, LTC, GBP:
		return true
	default:
		return false
	}
}

func availableSymbols() []CurrencySymbol {
	return []CurrencySymbol{BTC, USD, EUR, ETH, LTC, GBP}
}

func availableSymbolsString() string {
	var sb strings.Builder
	symbols := availableSymbols()
	for idx, symbol := range symbols {
		sb.WriteString(symbol.String())
		if idx < len(symbols)-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}
