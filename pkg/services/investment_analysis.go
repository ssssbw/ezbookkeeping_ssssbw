package services

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// InvestmentAnalysisService handles investment analysis and aggregation
type InvestmentAnalysisService struct {
	ServiceUsingDB
}

// InvestmentAnalysis is the singleton instance of InvestmentAnalysisService
var InvestmentAnalysis = &InvestmentAnalysisService{
	ServiceUsingDB: ServiceUsingDB{
		container: datastore.Container,
	},
}

// GetHoldings returns computed holding info for all active user assets
func (s *InvestmentAnalysisService) GetHoldings(c core.Context, uid int64) ([]*models.InvestmentHoldingInfo, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	// 1. Get all active user assets
	var userAssets []*models.UserAsset
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND is_active=?", uid, false, true).Find(&userAssets)
	if err != nil {
		return nil, err
	}

	if len(userAssets) == 0 {
		return []*models.InvestmentHoldingInfo{}, nil
	}

	var holdings []*models.InvestmentHoldingInfo

	for _, ua := range userAssets {
		// 2a. Get asset info
		asset := &models.Asset{}
		has, err := s.UserDataDB(0).NewSession(c).ID(ua.AssetId).Get(asset)
		if err != nil || !has {
			continue
		}

		// 2b. Get all non-deleted transactions for this asset, ordered by trade_time ASC
		var transactions []*models.InvestmentTransaction
		err = s.UserDataDB(uid).NewSession(c).
			Where("uid=? AND deleted=? AND asset_id=?", uid, false, ua.AssetId).
			OrderBy("trade_time asc").
			Find(&transactions)
		if err != nil || len(transactions) == 0 {
			continue
		}

		// 2c. Walk through transactions maintaining running totals
		var totalQuantity int64
		var totalCost     int64
		accountIdSet := make(map[int64]bool)
		var lastAccountId int64

		for _, tx := range transactions {
			accountIdSet[tx.AccountId] = true
			lastAccountId = tx.AccountId

			switch tx.Type {
			case models.INVESTMENT_TRANSACTION_TYPE_BUY:
				totalQuantity += tx.Quantity
				totalCost += tx.Amount + tx.Fee

			case models.INVESTMENT_TRANSACTION_TYPE_SELL:
				if totalQuantity > 0 {
					sellCost := totalCost * tx.Quantity / totalQuantity
					totalQuantity -= tx.Quantity
					totalCost -= sellCost
				}

			case models.INVESTMENT_TRANSACTION_TYPE_DIVIDEND_REINVEST:
				totalQuantity += tx.Quantity
				totalCost += tx.Amount

			case models.INVESTMENT_TRANSACTION_TYPE_SPLIT:
				totalQuantity += tx.Quantity
				// totalCost unchanged for splits

			case models.INVESTMENT_TRANSACTION_TYPE_CONVERSION_OUT:
				if totalQuantity > 0 {
					sellCost := totalCost * tx.Quantity / totalQuantity
					totalQuantity -= tx.Quantity
					totalCost -= sellCost
				}

			case models.INVESTMENT_TRANSACTION_TYPE_CONVERSION_IN:
				totalQuantity += tx.Quantity
				totalCost += tx.Amount + tx.Fee

			case models.INVESTMENT_TRANSACTION_TYPE_DIVIDEND_CASH:
				// No holding change
			}
		}

		// 2d. Skip if no holdings
		if totalQuantity <= 0 {
			continue
		}

		// 2e. Get latest market price
		marketData := &models.MarketData{}
		_, _ = s.UserDataDB(uid).NewSession(c).
			Where("asset_id=?", ua.AssetId).
			OrderBy("date desc").
			Limit(1).
			Get(marketData)

		currentPrice := marketData.Price

		// 2f. Compute derived metrics
		marketValue := totalQuantity * currentPrice / 10000
		unrealizedPnl := marketValue - totalCost
		var returnRate int64
		if totalCost > 0 {
			returnRate = unrealizedPnl * 10000 / totalCost
		}

		// 2g. Determine accountId
		accountId := lastAccountId
		if len(accountIdSet) == 1 {
			for id := range accountIdSet {
				accountId = id
				break
			}
		}

		// 2h. Build holding info
		var avgCostPrice int64
		if totalQuantity > 0 {
			avgCostPrice = totalCost * 10000 / totalQuantity
		}

		holding := &models.InvestmentHoldingInfo{
			AssetId:       ua.AssetId,
			AssetCode:     asset.Code,
			AssetName:     asset.Name,
			Category:      asset.Category,
			Currency:      asset.Currency,
			Market:        asset.Market,
			AccountId:     accountId,
			Quantity:      totalQuantity,
			AvgCostPrice:  avgCostPrice,
			TotalCost:     totalCost,
			CurrentPrice:  currentPrice,
			MarketValue:   marketValue,
			UnrealizedPnl: unrealizedPnl,
			ReturnRate:    returnRate,
		}

		holdings = append(holdings, holding)
	}

	return holdings, nil
}

// GetOverview returns aggregated portfolio summary
func (s *InvestmentAnalysisService) GetOverview(c core.Context, uid int64) (*models.InvestmentOverviewResponse, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	holdings, err := s.GetHoldings(c, uid)
	if err != nil {
		return nil, err
	}

	var totalInvestment  int64
	var totalMarketValue int64

	// Group by category
	categoryValues := make(map[string]int64)

	for _, h := range holdings {
		totalInvestment += h.TotalCost
		totalMarketValue += h.MarketValue
		categoryValues[string(h.Category)] += h.MarketValue
	}

	totalUnrealizedPnl := totalMarketValue - totalInvestment
	var totalReturnRate int64
	if totalInvestment > 0 {
		totalReturnRate = totalUnrealizedPnl * 10000 / totalInvestment
	}

	// Build allocation items
	allocations := make([]*models.InvestmentAllocationItem, 0, len(categoryValues))
	for category, value := range categoryValues {
		var percentage int64
		if totalMarketValue > 0 {
			percentage = value * 10000 / totalMarketValue
		}
		allocations = append(allocations, &models.InvestmentAllocationItem{
			Category:   category,
			Value:      value,
			Percentage: percentage,
		})
	}

	return &models.InvestmentOverviewResponse{
		TotalInvestment:    totalInvestment,
		TotalMarketValue:   totalMarketValue,
		TotalUnrealizedPnl: totalUnrealizedPnl,
		TotalReturnRate:    totalReturnRate,
		Allocations:        allocations,
	}, nil
}
