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

// GetHoldings returns computed holding info for all active user assets, grouped by asset+account
func (s *InvestmentAnalysisService) GetHoldings(c core.Context, uid int64) ([]*models.InvestmentHoldingInfo, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	// 1. Get all active user assets
	var userAssets []*models.UserAsset
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND is_active=? AND is_watchlist=?", uid, false, true, false).Find(&userAssets)
	if err != nil {
		return nil, err
	}

	if len(userAssets) == 0 {
		return []*models.InvestmentHoldingInfo{}, nil
	}

	// 2. Batch load account names
	var accounts []*models.Account
	err = s.UserDataDB(uid).NewSession(c).
		Where("uid=? AND deleted=? AND category=?", uid, false, models.ACCOUNT_CATEGORY_INVESTMENT).
		Find(&accounts)
	if err != nil {
		return nil, err
	}

	accountNameMap := make(map[int64]string, len(accounts))
	for _, a := range accounts {
		accountNameMap[a.AccountId] = a.Name
	}

	var holdings []*models.InvestmentHoldingInfo

	for _, ua := range userAssets {
		// 3a. Get asset info
		asset := &models.Asset{}
		has, err := s.UserDataDB(0).NewSession(c).ID(ua.AssetId).Get(asset)
		if err != nil || !has {
			continue
		}

		// 3b. Get all non-deleted transactions for this asset, ordered by trade_time ASC
		var transactions []*models.InvestmentTransaction
		err = s.UserDataDB(uid).NewSession(c).
			Where("uid=? AND deleted=? AND asset_id=?", uid, false, ua.AssetId).
			OrderBy("trade_time asc").
			Find(&transactions)
		if err != nil || len(transactions) == 0 {
			continue
		}

		// 3c. Group transactions by accountId
		type accountState struct {
			totalQuantity int64
			totalCost     int64
		}
		accountStates := make(map[int64]*accountState)

		for _, tx := range transactions {
			state, ok := accountStates[tx.AccountId]
			if !ok {
				state = &accountState{}
				accountStates[tx.AccountId] = state
			}

			switch tx.Type {
			case models.INVESTMENT_TRANSACTION_TYPE_BUY:
				state.totalQuantity += tx.Quantity
				state.totalCost += tx.Amount + tx.Fee

			case models.INVESTMENT_TRANSACTION_TYPE_SELL:
				if state.totalQuantity > 0 {
					sellCost := state.totalCost * tx.Quantity / state.totalQuantity
					state.totalQuantity -= tx.Quantity
					state.totalCost -= sellCost
				}

			case models.INVESTMENT_TRANSACTION_TYPE_DIVIDEND_REINVEST:
				state.totalQuantity += tx.Quantity
				state.totalCost += tx.Amount

			case models.INVESTMENT_TRANSACTION_TYPE_SPLIT:
				state.totalQuantity += tx.Quantity
				// totalCost unchanged for splits

			case models.INVESTMENT_TRANSACTION_TYPE_CONVERSION_OUT:
				if state.totalQuantity > 0 {
					sellCost := state.totalCost * tx.Quantity / state.totalQuantity
					state.totalQuantity -= tx.Quantity
					state.totalCost -= sellCost
				}

			case models.INVESTMENT_TRANSACTION_TYPE_CONVERSION_IN:
				state.totalQuantity += tx.Quantity
				state.totalCost += tx.Amount + tx.Fee

			case models.INVESTMENT_TRANSACTION_TYPE_DIVIDEND_CASH:
				// No holding change
			}
		}

		// 3d. Get latest market price
		marketData := &models.MarketData{}
		_, _ = s.UserDataDB(uid).NewSession(c).
			Where("asset_id=?", ua.AssetId).
			OrderBy("date desc").
			Limit(1).
			Get(marketData)

		currentPrice := marketData.Price

		// 3e. Build one holding per account
		for accountId, state := range accountStates {
			if state.totalQuantity <= 0 {
				continue
			}

			marketValue := state.totalQuantity * currentPrice / 10000
			unrealizedPnl := marketValue - state.totalCost
			var returnRate int64
			if state.totalCost > 0 {
				returnRate = unrealizedPnl * 10000 / state.totalCost
			}

			var avgCostPrice int64
			if state.totalQuantity > 0 {
				avgCostPrice = state.totalCost * 10000 / state.totalQuantity
			}

			holdings = append(holdings, &models.InvestmentHoldingInfo{
				AssetId:       ua.AssetId,
				AssetCode:     asset.Code,
				AssetName:     asset.Name,
				Category:      asset.Category,
				Currency:      asset.Currency,
				Market:        asset.Market,
				AccountId:     accountId,
				AccountName:   accountNameMap[accountId],
				Quantity:      state.totalQuantity,
				AvgCostPrice:  avgCostPrice,
				TotalCost:     state.totalCost,
				CurrentPrice:      currentPrice,
				CurrentPriceDate:  marketData.Date,
				MarketValue:   marketValue,
				UnrealizedPnl: unrealizedPnl,
				ReturnRate:    returnRate,
			})
		}
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
