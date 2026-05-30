package models

// UserAsset represents the relationship between user and asset
type UserAsset struct {
	Id              int64 `xorm:"PK comment('记录ID')"`
	Uid             int64 `xorm:"INDEX(IDX_user_asset_uid_asset_id) NOT NULL comment('用户ID')"`
	AssetId         int64 `xorm:"INDEX(IDX_user_asset_uid_asset_id) NOT NULL comment('资产ID')"`
	IsActive        bool  `xorm:"NOT NULL comment('是否活跃')"`
	AddedUnixTime   int64 `comment('添加时间')"`
}

// UserAssetListRequest represents all parameters of user asset listing request
type UserAssetListRequest struct {
	IsActive *bool `form:"is_active"`
}

// UserAssetAddRequest represents all parameters of adding user asset
type UserAssetAddRequest struct {
	AssetId int64 `json:"assetId,string" binding:"required,min=1"`
}

// UserAssetRemoveRequest represents all parameters of removing user asset
type UserAssetRemoveRequest struct {
	AssetId int64 `json:"assetId,string" binding:"required,min=1"`
}

// UserAssetInfoResponse represents a view-object of user asset
type UserAssetInfoResponse struct {
	Id       int64          `json:"id,string"`
	AssetId  int64          `json:"assetId,string"`
	IsActive bool           `json:"isActive"`
	Asset    *AssetInfoResponse `json:"asset,omitempty"`
}

// ToUserAssetInfoResponse returns a view-object according to database model
func (ua *UserAsset) ToUserAssetInfoResponse() *UserAssetInfoResponse {
	return &UserAssetInfoResponse{
		Id:       ua.Id,
		AssetId:  ua.AssetId,
		IsActive: ua.IsActive,
	}
}
