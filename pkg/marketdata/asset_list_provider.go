package marketdata

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
)

// AssetMeta represents metadata for an asset fetched from external data source
type AssetMeta struct {
	Code        string
	Name        string
	Market      byte   // 1=CN, 2=HK, 3=US
	Category    string // "equity", "fixed_income", "commodity", "digital"
	SubCategory string // 原始分类文本，如"混合型-灵活"，直接显示不做映射
}

// AssetListProvider is the interface for fetching asset lists from external data sources
type AssetListProvider interface {
	ListFunds(ctx core.Context) ([]*AssetMeta, error)
	ListStocks(ctx core.Context) ([]*AssetMeta, error)
	ListETFs(ctx core.Context) ([]*AssetMeta, error)
}
