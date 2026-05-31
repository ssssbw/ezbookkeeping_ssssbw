package models

// InvestmentOverviewResponse represents aggregated portfolio summary
type InvestmentOverviewResponse struct {
	TotalInvestment    int64                       `json:"totalInvestment"`
	TotalMarketValue   int64                       `json:"totalMarketValue"`
	TotalUnrealizedPnl int64                       `json:"totalUnrealizedPnl"`
	TotalReturnRate    int64                       `json:"totalReturnRate"`
	Allocations        []*InvestmentAllocationItem `json:"allocations"`
}

// InvestmentAllocationItem represents per-category allocation
type InvestmentAllocationItem struct {
	Category   string `json:"category"`
	Value      int64  `json:"value"`
	Percentage int64  `json:"percentage"`
}

// InvestmentHoldingInfo represents per-asset holding with computed metrics
type InvestmentHoldingInfo struct {
	AssetId       int64            `json:"assetId,string"`
	AssetCode     string           `json:"assetCode"`
	AssetName     string           `json:"assetName"`
	Category      AssetCategory    `json:"category"`
	Currency      string           `json:"currency"`
	Market        InvestmentMarket `json:"market"`
	AccountId     int64            `json:"accountId,string"`
	Quantity      int64            `json:"quantity"`
	AvgCostPrice  int64            `json:"avgCostPrice"`
	TotalCost     int64            `json:"totalCost"`
	CurrentPrice  int64            `json:"currentPrice"`
	MarketValue   int64            `json:"marketValue"`
	UnrealizedPnl int64            `json:"unrealizedPnl"`
	ReturnRate    int64            `json:"returnRate"`
}
