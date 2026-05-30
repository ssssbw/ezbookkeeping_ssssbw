package models

// AssetCategory represents the category of an asset
type AssetCategory string

const (
	AssetCategoryEquity      AssetCategory = "equity"       // 权益类
	AssetCategoryFixedIncome AssetCategory = "fixed_income" // 固定收益类
	AssetCategoryCommodity   AssetCategory = "commodity"    // 商品类
	AssetCategoryDigital     AssetCategory = "digital"      // 数字资产类
)

// Asset represents a global asset stored in database
type Asset struct {
	AssetId         int64         `xorm:"PK comment('资产ID')"`
	Code            string        `xorm:"INDEX(IDX_asset_code_market) NOT NULL comment('资产代码, 如 005827')"`
	Market          InvestmentMarket `xorm:"INDEX(IDX_asset_code_market) NOT NULL comment('市场: 1=中国, 2=香港, 3=美国')"`
	Name            string        `xorm:"NOT NULL comment('资产名称')"`
	Category        AssetCategory `xorm:"INDEX NOT NULL comment('类别: equity/fixed_income/commodity/digital')"`
	Currency        string        `xorm:"NOT NULL comment('计价货币: CNY/USD/HKD')"`
	Industry        string        `xorm:"INDEX comment('行业分类: technology/healthcare/consumer/financial/...')"`
	Tags            string        `xorm:"TEXT comment('标签JSON数组, 用于搜索')"`
	ExtraInfo       string        `xorm:"TEXT comment('扩展信息JSON: 基金公司/经理/费率等')"`
	CreatedUnixTime int64         `comment('创建时间')"`
	UpdatedUnixTime int64         `comment('更新时间')"`
}

// AssetListRequest represents all parameters of asset listing request
type AssetListRequest struct {
	Category AssetCategory    `form:"category"`
	Market   InvestmentMarket `form:"market"`
	Industry string           `form:"industry"`
}

// AssetGetRequest represents all parameters of asset getting request
type AssetGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// AssetSearchRequest represents all parameters of asset search request
type AssetSearchRequest struct {
	Keyword string `form:"keyword" binding:"required"`
	Limit   int    `form:"limit"`
}

// AssetCreateRequest represents all parameters of asset creation request
type AssetCreateRequest struct {
	Code      string           `json:"code" binding:"required,notBlank,max=20"`
	Market    InvestmentMarket `json:"market" binding:"required,min=1"`
	Name      string           `json:"name" binding:"required,notBlank,max=64"`
	Category  AssetCategory    `json:"category" binding:"required"`
	Currency  string           `json:"currency" binding:"required,len=3,validCurrency"`
	Industry  string           `json:"industry"`
	Tags      string           `json:"tags"`
	ExtraInfo string           `json:"extraInfo"`
}

// AssetModifyRequest represents all parameters of asset modification request
type AssetModifyRequest struct {
	Id        int64            `json:"id,string" binding:"required,min=1"`
	Code      string           `json:"code" binding:"required,notBlank,max=20"`
	Market    InvestmentMarket `json:"market" binding:"required,min=1"`
	Name      string           `json:"name" binding:"required,notBlank,max=64"`
	Category  AssetCategory    `json:"category" binding:"required"`
	Currency  string           `json:"currency" binding:"required,len=3,validCurrency"`
	Industry  string           `json:"industry"`
	Tags      string           `json:"tags"`
	ExtraInfo string           `json:"extraInfo"`
}

// AssetInfoResponse represents a view-object of asset
type AssetInfoResponse struct {
	Id        int64            `json:"id,string"`
	Code      string           `json:"code"`
	Market    InvestmentMarket `json:"market"`
	Name      string           `json:"name"`
	Category  AssetCategory    `json:"category"`
	Currency  string           `json:"currency"`
	Industry  string           `json:"industry,omitempty"`
	Tags      string           `json:"tags,omitempty"`
	ExtraInfo string           `json:"extraInfo,omitempty"`
}

// ToAssetInfoResponse returns a view-object according to database model
func (a *Asset) ToAssetInfoResponse() *AssetInfoResponse {
	return &AssetInfoResponse{
		Id:        a.AssetId,
		Code:      a.Code,
		Market:    a.Market,
		Name:      a.Name,
		Category:  a.Category,
		Currency:  a.Currency,
		Industry:  a.Industry,
		Tags:      a.Tags,
		ExtraInfo: a.ExtraInfo,
	}
}
