package services

import (
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

type AssetService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

var Assets = &AssetService{
	ServiceUsingDB: ServiceUsingDB{
		container: datastore.Container,
	},
	ServiceUsingUuid: ServiceUsingUuid{
		container: uuid.Container,
	},
}

func (s *AssetService) GetAllAssets(c core.Context, category models.AssetCategory, market models.InvestmentMarket, industry string) ([]*models.Asset, error) {
	condition := "1=1"
	conditionParams := make([]any, 0, 3)

	if category != "" {
		condition = condition + " AND category=?"
		conditionParams = append(conditionParams, category)
	}

	if market > 0 {
		condition = condition + " AND market=?"
		conditionParams = append(conditionParams, market)
	}

	if industry != "" {
		condition = condition + " AND industry=?"
		conditionParams = append(conditionParams, industry)
	}

	var assets []*models.Asset
	err := s.UserDataDB(0).NewSession(c).Where(condition, conditionParams...).OrderBy("code asc").Find(&assets)

	return assets, err
}

func (s *AssetService) GetAssetByAssetId(c core.Context, assetId int64) (*models.Asset, error) {
	if assetId <= 0 {
		return nil, errs.ErrInvestmentAssetIdInvalid
	}

	asset := &models.Asset{}
	has, err := s.UserDataDB(0).NewSession(c).ID(assetId).Get(asset)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrInvestmentAssetNotFound
	}

	return asset, nil
}

func (s *AssetService) GetAssetByCodeAndMarket(c core.Context, code string, market models.InvestmentMarket) (*models.Asset, error) {
	if code == "" {
		return nil, errs.ErrInvestmentAssetIdInvalid
	}

	asset := &models.Asset{}
	has, err := s.UserDataDB(0).NewSession(c).Where("code=? AND market=?", code, market).Get(asset)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrInvestmentAssetNotFound
	}

	return asset, nil
}

func (s *AssetService) SearchAssets(c core.Context, keyword string, limit int) ([]*models.Asset, error) {
	if keyword == "" {
		return nil, errs.ErrInvestmentAssetIdInvalid
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	condition := "code LIKE ? OR name LIKE ?"
	keywordPattern := "%" + keyword + "%"

	var assets []*models.Asset
	err := s.UserDataDB(0).NewSession(c).Where(condition, keywordPattern, keywordPattern).Limit(limit).OrderBy("code asc").Find(&assets)

	return assets, err
}

func (s *AssetService) CreateAsset(c core.Context, asset *models.Asset) error {
	if asset.Code == "" || asset.Name == "" {
		return errs.ErrInvestmentAssetIdInvalid
	}

	asset.AssetId = s.GenerateUuid(uuid.UUID_TYPE_ASSET)

	if asset.AssetId < 1 {
		return errs.ErrSystemIsBusy
	}

	now := time.Now().Unix()
	asset.CreatedUnixTime = now
	asset.UpdatedUnixTime = now

	return s.UserDataDB(0).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(asset)
		return err
	})
}

func (s *AssetService) ModifyAsset(c core.Context, asset *models.Asset) error {
	if asset.AssetId <= 0 {
		return errs.ErrInvestmentAssetIdInvalid
	}

	asset.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(0).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(asset.AssetId).Cols("code", "market", "name", "category", "currency", "industry", "tags", "extra_info", "updated_unix_time").Update(asset)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrInvestmentAssetNotFound
		}

		return nil
	})
}
