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
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

type InvestmentApi struct {
	ApiUsingConfig
	ApiUsingDuplicateChecker
	assets       *services.InvestmentAssetService
	transactions *services.InvestmentTransactionService
	marketData   *services.MarketDataService
	globalAssets *services.AssetService
	userAssets   *services.UserAssetService
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
	assets:       services.InvestmentAssets,
	transactions: services.InvestmentTransactions,
	marketData:   services.MarketData,
	globalAssets: services.Assets,
	userAssets:   services.UserAssets,
}

// Asset handlers

func (a *InvestmentApi) AssetListHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.InvestmentAssetListRequest
	err := c.ShouldBindQuery(&req)

	if err != nil {
		log.Warnf(c, "[investment.AssetListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	assets, err := a.assets.GetAllAssetsByUid(c, uid, req.Type, req.Market)

	if err != nil {
		log.Errorf(c, "[investment.AssetListHandler] failed to get assets for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	assetResps := make([]*models.InvestmentAssetInfoResponse, len(assets))
	for i, asset := range assets {
		assetResps[i] = asset.ToInvestmentAssetInfoResponse()
	}

	return assetResps, nil
}

func (a *InvestmentApi) AssetGetHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.InvestmentAssetGetRequest
	err := c.ShouldBindQuery(&req)

	if err != nil {
		log.Warnf(c, "[investment.AssetGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	asset, err := a.assets.GetAssetByAssetId(c, uid, req.Id)

	if err != nil {
		log.Errorf(c, "[investment.AssetGetHandler] failed to get asset \"id:%d\" for user \"uid:%d\", because %s", req.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return asset.ToInvestmentAssetInfoResponse(), nil
}

func (a *InvestmentApi) AssetCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.InvestmentAssetCreateRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Warnf(c, "[investment.AssetCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	asset := &models.InvestmentAsset{
		Uid:       uid,
		Type:      req.Type,
		Market:    req.Market,
		Code:      req.Code,
		Name:      req.Name,
		Currency:  req.Currency,
		ExtraInfo: req.ExtraInfo,
		Comment:   req.Comment,
	}

	if a.CurrentConfig().EnableDuplicateSubmissionsCheck && req.ClientSessionId != "" {
		found, remark := a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_TRANSACTION, uid, req.ClientSessionId)

		if found {
			log.Infof(c, "[investment.AssetCreateHandler] another asset \"id:%s\" has been created for user \"uid:%d\"", remark, uid)
			assetId, err := utils.StringToInt64(remark)

			if err == nil {
				asset, err = a.assets.GetAssetByAssetId(c, uid, assetId)

				if err != nil {
					log.Errorf(c, "[investment.AssetCreateHandler] failed to get existed asset \"id:%d\" for user \"uid:%d\", because %s", assetId, uid, err.Error())
					return nil, errs.Or(err, errs.ErrOperationFailed)
				}

				return asset.ToInvestmentAssetInfoResponse(), nil
			}
		}
	}

	err = a.assets.CreateAsset(c, asset)

	if err != nil {
		log.Errorf(c, "[investment.AssetCreateHandler] failed to create asset for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[investment.AssetCreateHandler] user \"uid:%d\" has created a new asset \"id:%d\" successfully", uid, asset.AssetId)

	a.SetSubmissionRemarkIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_TRANSACTION, uid, req.ClientSessionId, utils.Int64ToString(asset.AssetId))

	return asset.ToInvestmentAssetInfoResponse(), nil
}

func (a *InvestmentApi) AssetModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.InvestmentAssetModifyRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Warnf(c, "[investment.AssetModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	asset, err := a.assets.GetAssetByAssetId(c, uid, req.Id)

	if err != nil {
		log.Errorf(c, "[investment.AssetModifyHandler] failed to get asset \"id:%d\" for user \"uid:%d\", because %s", req.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	asset.Type = req.Type
	asset.Market = req.Market
	asset.Code = req.Code
	asset.Name = req.Name
	asset.Currency = req.Currency
	asset.IsActive = req.IsActive
	asset.ExtraInfo = req.ExtraInfo
	asset.Comment = req.Comment

	err = a.assets.ModifyAsset(c, asset)

	if err != nil {
		log.Errorf(c, "[investment.AssetModifyHandler] failed to update asset \"id:%d\" for user \"uid:%d\", because %s", asset.AssetId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[investment.AssetModifyHandler] user \"uid:%d\" has updated asset \"id:%d\" successfully", uid, asset.AssetId)

	return asset.ToInvestmentAssetInfoResponse(), nil
}

func (a *InvestmentApi) AssetDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.InvestmentAssetDeleteRequest
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Warnf(c, "[investment.AssetDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.assets.DeleteAsset(c, uid, req.Id)

	if err != nil {
		log.Errorf(c, "[investment.AssetDeleteHandler] failed to delete asset \"id:%d\" for user \"uid:%d\", because %s", req.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[investment.AssetDeleteHandler] user \"uid:%d\" has deleted asset \"id:%d\"", uid, req.Id)
	return true, nil
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

	txResps := make([]*models.InvestmentTransactionInfoResponse, len(transactions))
	for i, tx := range transactions {
		txResps[i] = tx.ToInvestmentTransactionInfoResponse()
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

	asset, err := a.assets.GetAssetByAssetCode(c, uid, req.AssetCode)
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

	asset, err := a.assets.GetAssetByAssetCode(c, uid, req.AssetCode)
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

// User Asset handlers

func (a *InvestmentApi) UserAssetListHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.UserAssetListRequest
	err := c.ShouldBindQuery(&req)

	if err != nil {
		log.Warnf(c, "[investment.UserAssetListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	userAssets, err := a.userAssets.GetUserAssetsByUid(c, uid, req.IsActive)
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

