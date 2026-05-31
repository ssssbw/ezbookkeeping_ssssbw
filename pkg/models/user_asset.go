package models

// UserAsset represents the relationship between user and asset
type UserAsset struct {
	Id              int64 `xorm:"PK comment('记录ID')"`
	Uid             int64 `xorm:"INDEX(IDX_user_asset_uid_asset_id) NOT NULL comment('用户ID')"`
	AssetId         int64 `xorm:"INDEX(IDX_user_asset_uid_asset_id) NOT NULL comment('资产ID')"`
	Deleted         bool  `xorm:"NOT NULL comment('是否删除')"`
	IsActive        bool  `xorm:"NOT NULL comment('是否活跃')"`
	Comment         string `xorm:"VARCHAR(255) NOT NULL comment('备注')"`
	CreatedUnixTime int64 `comment('创建时间')"`
	UpdatedUnixTime int64 `comment('更新时间')"`
	DeletedUnixTime int64 `comment('删除时间')"`
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
	Id       int64             `json:"id,string"`
	AssetId  int64             `json:"assetId,string"`
	IsActive bool              `json:"isActive"`
	Comment  string            `json:"comment"`
	Asset    *AssetInfoResponse `json:"asset,omitempty"`
}

func (ua *UserAsset) ToUserAssetInfoResponse() *UserAssetInfoResponse {
	return &UserAssetInfoResponse{
		Id:       ua.Id,
		AssetId:  ua.AssetId,
		IsActive: ua.IsActive,
		Comment:  ua.Comment,
	}
}
