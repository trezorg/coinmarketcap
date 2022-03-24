package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/trezorg/coinmarketcap"
)

func convert(args []string, callback func(prices coinmarketcap.Prices)) error {
	amount, err := coinmarketcap.NewAmount(args[0])
	if err != nil {
		return err
	}
	fromCurrencySymbol, err := coinmarketcap.NewCurrencySymbol(args[1])
	if err != nil {
		return err
	}
	toCurrencySymbol, err := coinmarketcap.NewCurrencySymbol(args[2])
	if err != nil {
		return err
	}
	req, err := coinmarketcap.NewPriceConversionRequest(amount, fromCurrencySymbol, toCurrencySymbol)
	if err != nil {
		return err
	}
	ctx, done := context.WithTimeout(context.Background(), time.Second*5)
	defer done()

	repository := coinmarketcap.CoinMarketRepository{}
	service := coinmarketcap.CoinMarketService{Repository: repository}
	prices, err := service.ConvertionRequest(ctx, req)
	if err != nil {
		return err
	}
	callback(prices)
	return nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:  "converter AMOUNT FROM_CURRENCY TO_CURRENCY",
		Long: "Convert fiat and cryptocurrencies. Example: ./converter 123.45 USD BTC",
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if err := convert(args, func(prices coinmarketcap.Prices) {
				for _, p := range prices {
					fmt.Printf("%0.2f\n", p.Price)
				}
			}); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
