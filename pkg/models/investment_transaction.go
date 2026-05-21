package models

// InvestmentTransaction represents investment transaction data stored in database
type InvestmentTransaction struct {
	TransactionId        int64                     `xorm:"PK"`
	Uid                  int64                     `xorm:"INDEX(IDX_invest_trans_uid_deleted_asset_id) INDEX(IDX_invest_trans_uid_deleted_account_id) NOT NULL"`
	Deleted              bool                      `xorm:"INDEX(IDX_invest_trans_uid_deleted_asset_id) INDEX(IDX_invest_trans_uid_deleted_account_id) NOT NULL"`
	AssetId              int64                     `xorm:"INDEX(IDX_invest_trans_uid_deleted_asset_id) NOT NULL"`
	AccountId            int64                     `xorm:"INDEX(IDX_invest_trans_uid_deleted_account_id) NOT NULL"`
	Type                 InvestmentTransactionType `xorm:"NOT NULL"`
	TradeTime            int64                     `xorm:"NOT NULL"`
	ConfirmTime          int64
	Quantity             int64 `xorm:"NOT NULL"`
	Price                int64 `xorm:"NOT NULL"`
	Amount               int64 `xorm:"NOT NULL"`
	Fee                  int64 `xorm:"NOT NULL DEFAULT 0"`
	RelatedTransactionId int64
	TimezoneUtcOffset    int16  `xorm:"NOT NULL"`
	Comment              string `xorm:"VARCHAR(255) NOT NULL"`
	CreatedUnixTime      int64
	UpdatedUnixTime      int64
	DeletedUnixTime      int64
}

// InvestmentTransactionListRequest represents all parameters of investment transaction listing request
type InvestmentTransactionListRequest struct {
	AssetId   int64                     `form:"asset_id,string"`
	AccountId int64                     `form:"account_id,string"`
	Type      InvestmentTransactionType `form:"type"`
	StartTime int64                     `form:"start_time"`
	EndTime   int64                     `form:"end_time"`
}

// InvestmentTransactionGetRequest represents all parameters of investment transaction getting request
type InvestmentTransactionGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// InvestmentTransactionCreateRequest represents all parameters of investment transaction creation request
type InvestmentTransactionCreateRequest struct {
	AssetId              int64                     `json:"assetId,string" binding:"required,min=1"`
	AccountId            int64                     `json:"accountId,string" binding:"required,min=1"`
	Type                 InvestmentTransactionType `json:"type" binding:"required,min=1"`
	TradeTime            int64                     `json:"tradeTime" binding:"required,min=1"`
	ConfirmTime          int64                     `json:"confirmTime"`
	Quantity             int64                     `json:"quantity" binding:"required,min=1"`
	Price                int64                     `json:"price" binding:"required,min=1"`
	Amount               int64                     `json:"amount" binding:"required,min=1"`
	Fee                  int64                     `json:"fee" binding:"min=0"`
	RelatedTransactionId int64                     `json:"relatedTransactionId,string"`
	TimezoneUtcOffset    int16                     `json:"utcOffset" binding:"min=-720,max=840"`
	Comment              string                    `json:"comment" binding:"max=255"`
	ClientSessionId      string                    `json:"clientSessionId"`
}

// InvestmentTransactionModifyRequest represents all parameters of investment transaction modification request
type InvestmentTransactionModifyRequest struct {
	Id                   int64                     `json:"id,string" binding:"required,min=1"`
	AssetId              int64                     `json:"assetId,string" binding:"required,min=1"`
	AccountId            int64                     `json:"accountId,string" binding:"required,min=1"`
	Type                 InvestmentTransactionType `json:"type" binding:"required,min=1"`
	TradeTime            int64                     `json:"tradeTime" binding:"required,min=1"`
	ConfirmTime          int64                     `json:"confirmTime"`
	Quantity             int64                     `json:"quantity" binding:"required,min=1"`
	Price                int64                     `json:"price" binding:"required,min=1"`
	Amount               int64                     `json:"amount" binding:"required,min=1"`
	Fee                  int64                     `json:"fee" binding:"min=0"`
	RelatedTransactionId int64                     `json:"relatedTransactionId,string"`
	TimezoneUtcOffset    int16                     `json:"utcOffset" binding:"min=-720,max=840"`
	Comment              string                    `json:"comment" binding:"max=255"`
}

// InvestmentTransactionDeleteRequest represents all parameters of investment transaction deleting request
type InvestmentTransactionDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// InvestmentTransactionInfoResponse represents a view-object of investment transaction
type InvestmentTransactionInfoResponse struct {
	Id                   int64                     `json:"id,string"`
	AssetId              int64                     `json:"assetId,string"`
	AccountId            int64                     `json:"accountId,string"`
	Type                 InvestmentTransactionType `json:"type"`
	TradeTime            int64                     `json:"tradeTime"`
	ConfirmTime          int64                     `json:"confirmTime,omitempty"`
	Quantity             int64                     `json:"quantity"`
	Price                int64                     `json:"price"`
	Amount               int64                     `json:"amount"`
	Fee                  int64                     `json:"fee"`
	RelatedTransactionId int64                     `json:"relatedTransactionId,string,omitempty"`
	TimezoneUtcOffset    int16                     `json:"utcOffset"`
	Comment              string                    `json:"comment"`
}

// ToInvestmentTransactionInfoResponse returns a view-object according to database model
func (t *InvestmentTransaction) ToInvestmentTransactionInfoResponse() *InvestmentTransactionInfoResponse {
	return &InvestmentTransactionInfoResponse{
		Id:                   t.TransactionId,
		AssetId:              t.AssetId,
		AccountId:            t.AccountId,
		Type:                 t.Type,
		TradeTime:            t.TradeTime,
		ConfirmTime:          t.ConfirmTime,
		Quantity:             t.Quantity,
		Price:                t.Price,
		Amount:               t.Amount,
		Fee:                  t.Fee,
		RelatedTransactionId: t.RelatedTransactionId,
		TimezoneUtcOffset:    t.TimezoneUtcOffset,
		Comment:              t.Comment,
	}
}
