package services

import (
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

type MarketDataService struct {
	ServiceUsingDB
}

var MarketData = &MarketDataService{
	ServiceUsingDB: ServiceUsingDB{
		container: datastore.Container,
	},
}

func (s *MarketDataService) GetLatestPrice(c core.Context, uid int64, assetId int64) (*models.MarketData, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if assetId <= 0 {
		return nil, errs.ErrMarketDataAssetIdInvalid
	}

	data := &models.MarketData{}
	has, err := s.UserDataDB(uid).NewSession(c).Where("asset_id=?", assetId).OrderBy("date desc").Limit(1).Get(data)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrMarketDataNotFound
	}

	return data, nil
}

func (s *MarketDataService) GetMarketDataByAssetId(c core.Context, uid int64, assetId int64, startTime int64, endTime int64) ([]*models.MarketData, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if assetId <= 0 {
		return nil, errs.ErrMarketDataAssetIdInvalid
	}

	condition := "asset_id=?"
	conditionParams := make([]any, 0, 3)
	conditionParams = append(conditionParams, assetId)

	if startTime > 0 {
		condition = condition + " AND date>=?"
		conditionParams = append(conditionParams, startTime)
	}

	if endTime > 0 {
		condition = condition + " AND date<=?"
		conditionParams = append(conditionParams, endTime)
	}

	var data []*models.MarketData
	err := s.UserDataDB(uid).NewSession(c).Where(condition, conditionParams...).OrderBy("date asc").Find(&data)

	return data, err
}

func (s *MarketDataService) CreateMarketData(c core.Context, uid int64, data *models.MarketData) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if data.AssetId <= 0 {
		return errs.ErrMarketDataAssetIdInvalid
	}

	now := time.Now().Unix()
	data.CreatedUnixTime = now
	data.UpdatedUnixTime = now

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		existing := &models.MarketData{}
		has, err := sess.Where("asset_id=? AND date=?", data.AssetId, data.Date).Get(existing)

		if err != nil {
			return err
		}

		if has {
			existing.Price = data.Price
			existing.Volume = data.Volume
			existing.UpdatedUnixTime = now

			_, err = sess.ID(existing.DataId).Cols("price", "volume", "updated_unix_time").Update(existing)
			return err
		}

		_, err = sess.Insert(data)
		return err
	})
}

func (s *MarketDataService) ModifyMarketData(c core.Context, uid int64, data *models.MarketData) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	data.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(data.DataId).Cols("price", "volume", "updated_unix_time").Where("asset_id=?", data.AssetId).Update(data)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrMarketDataNotFound
		}

		return nil
	})
}
