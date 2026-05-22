package models

// InvestmentAsset represents investment asset data stored in database
type InvestmentAsset struct {
	AssetId         int64               `xorm:"PK comment('资产ID')"`
	Uid             int64               `xorm:"INDEX(IDX_investment_asset_uid_deleted_type_market) NOT NULL comment('用户ID')"`
	Deleted         bool                `xorm:"INDEX(IDX_investment_asset_uid_deleted_type_market) NOT NULL comment('是否删除')"`
	Type            InvestmentAssetType `xorm:"INDEX(IDX_investment_asset_uid_deleted_type_market) NOT NULL comment('资产类型: 1=基金, 2=股票, 3=ETF, 4=债券, 5=加密货币')"`
	Market          InvestmentMarket    `xorm:"INDEX(IDX_investment_asset_uid_deleted_type_market) NOT NULL comment('市场: 1=中国, 2=香港, 3=美国')"`
	Code            string              `xorm:"VARCHAR(20) NOT NULL comment('资产代码, 如 005827')"`
	Name            string              `xorm:"VARCHAR(64) NOT NULL comment('资产名称')"`
	Currency        string              `xorm:"VARCHAR(3) NOT NULL comment('计价货币, 如 CNY/USD/HKD')"`
	IsActive        bool                `xorm:"NOT NULL comment('是否活跃')"`
	ExtraInfo       string              `xorm:"BLOB comment('扩展信息JSON: 行业/基金公司/经理/费率/持仓分布')"`
	Comment         string              `xorm:"VARCHAR(255) NOT NULL comment('备注')"`
	CreatedUnixTime int64               `comment('创建时间')"`
	UpdatedUnixTime int64               `comment('更新时间')"`
	DeletedUnixTime int64               `comment('删除时间')"`
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
