package services

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/marketdata"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// AssetSyncService synchronizes global assets from external data sources
type AssetSyncService struct {
	assets   *AssetService
	provider marketdata.AssetListProvider
}

// AssetSyncResult represents the result of an asset sync operation
type AssetSyncResult struct {
	FundsAdded  int64 `json:"fundsAdded"`
	StocksAdded int64 `json:"stocksAdded"`
	ETFsAdded   int64 `json:"etfsAdded"`
	TotalAdded  int64 `json:"totalAdded"`
}

// AssetSync is the package-level instance of AssetSyncService
var AssetSync = &AssetSyncService{
	assets:   Assets,
	provider: marketdata.NewEastMoneyAssetListProvider(),
}

// SyncAllAssets runs all three asset sync sub-tasks and returns the combined result
func (s *AssetSyncService) SyncAllAssets(c core.Context) (*AssetSyncResult, error) {
	var result AssetSyncResult
	var err error

	result.FundsAdded, err = s.syncFunds(c)
	if err != nil {
		log.Warnf(c, "[asset_sync] sync funds failed: %s", err.Error())
	}

	result.StocksAdded, err = s.syncStocks(c)
	if err != nil {
		log.Warnf(c, "[asset_sync] sync stocks failed: %s", err.Error())
	}

	result.ETFsAdded, err = s.syncETFs(c)
	if err != nil {
		log.Warnf(c, "[asset_sync] sync ETFs failed: %s", err.Error())
	}

	result.TotalAdded = result.FundsAdded + result.StocksAdded + result.ETFsAdded

	log.Infof(c, "[asset_sync] sync completed: funds=%d, stocks=%d, etfs=%d, total=%d",
		result.FundsAdded, result.StocksAdded, result.ETFsAdded, result.TotalAdded)

	return &result, nil
}

// syncFunds fetches all funds from the provider and inserts new ones
func (s *AssetSyncService) syncFunds(c core.Context) (int64, error) {
	metas, err := s.provider.ListFunds(c)
	if err != nil {
		return 0, err
	}

	var added int64
	for _, meta := range metas {
		_, err := s.assets.GetAssetByCodeAndMarket(c, meta.Code, models.InvestmentMarket(meta.Market))
		if err == nil {
			// already exists, skip
			continue
		}
		if err != errs.ErrInvestmentAssetNotFound {
			log.Warnf(c, "[asset_sync] failed to check fund %s: %s", meta.Code, err.Error())
			continue
		}

		asset := &models.Asset{
			Code:     meta.Code,
			Market:   models.InvestmentMarket(meta.Market),
			Name:     meta.Name,
			Category:    models.AssetCategory(meta.Category),
			SubCategory: meta.SubCategory,
			Currency:    "CNY",
			Industry: "other",
			Tags:     "[]",
		}

		if createErr := s.assets.CreateAsset(c, asset); createErr != nil {
			log.Warnf(c, "[asset_sync] failed to create fund %s: %s", meta.Code, createErr.Error())
			continue
		}
		added++
	}

	return added, nil
}

// syncStocks fetches all stocks from the provider and inserts new ones
func (s *AssetSyncService) syncStocks(c core.Context) (int64, error) {
	metas, err := s.provider.ListStocks(c)
	if err != nil {
		return 0, err
	}

	var added int64
	for _, meta := range metas {
		_, err := s.assets.GetAssetByCodeAndMarket(c, meta.Code, models.InvestmentMarket(meta.Market))
		if err == nil {
			continue
		}
		if err != errs.ErrInvestmentAssetNotFound {
			log.Warnf(c, "[asset_sync] failed to check stock %s: %s", meta.Code, err.Error())
			continue
		}

		asset := &models.Asset{
			Code:     meta.Code,
			Market:   models.InvestmentMarket(meta.Market),
			Name:     meta.Name,
			Category:    models.AssetCategory(meta.Category),
			SubCategory: meta.SubCategory,
			Currency:    "CNY",
			Industry: "other",
			Tags:     "[]",
		}

		if createErr := s.assets.CreateAsset(c, asset); createErr != nil {
			log.Warnf(c, "[asset_sync] failed to create stock %s: %s", meta.Code, createErr.Error())
			continue
		}
		added++
	}

	return added, nil
}

// syncETFs fetches all ETFs from the provider and inserts new ones
func (s *AssetSyncService) syncETFs(c core.Context) (int64, error) {
	metas, err := s.provider.ListETFs(c)
	if err != nil {
		return 0, err
	}

	var added int64
	for _, meta := range metas {
		_, err := s.assets.GetAssetByCodeAndMarket(c, meta.Code, models.InvestmentMarket(meta.Market))
		if err == nil {
			continue
		}
		if err != errs.ErrInvestmentAssetNotFound {
			log.Warnf(c, "[asset_sync] failed to check ETF %s: %s", meta.Code, err.Error())
			continue
		}

		asset := &models.Asset{
			Code:     meta.Code,
			Market:   models.InvestmentMarket(meta.Market),
			Name:     meta.Name,
			Category:    models.AssetCategory(meta.Category),
			SubCategory: meta.SubCategory,
			Currency:    "CNY",
			Industry: "other",
			Tags:     "[]",
		}

		if createErr := s.assets.CreateAsset(c, asset); createErr != nil {
			log.Warnf(c, "[asset_sync] failed to create ETF %s: %s", meta.Code, createErr.Error())
			continue
		}
		added++
	}

	return added, nil
}
