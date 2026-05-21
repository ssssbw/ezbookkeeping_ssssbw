package models

// InvestmentAsset represents investment asset data stored in database
type InvestmentAsset struct {
	AssetId         int64               `xorm:"PK"`
	Uid             int64               `xorm:"INDEX(IDX_investment_asset_uid_deleted_type_market) NOT NULL"`
	Deleted         bool                `xorm:"INDEX(IDX_investment_asset_uid_deleted_type_market) NOT NULL"`
	Type            InvestmentAssetType `xorm:"INDEX(IDX_investment_asset_uid_deleted_type_market) NOT NULL"`
	Market          InvestmentMarket    `xorm:"INDEX(IDX_investment_asset_uid_deleted_type_market) NOT NULL"`
	Code            string              `xorm:"VARCHAR(20) NOT NULL"`
	Name            string              `xorm:"VARCHAR(64) NOT NULL"`
	Currency        string              `xorm:"VARCHAR(3) NOT NULL"`
	IsActive        bool                `xorm:"NOT NULL"`
	ExtraInfo       string              `xorm:"BLOB"`
	Comment         string              `xorm:"VARCHAR(255) NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}

// InvestmentAssetListRequest represents all parameters of investment asset listing request
type InvestmentAssetListRequest struct {
	Type     InvestmentAssetType `form:"type"`
	Market   InvestmentMarket    `form:"market"`
	IsActive *bool               `form:"is_active"`
}

// InvestmentAssetGetRequest represents all parameters of investment asset getting request
type InvestmentAssetGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// InvestmentAssetCreateRequest represents all parameters of investment asset creation request
type InvestmentAssetCreateRequest struct {
	Type            InvestmentAssetType `json:"type" binding:"required,min=1"`
	Market          InvestmentMarket    `json:"market" binding:"required,min=1"`
	Code            string              `json:"code" binding:"required,notBlank,max=20"`
	Name            string              `json:"name" binding:"required,notBlank,max=64"`
	Currency        string              `json:"currency" binding:"required,len=3,validCurrency"`
	ExtraInfo       string              `json:"extraInfo"`
	Comment         string              `json:"comment" binding:"max=255"`
	ClientSessionId string              `json:"clientSessionId"`
}

// InvestmentAssetModifyRequest represents all parameters of investment asset modification request
type InvestmentAssetModifyRequest struct {
	Id        int64               `json:"id,string" binding:"required,min=1"`
	Type      InvestmentAssetType `json:"type" binding:"required,min=1"`
	Market    InvestmentMarket    `json:"market" binding:"required,min=1"`
	Code      string              `json:"code" binding:"required,notBlank,max=20"`
	Name      string              `json:"name" binding:"required,notBlank,max=64"`
	Currency  string              `json:"currency" binding:"required,len=3,validCurrency"`
	IsActive  bool                `json:"isActive"`
	ExtraInfo string              `json:"extraInfo"`
	Comment   string              `json:"comment" binding:"max=255"`
}

// InvestmentAssetDeleteRequest represents all parameters of investment asset deleting request
type InvestmentAssetDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// InvestmentAssetInfoResponse represents a view-object of investment asset
type InvestmentAssetInfoResponse struct {
	Id        int64               `json:"id,string"`
	Type      InvestmentAssetType `json:"type"`
	Market    InvestmentMarket    `json:"market"`
	Code      string              `json:"code"`
	Name      string              `json:"name"`
	Currency  string              `json:"currency"`
	IsActive  bool                `json:"isActive"`
	ExtraInfo string              `json:"extraInfo,omitempty"`
	Comment   string              `json:"comment"`
}

// ToInvestmentAssetInfoResponse returns a view-object according to database model
func (a *InvestmentAsset) ToInvestmentAssetInfoResponse() *InvestmentAssetInfoResponse {
	return &InvestmentAssetInfoResponse{
		Id:        a.AssetId,
		Type:      a.Type,
		Market:    a.Market,
		Code:      a.Code,
		Name:      a.Name,
		Currency:  a.Currency,
		IsActive:  a.IsActive,
		ExtraInfo: a.ExtraInfo,
		Comment:   a.Comment,
	}
}
