package models

// MarketData represents market data stored in database
type MarketData struct {
	DataId          int64 `xorm:"PK comment('行情ID')"`
	AssetId         int64 `xorm:"UNIQUE(UQE_market_data_asset_id_date) NOT NULL comment('关联资产ID')"`
	Date            int64 `xorm:"UNIQUE(UQE_market_data_asset_id_date) NOT NULL comment('日期, Unix时间戳取0点')"`
	Price           int64 `xorm:"NOT NULL comment('当日净值/收盘价, 精度 x10000')"`
	Volume          int64 `comment('成交量, ETF/股票有效')"`
	CreatedUnixTime int64 `comment('创建时间')"`
	UpdatedUnixTime int64 `comment('更新时间')"`
}

// MarketDataListRequest represents all parameters of market data listing request
type MarketDataListRequest struct {
	AssetId   int64 `form:"asset_id,string" binding:"required,min=1"`
	StartTime int64 `form:"start_time"`
	EndTime   int64 `form:"end_time"`
}

// MarketDataGetRequest represents all parameters of market data getting request
type MarketDataGetRequest struct {
	AssetId int64 `form:"asset_id,string" binding:"required,min=1"`
	Date    int64 `form:"date" binding:"required,min=1"`
}

// MarketDataCreateRequest represents all parameters of market data creation request
type MarketDataCreateRequest struct {
	AssetId int64 `json:"assetId,string" binding:"required,min=1"`
	Date    int64 `json:"date" binding:"required,min=1"`
	Price   int64 `json:"price" binding:"required,min=1"`
	Volume  int64 `json:"volume"`
}

// MarketDataModifyRequest represents all parameters of market data modification request
type MarketDataModifyRequest struct {
	AssetId int64 `json:"assetId,string" binding:"required,min=1"`
	Date    int64 `json:"date" binding:"required,min=1"`
	Price   int64 `json:"price" binding:"required,min=1"`
	Volume  int64 `json:"volume"`
}

// MarketDataDeleteRequest represents all parameters of market data deleting request
type MarketDataDeleteRequest struct {
	AssetId int64 `json:"assetId,string" binding:"required,min=1"`
	Date    int64 `json:"date" binding:"required,min=1"`
}

// MarketDataInfoResponse represents a view-object of market data
type MarketDataInfoResponse struct {
	AssetId int64 `json:"assetId,string"`
	Date    int64 `json:"date"`
	Price   int64 `json:"price"`
	Volume  int64 `json:"volume,omitempty"`
}

// ToMarketDataInfoResponse returns a view-object according to database model
func (m *MarketData) ToMarketDataInfoResponse() *MarketDataInfoResponse {
	return &MarketDataInfoResponse{
		AssetId: m.AssetId,
		Date:    m.Date,
		Price:   m.Price,
		Volume:  m.Volume,
	}
}
