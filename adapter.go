package coinmarketcap

import "context"

type Repository interface {
	ConvertionRequest(ctx context.Context, request priceConversionRequest) (Prices, error)
}

type Service interface {
	ConvertionRequest(ctx context.Context, request priceConversionRequest) (Prices, error)
}

type CoinMarketRepository struct {
}

func (ca CoinMarketRepository) ConvertionRequest(ctx context.Context, request priceConversionRequest) (Prices, error) {
	req, err := coinMarketConvertionRequest(request)
	if err != nil {
		return nil, err
	}
	resp, err := makeConversionRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.prices(), err
}

type CoinMarketService struct {
	Repository Repository
}

func (cs CoinMarketService) ConvertionRequest(ctx context.Context, request priceConversionRequest) (Prices, error) {
	return cs.Repository.ConvertionRequest(ctx, request)
}
