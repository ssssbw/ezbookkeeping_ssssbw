package marketdata

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// MarketDataProvider defines the structure of market data provider
type MarketDataProvider interface {
	GetDataSourceName() string
	// GetLatestPrice returns the latest confirmed price for a specific asset
	GetLatestPrice(c core.Context, assetCode string, market string) (*models.MarketData, error)

	// GetRealtimeEstimate returns the realtime estimated price (nil if not available, e.g. QDII)
	GetRealtimeEstimate(c core.Context, assetCode string, market string) (*models.MarketData, error)

	// GetHistoricalPrices returns historical prices for a specific asset
	GetHistoricalPrices(c core.Context, assetCode string, market string, startTime int64, endTime int64) ([]*models.MarketData, error)

	// GetAllFundNames returns all fund names for search functionality
	GetAllFundNames(c core.Context) (map[string]string, error)
}
