package marketdata

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

var (
	Container = &MarketDataProviderContainer{}
)

type MarketDataProviderContainer struct {
	current MarketDataProvider
}

func InitializeMarketDataSource(config *settings.Config) error {
	switch config.MarketDataSource {
	case settings.AkshareMarketDataSource:
		Container.current = NewAkshareMarketDataProvider()
	case settings.EastMoneyMarketDataSource:
		Container.current = NewEastMoneyMarketDataProvider()
	default:
		return errs.ErrInvalidMarketDataSource
	}

	return nil
}

func (c *MarketDataProviderContainer) GetLatestPrice(assetCode string, market string) (*MarketDataResult, error) {
	if c.current == nil {
		return nil, errs.ErrSystemIsBusy
	}

	ctx := core.NewNullContext()
	result, err := c.current.GetLatestPrice(ctx, assetCode, market)
	if err != nil {
		return nil, err
	}

	return &MarketDataResult{
		AssetCode: assetCode,
		Market:    market,
		Data:      result,
	}, nil
}

func (c *MarketDataProviderContainer) GetHistoricalPrices(assetCode string, market string, startTime int64, endTime int64) ([]*MarketDataResult, error) {
	if c.current == nil {
		return nil, errs.ErrSystemIsBusy
	}

	ctx := core.NewNullContext()
	dataList, err := c.current.GetHistoricalPrices(ctx, assetCode, market, startTime, endTime)
	if err != nil {
		return nil, err
	}

	var results []*MarketDataResult
	for _, data := range dataList {
		results = append(results, &MarketDataResult{
			AssetCode: assetCode,
			Market:    market,
			Data:      data,
		})
	}

	return results, nil
}

func (c *MarketDataProviderContainer) GetAllFundNames() (map[string]string, error) {
	if c.current == nil {
		return nil, errs.ErrSystemIsBusy
	}

	ctx := core.NewNullContext()
	return c.current.GetAllFundNames(ctx)
}

type MarketDataResult struct {
	AssetCode string
	Market    string
	Data      interface{}
}
