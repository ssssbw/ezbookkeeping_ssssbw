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

type InvestmentAssetService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

var InvestmentAssets = &InvestmentAssetService{
	ServiceUsingDB: ServiceUsingDB{
		container: datastore.Container,
	},
	ServiceUsingUuid: ServiceUsingUuid{
		container: uuid.Container,
	},
}

func (s *InvestmentAssetService) GetAllAssetsByUid(c core.Context, uid int64, assetType models.InvestmentAssetType, market models.InvestmentMarket) ([]*models.InvestmentAsset, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	condition := "uid=? AND deleted=?"
	conditionParams := make([]any, 0, 4)
	conditionParams = append(conditionParams, uid)
	conditionParams = append(conditionParams, false)

	if assetType > 0 {
		condition = condition + " AND type=?"
		conditionParams = append(conditionParams, assetType)
	}

	if market > 0 {
		condition = condition + " AND market=?"
		conditionParams = append(conditionParams, market)
	}

	var assets []*models.InvestmentAsset
	err := s.UserDataDB(uid).NewSession(c).Where(condition, conditionParams...).OrderBy("type asc, market asc, code asc").Find(&assets)

	return assets, err
}

func (s *InvestmentAssetService) GetAssetByAssetId(c core.Context, uid int64, assetId int64) (*models.InvestmentAsset, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if assetId <= 0 {
		return nil, errs.ErrInvestmentAssetIdInvalid
	}

	asset := &models.InvestmentAsset{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(assetId).Where("uid=? AND deleted=?", uid, false).Get(asset)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrInvestmentAssetNotFound
	}

	return asset, nil
}

func (s *InvestmentAssetService) GetAssetByAssetCode(c core.Context, uid int64, assetCode string) (*models.InvestmentAsset, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if assetCode == "" {
		return nil, errs.ErrInvestmentAssetIdInvalid
	}

	asset := &models.InvestmentAsset{}
	has, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND code=?", uid, false, assetCode).Get(asset)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrInvestmentAssetNotFound
	}

	return asset, nil
}

func (s *InvestmentAssetService) GetAssetsByAssetIds(c core.Context, uid int64, assetIds []int64) (map[int64]*models.InvestmentAsset, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if len(assetIds) <= 0 {
		return nil, errs.ErrInvestmentAssetIdInvalid
	}

	var assets []*models.InvestmentAsset
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).In("asset_id", assetIds).Find(&assets)

	if err != nil {
		return nil, err
	}

	assetMap := make(map[int64]*models.InvestmentAsset, len(assets))
	for _, asset := range assets {
		assetMap[asset.AssetId] = asset
	}

	return assetMap, nil
}

func (s *InvestmentAssetService) CreateAsset(c core.Context, asset *models.InvestmentAsset) error {
	if asset.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	asset.AssetId = s.GenerateUuid(uuid.UUID_TYPE_INVESTMENT_ASSET)

	if asset.AssetId < 1 {
		return errs.ErrSystemIsBusy
	}

	asset.Deleted = false
	asset.IsActive = true
	asset.CreatedUnixTime = time.Now().Unix()
	asset.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(asset.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(asset)
		return err
	})
}

func (s *InvestmentAssetService) ModifyAsset(c core.Context, asset *models.InvestmentAsset) error {
	if asset.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	asset.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(asset.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(asset.AssetId).Cols("type", "market", "code", "name", "currency", "is_active", "extra_info", "comment", "updated_unix_time").Where("uid=? AND deleted=?", asset.Uid, false).Update(asset)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrInvestmentAssetNotFound
		}

		return nil
	})
}

func (s *InvestmentAssetService) DeleteAsset(c core.Context, uid int64, assetId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.InvestmentAsset{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		exists, err := sess.Cols("uid", "deleted", "asset_id").Where("uid=? AND deleted=? AND asset_id=?", uid, false, assetId).Limit(1).Exist(&models.InvestmentTransaction{})

		if err != nil {
			return err
		} else if exists {
			return errs.ErrInvestmentAssetInUseCannotBeDeleted
		}

		deletedRows, err := sess.ID(assetId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrInvestmentAssetNotFound
		}

		return nil
	})
}
