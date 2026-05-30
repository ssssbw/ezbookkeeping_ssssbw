package marketdata

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

var (
	Container = &MarketDataProviderContainer{}
)

type MarketDataProviderContainer struct {
	primary   MarketDataProvider
	fallback  MarketDataProvider
}

func InitializeMarketDataSource(config *settings.Config) error {
	switch config.MarketDataSource {
	case settings.AkshareMarketDataSource:
		Container.primary = NewAkshareMarketDataProvider()
		Container.fallback = NewEastMoneyMarketDataProvider()
	case settings.EastMoneyMarketDataSource:
		Container.primary = NewEastMoneyMarketDataProvider()
		Container.fallback = NewAkshareMarketDataProvider()
	default:
		return errs.ErrInvalidMarketDataSource
	}

	return nil
}

func (c *MarketDataProviderContainer) GetLatestPrice(assetCode string, market string) (*MarketDataResult, error) {
	if c.primary == nil {
		return nil, errs.ErrSystemIsBusy
	}

	ctx := core.NewNullContext()
	result, err := c.primary.GetLatestPrice(ctx, assetCode, market)
	if err != nil {
		log.Warnf(ctx, "[marketdata.GetLatestPrice] primary datasource: %s, failed for %s: %s, trying fallback", c.primary.GetDataSourceName(), assetCode, err.Error())
		if c.fallback != nil {
			result, err = c.fallback.GetLatestPrice(ctx, assetCode, market)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &MarketDataResult{
		AssetCode: assetCode,
		Market:    market,
		Data:      result,
	}, nil
}

func (c *MarketDataProviderContainer) GetRealtimeEstimate(assetCode string, market string) (*MarketDataResult, error) {
	if c.primary == nil {
		return nil, errs.ErrSystemIsBusy
	}

	ctx := core.NewNullContext()
	result, err := c.primary.GetRealtimeEstimate(ctx, assetCode, market)
	if err != nil {
		if c.fallback != nil {
			result, err = c.fallback.GetRealtimeEstimate(ctx, assetCode, market)
		}
	}

	if result == nil && err == nil {
		return nil, nil
	}

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
	if c.primary == nil {
		return nil, errs.ErrSystemIsBusy
	}

	ctx := core.NewNullContext()
	dataList, err := c.primary.GetHistoricalPrices(ctx, assetCode, market, startTime, endTime)
	if err != nil && c.fallback != nil {
		dataList, err = c.fallback.GetHistoricalPrices(ctx, assetCode, market, startTime, endTime)
	}
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
	if c.primary == nil {
		return nil, errs.ErrSystemIsBusy
	}

	ctx := core.NewNullContext()
	names, err := c.primary.GetAllFundNames(ctx)
	if err != nil && c.fallback != nil {
		names, err = c.fallback.GetAllFundNames(ctx)
	}
	return names, err
}

type MarketDataResult struct {
	AssetCode string
	Market    string
	Data      interface{}
}
