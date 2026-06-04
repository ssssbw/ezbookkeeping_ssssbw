package api

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/duplicatechecker"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/marketdata"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

type InvestmentApi struct {
	ApiUsingConfig
	ApiUsingDuplicateChecker
	transactions *services.InvestmentTransactionService
	marketData   *services.MarketDataService
	globalAssets *services.AssetService
	userAssets   *services.UserAssetService
	analysis     *services.InvestmentAnalysisService
}

var Investment = &InvestmentApi{
	ApiUsingConfig: ApiUsingConfig{
		container: settings.Container,
	},
	ApiUsingDuplicateChecker: ApiUsingDuplicateChecker{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		container: duplicatechecker.Container,
	},
	transactions: services.InvestmentTransactions,
	marketData:   services.MarketData,
	globalAssets: services.Assets,
	userAssets:   services.UserAssets,
	analysis:     services.InvestmentAnalysis,
}

// Transaction handlers

func (a *InvestmentApi) TransactionListHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.InvestmentTransactionListRequest
	err := c.ShouldBindQuery(&req)

	if err != nil {
		log.Warnf(c, "[investment.TransactionListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	transactions, err := a.transactions.GetAllTransactionsByUid(c, uid, req.AssetId, req.AccountId, req.Type, req.StartTime, req.EndTime)

	if err != nil {
		log.Errorf(c, "[investment.TransactionListHandler] failed to get transactions for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	// Collect unique assetIds and accountIds for batch lookup
	assetIdSet := make(map[int64]bool)
	accountIdSet := make(map[int64]bool)
	for _, tx := range transactions {
		assetIdSet[tx.AssetId] = true
		accountIdSet[tx.AccountId] = true
	}

	// Batch load assets
	assetMap := make(map[int64]*models.Asset)
	for assetId := range assetIdSet {
		asset, err := a.globalAssets.GetAssetByAssetId(c, assetId)
		if err == nil && asset != nil {
			assetMap[assetId] = asset
		}
	}

	// Batch load accounts
	accountMap := make(map[int64]string)
	if len(accountIdSet) > 0 {
		ids := make([]int64, 0, len(accountIdSet))
		for id := range accountIdSet {
			ids = append(ids, id)
		}
		var accounts []*models.Account
		err := a.transactions.UserDataDB(uid).NewSession(c).In("account_id", ids).Where("uid=? AND deleted=?", uid, false).Find(&accounts)
		if err == nil {
			for _, acc := range accounts {
				accountMap[acc.AccountId] = acc.Name
			}
		}
	}

	// Build response with embedded info
	txResps := make([]*models.InvestmentTransactionInfoResponse, len(transactions))
	for i, tx := range transactions {
		var assetName, assetCode string
		if asset, ok := assetMap[tx.AssetId]; ok {
			assetName = asset.Name
			assetCode = asset.Code
		}
		accountName := accountMap[tx.AccountId]
		txResps[i] = tx.ToInvestmentTransactionInfoResponseWithInfo(assetName, assetCode, accountName)
	}

	return txResps, nil
}

func (a *InvestmentApi) TransactionGetHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.InvestmentTransactionGetRequest
	err := c.ShouldBindQuery(&req)

	if err != nil {
		log.Warnf(c, "[investment.TransactionGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	tx, err := a.transactions.GetTransactionByTransactionId(c, uid, req.Id)

	if err != nil {
		log.Errorf(c, "[investment.TransactionGetHandler] failed to get transaction \"id:%d\" for user \"uid:%d\", because %s", req.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return tx.ToInvestmentTransactionInfoResponse(), nil
}

func (a *InvestmentApi) TransactionCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.InvestmentTransactionCreateRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Warnf(c, "[investment.TransactionCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	tx := &models.InvestmentTransaction{
		Uid:                  uid,
		AssetId:              req.AssetId,
		AccountId:            req.AccountId,
		Type:                 req.Type,
		TradeTime:            req.TradeTime,
		ConfirmTime:          req.ConfirmTime,
		Quantity:             req.Quantity,
		Price:                req.Price,
		Amount:               req.Amount,
		Fee:                  req.Fee,
		RelatedTransactionId: req.RelatedTransactionId,
		TimezoneUtcOffset:    req.TimezoneUtcOffset,
		Comment:              req.Comment,
	}

	err = a.transactions.CreateTransaction(c, tx)

	if err != nil {
		log.Errorf(c, "[investment.TransactionCreateHandler] failed to create transaction for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[investment.TransactionCreateHandler] user \"uid:%d\" has created a new transaction \"id:%d\" successfully", uid, tx.TransactionId)

	return tx.ToInvestmentTransactionInfoResponse(), nil
}

func (a *InvestmentApi) TransactionModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.InvestmentTransactionModifyRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Warnf(c, "[investment.TransactionModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	tx, err := a.transactions.GetTransactionByTransactionId(c, uid, req.Id)

	if err != nil {
		log.Errorf(c, "[investment.TransactionModifyHandler] failed to get transaction \"id:%d\" for user \"uid:%d\", because %s", req.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	tx.AssetId = req.AssetId
	tx.AccountId = req.AccountId
	tx.Type = req.Type
	tx.TradeTime = req.TradeTime
	tx.ConfirmTime = req.ConfirmTime
	tx.Quantity = req.Quantity
	tx.Price = req.Price
	tx.Amount = req.Amount
	tx.Fee = req.Fee
	tx.RelatedTransactionId = req.RelatedTransactionId
	tx.TimezoneUtcOffset = req.TimezoneUtcOffset
	tx.Comment = req.Comment

	err = a.transactions.ModifyTransaction(c, tx)

	if err != nil {
		log.Errorf(c, "[investment.TransactionModifyHandler] failed to update transaction \"id:%d\" for user \"uid:%d\", because %s", tx.TransactionId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[investment.TransactionModifyHandler] user \"uid:%d\" has updated transaction \"id:%d\" successfully", uid, tx.TransactionId)

	return tx.ToInvestmentTransactionInfoResponse(), nil
}

func (a *InvestmentApi) TransactionDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.InvestmentTransactionDeleteRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Warnf(c, "[investment.TransactionDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.transactions.DeleteTransaction(c, uid, req.Id)

	if err != nil {
		log.Errorf(c, "[investment.TransactionDeleteHandler] failed to delete transaction \"id:%d\" for user \"uid:%d\", because %s", req.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[investment.TransactionDeleteHandler] user \"uid:%d\" has deleted transaction \"id:%d\"", uid, req.Id)
	return true, nil
}

// MarketData handlers

func (a *InvestmentApi) MarketDataLatestHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.MarketDataGetRequest
	err := c.ShouldBindQuery(&req)

	if err != nil {
		log.Warnf(c, "[investment.MarketDataLatestHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	data, err := a.marketData.GetLatestPrice(c, uid, req.AssetId)

	if err != nil {
		log.Errorf(c, "[investment.MarketDataLatestHandler] failed to get latest price for asset \"id:%d\", because %s", req.AssetId, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return data.ToMarketDataInfoResponse(), nil
}

func (a *InvestmentApi) MarketDataListHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.MarketDataListRequest
	err := c.ShouldBindQuery(&req)

	if err != nil {
		log.Warnf(c, "[investment.MarketDataListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	dataList, err := a.marketData.GetMarketDataByAssetId(c, uid, req.AssetId, req.StartTime, req.EndTime)

	if err != nil {
		log.Errorf(c, "[investment.MarketDataListHandler] failed to get market data for asset \"id:%d\", because %s", req.AssetId, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	dataResps := make([]*models.MarketDataInfoResponse, len(dataList))
	for i, data := range dataList {
		dataResps[i] = data.ToMarketDataInfoResponse()
	}

	return dataResps, nil
}

func (a *InvestmentApi) MarketDataCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.MarketDataCreateRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Warnf(c, "[investment.MarketDataCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	data := &models.MarketData{
		AssetId: req.AssetId,
		Date:    req.Date,
		Price:   req.Price,
		Volume:  req.Volume,
	}

	err = a.marketData.CreateMarketData(c, uid, data)

	if err != nil {
		log.Errorf(c, "[investment.MarketDataCreateHandler] failed to create market data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[investment.MarketDataCreateHandler] user \"uid:%d\" has created market data for asset \"id:%d\" successfully", uid, data.AssetId)

	return data.ToMarketDataInfoResponse(), nil
}

func (a *InvestmentApi) MarketDataModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.MarketDataModifyRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Warnf(c, "[investment.MarketDataModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	data := &models.MarketData{
		AssetId: req.AssetId,
		Date:    req.Date,
		Price:   req.Price,
		Volume:  req.Volume,
	}

	err = a.marketData.ModifyMarketData(c, uid, data)

	if err != nil {
		log.Errorf(c, "[investment.MarketDataModifyHandler] failed to update market data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[investment.MarketDataModifyHandler] user \"uid:%d\" has updated market data for asset \"id:%d\" successfully", uid, data.AssetId)

	return data.ToMarketDataInfoResponse(), nil
}

func (a *InvestmentApi) MarketDataRefreshHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()

	err := a.marketData.FetchAllActiveAssetsMarketData(c)

	if err != nil {
		log.Errorf(c, "[investment.MarketDataRefreshHandler] failed to refresh market data for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[investment.MarketDataRefreshHandler] user \"uid:%d\" has refreshed market data successfully", uid)

	return "ok", nil
}

func (a *InvestmentApi) MarketDataInitHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.MarketDataInitRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Warnf(c, "[investment.MarketDataInitHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	asset, err := a.globalAssets.GetAssetByCodeAndMarket(c, req.AssetCode, models.INVESTMENT_MARKET_CN)
	if err != nil {
		log.Errorf(c, "[investment.MarketDataInitHandler] failed to get asset for user \"uid:%d\", code \"%s\", because %s", uid, req.AssetCode, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	count, err := a.marketData.InitAssetMarketData(c, uid, asset.AssetId, asset.Code, string(asset.Market), req.TradeTime)
	if err != nil {
		log.Errorf(c, "[investment.MarketDataInitHandler] failed to init market data for user \"uid:%d\", asset \"%s\", because %s", uid, req.AssetCode, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[investment.MarketDataInitHandler] user \"uid:%d\" has initialized %d market data records for asset \"%s\" successfully", uid, count, req.AssetCode)

	return &models.MarketDataInitResponse{
		Count:     count,
		StartTime: req.TradeTime,
		EndTime:   time.Now().Unix(),
	}, nil
}

func (a *InvestmentApi) MarketDataEstimateHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.MarketDataEstimateRequest
	err := c.ShouldBindQuery(&req)

	if err != nil {
		log.Warnf(c, "[investment.MarketDataEstimateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	asset, err := a.globalAssets.GetAssetByCodeAndMarket(c, req.AssetCode, models.INVESTMENT_MARKET_CN)
	if err != nil {
		log.Errorf(c, "[investment.MarketDataEstimateHandler] failed to get asset for user \"uid:%d\", code \"%s\", because %s", uid, req.AssetCode, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	result, err := marketdata.Container.GetRealtimeEstimate(asset.Code, string(asset.Market))
	if err != nil {
		log.Errorf(c, "[investment.MarketDataEstimateHandler] failed to get estimate for user \"uid:%d\", asset \"%s\", because %s", uid, req.AssetCode, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if result == nil {
		return nil, nil
	}

	marketData, ok := result.Data.(*models.MarketData)
	if !ok {
		return nil, nil
	}

	log.Infof(c, "[investment.MarketDataEstimateHandler] user \"uid:%d\" got estimate for asset \"%s\": %d", uid, req.AssetCode, marketData.Price)

	return marketData.ToMarketDataInfoResponse(), nil
}

// Global Asset handlers

func (a *InvestmentApi) AssetSearchHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.AssetSearchRequest
	err := c.ShouldBindQuery(&req)

	if err != nil {
		log.Warnf(c, "[investment.AssetSearchHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	assets, err := a.globalAssets.SearchAssets(c, req.Keyword, req.Limit)
	if err != nil {
		log.Errorf(c, "[investment.AssetSearchHandler] failed to search assets, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	var response []*models.AssetInfoResponse
	for _, asset := range assets {
		response = append(response, asset.ToAssetInfoResponse())
	}

	return response, nil
}

func (a *InvestmentApi) GlobalAssetGetHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.AssetGetRequest
	err := c.ShouldBindQuery(&req)

	if err != nil {
		log.Warnf(c, "[investment.GlobalAssetGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	asset, err := a.globalAssets.GetAssetByAssetId(c, req.Id)
	if err != nil {
		log.Errorf(c, "[investment.GlobalAssetGetHandler] failed to get asset, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return asset.ToAssetInfoResponse(), nil
}

func (a *InvestmentApi) GlobalAssetCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.AssetCreateRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Warnf(c, "[investment.GlobalAssetCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	asset := &models.Asset{
		Code:     req.Code,
		Market:   req.Market,
		Name:     req.Name,
		Category: req.Category,
		Currency: req.Currency,
		Industry: req.Industry,
		Tags:     req.Tags,
		ExtraInfo: req.ExtraInfo,
	}

	err = a.globalAssets.CreateAsset(c, asset)
	if err != nil {
		log.Errorf(c, "[investment.AssetCreateHandler] failed to create asset, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[investment.AssetCreateHandler] asset \"%s\" has been created successfully", asset.Code)

	return asset.ToAssetInfoResponse(), nil
}

func (a *InvestmentApi) GlobalAssetListHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.AssetListRequest
	err := c.ShouldBindQuery(&req)

	if err != nil {
		log.Warnf(c, "[investment.GlobalAssetListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 50
	}

	totalCount, err := a.globalAssets.GetAllAssetsCount(c, req.Category, req.Market, req.Industry, req.Keyword)
	if err != nil {
		log.Errorf(c, "[investment.GlobalAssetListHandler] failed to count assets, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	assets, err := a.globalAssets.GetAllAssets(c, req.Category, req.Market, req.Industry, req.Keyword)
	if err != nil {
		log.Errorf(c, "[investment.GlobalAssetListHandler] failed to list assets, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	assetResps := make([]*models.AssetInfoResponse, len(assets))
	for i, asset := range assets {
		assetResps[i] = asset.ToAssetInfoResponse()
	}

	return &models.AssetListResponse{
		TotalCount: totalCount,
		Assets:     assetResps,
	}, nil
}

func (a *InvestmentApi) GlobalAssetModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.AssetModifyRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Warnf(c, "[investment.GlobalAssetModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	asset := &models.Asset{
		AssetId:  req.Id,
		Code:     req.Code,
		Market:   req.Market,
		Name:     req.Name,
		Category: req.Category,
		Currency: req.Currency,
		Industry: req.Industry,
		Tags:     req.Tags,
		ExtraInfo: req.ExtraInfo,
	}

	err = a.globalAssets.ModifyAsset(c, asset)
	if err != nil {
		log.Errorf(c, "[investment.GlobalAssetModifyHandler] failed to modify asset \"id:%d\", because %s", req.Id, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[investment.GlobalAssetModifyHandler] asset \"id:%d\" has been modified successfully", req.Id)

	return asset.ToAssetInfoResponse(), nil
}

func (a *InvestmentApi) GlobalAssetDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.AssetDeleteRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Warnf(c, "[investment.GlobalAssetDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	err = a.globalAssets.DeleteAsset(c, req.Id)
	if err != nil {
		log.Errorf(c, "[investment.GlobalAssetDeleteHandler] failed to delete asset \"id:%d\", because %s", req.Id, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[investment.GlobalAssetDeleteHandler] asset \"id:%d\" has been deleted successfully", req.Id)

	return true, nil
}

// User Asset handlers

func (a *InvestmentApi) UserAssetListHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.UserAssetListRequest
	err := c.ShouldBindQuery(&req)

	if err != nil {
		log.Warnf(c, "[investment.UserAssetListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	userAssets, err := a.userAssets.GetUserAssetsByUid(c, uid, req.IsActive, req.IsWatchlist)
	if err != nil {
		log.Errorf(c, "[investment.UserAssetListHandler] failed to get user assets, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	var response []*models.UserAssetInfoResponse
	for _, ua := range userAssets {
		resp := ua.ToUserAssetInfoResponse()

		asset, err := a.globalAssets.GetAssetByAssetId(c, ua.AssetId)
		if err == nil {
			resp.Asset = asset.ToAssetInfoResponse()
		}

		response = append(response, resp)
	}

	return response, nil
}

func (a *InvestmentApi) UserAssetAddHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.UserAssetAddRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Warnf(c, "[investment.UserAssetAddHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	err = a.userAssets.AddUserAsset(c, uid, req.AssetId)
	if err != nil {
		log.Errorf(c, "[investment.UserAssetAddHandler] failed to add user asset, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[investment.UserAssetAddHandler] user \"uid:%d\" has added asset \"id:%d\" successfully", uid, req.AssetId)

	return "ok", nil
}

func (a *InvestmentApi) UserAssetRemoveHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.UserAssetRemoveRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Warnf(c, "[investment.UserAssetRemoveHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	err = a.userAssets.RemoveUserAsset(c, uid, req.AssetId)
	if err != nil {
		log.Errorf(c, "[investment.UserAssetRemoveHandler] failed to remove user asset, because %s", err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[investment.UserAssetRemoveHandler] user \"uid:%d\" has removed asset \"id:%d\" successfully", uid, req.AssetId)

	return "ok", nil
}

// HoldingsHandler returns computed holding info for all active user assets
func (a *InvestmentApi) HoldingsHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	holdings, err := a.analysis.GetHoldings(c, uid)

	if err != nil {
		log.Errorf(c, "[investment.HoldingsHandler] failed to get holdings for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return holdings, nil
}

// OverviewHandler returns aggregated portfolio summary
func (a *InvestmentApi) OverviewHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	overview, err := a.analysis.GetOverview(c, uid)

	if err != nil {
		log.Errorf(c, "[investment.OverviewHandler] failed to get overview for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return overview, nil
}

// AdminCheckHandler checks if current user is investment admin
func (a *InvestmentApi) AdminCheckHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	config := a.CurrentConfig()
	isAdmin := config.InvestmentAdminUid > 0 && config.InvestmentAdminUid == uid
	log.Infof(c, "[investment.AdminCheckHandler] user \"uid:%d\" is investment admin: %v", uid, isAdmin)

	return map[string]interface{}{
		"isAdmin": isAdmin,
	}, nil
}

