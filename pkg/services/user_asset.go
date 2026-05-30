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

type UserAssetService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

var UserAssets = &UserAssetService{
	ServiceUsingDB: ServiceUsingDB{
		container: datastore.Container,
	},
	ServiceUsingUuid: ServiceUsingUuid{
		container: uuid.Container,
	},
}

func (s *UserAssetService) GetUserAssetsByUid(c core.Context, uid int64, isActive *bool) ([]*models.UserAsset, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	condition := "uid=?"
	conditionParams := make([]any, 0, 2)
	conditionParams = append(conditionParams, uid)

	if isActive != nil {
		condition = condition + " AND is_active=?"
		conditionParams = append(conditionParams, *isActive)
	}

	var userAssets []*models.UserAsset
	err := s.UserDataDB(uid).NewSession(c).Where(condition, conditionParams...).OrderBy("added_unix_time desc").Find(&userAssets)

	return userAssets, err
}

func (s *UserAssetService) GetUserAssetByAssetId(c core.Context, uid int64, assetId int64) (*models.UserAsset, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if assetId <= 0 {
		return nil, errs.ErrInvestmentAssetIdInvalid
	}

	userAsset := &models.UserAsset{}
	has, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND asset_id=?", uid, assetId).Get(userAsset)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrInvestmentAssetNotFound
	}

	return userAsset, nil
}

func (s *UserAssetService) AddUserAsset(c core.Context, uid int64, assetId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if assetId <= 0 {
		return errs.ErrInvestmentAssetIdInvalid
	}

	existing := &models.UserAsset{}
	has, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND asset_id=?", uid, assetId).Get(existing)
	if err != nil {
		return err
	}
	if has {
		return errs.ErrInvestmentAssetNotFound
	}

	userAsset := &models.UserAsset{
		Uid:           uid,
		AssetId:       assetId,
		IsActive:      true,
		AddedUnixTime: time.Now().Unix(),
	}

	userAsset.Id = s.GenerateUuid(uuid.UUID_TYPE_USER_ASSET)
	if userAsset.Id < 1 {
		return errs.ErrSystemIsBusy
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(userAsset)
		return err
	})
}

func (s *UserAssetService) RemoveUserAsset(c core.Context, uid int64, assetId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if assetId <= 0 {
		return errs.ErrInvestmentAssetIdInvalid
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		deletedRows, err := sess.Where("uid=? AND asset_id=?", uid, assetId).Delete(&models.UserAsset{})
		if err != nil {
			return err
		}
		if deletedRows < 1 {
			return errs.ErrInvestmentAssetNotFound
		}
		return nil
	})
}

func (s *UserAssetService) SetUserAssetActive(c core.Context, uid int64, assetId int64, isActive bool) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if assetId <= 0 {
		return errs.ErrInvestmentAssetIdInvalid
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.Where("uid=? AND asset_id=?", uid, assetId).Cols("is_active").Update(&models.UserAsset{IsActive: isActive})
		if err != nil {
			return err
		}
		if updatedRows < 1 {
			return errs.ErrInvestmentAssetNotFound
		}
		return nil
	})
}
