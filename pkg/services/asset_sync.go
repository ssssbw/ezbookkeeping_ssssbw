package services

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/marketdata"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

type AssetSyncService struct {
	assets   *AssetService
	provider marketdata.AssetListProvider
}

type AssetSyncResult struct {
	FundsAdded  int64 `json:"fundsAdded"`
	StocksAdded int64 `json:"stocksAdded"`
	ETFsAdded   int64 `json:"etfsAdded"`
	TotalAdded  int64 `json:"totalAdded"`
}

var AssetSync = &AssetSyncService{
	assets:   Assets,
	provider: marketdata.NewEastMoneyAssetListProvider(),
}

func (s *AssetSyncService) SyncAllAssets(c core.Context) (*AssetSyncResult, error) {
	var result AssetSyncResult

	result.FundsAdded = s.syncFunds(c)

	// TODO: push2.eastmoney.com 在海外被墙，暂时跳过
	// result.StocksAdded = s.syncStocks(c)
	// result.ETFsAdded = s.syncETFs(c)

	result.TotalAdded = result.FundsAdded + result.StocksAdded + result.ETFsAdded

	log.Infof(c, "[asset_sync] sync completed: funds=%d, stocks=%d, etfs=%d, total=%d",
		result.FundsAdded, result.StocksAdded, result.ETFsAdded, result.TotalAdded)

	return &result, nil
}

func (s *AssetSyncService) syncFunds(c core.Context) int64 {
	metas, err := s.provider.ListFunds(c)
	if err != nil {
		log.Warnf(c, "[asset_sync] list funds failed: %s", err.Error())
		return 0
	}

	existing, err := s.assets.GetAllAssetCodeAndMarketMap(c, models.INVESTMENT_MARKET_CN)
	if err != nil {
		log.Warnf(c, "[asset_sync] get existing assets failed: %s", err.Error())
		return 0
	}

	var added int64
	// 增量更新
	for _, meta := range metas {
		if existing[meta.Code] {
			continue
		}

		asset := &models.Asset{
			Code:        meta.Code,
			Market:      models.INVESTMENT_MARKET_CN,
			Name:        meta.Name,
			Category:    models.AssetCategory(meta.Category),
			SubCategory: meta.SubCategory,
			Currency:    "CNY",
			Industry:    "other",
			Tags:        "[]",
		}

		if createErr := s.assets.CreateAsset(c, asset); createErr != nil {
			log.Warnf(c, "[asset_sync] create fund %s failed: %s", meta.Code, createErr.Error())
			continue
		}
		existing[meta.Code] = true
		added++
	}

	log.Infof(c, "[asset_sync] sync funds done: added %d / %d", added, len(metas))
	return added
}
