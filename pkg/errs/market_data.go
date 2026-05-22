package errs

import "net/http"

// Error codes related to market data
var (
	ErrMarketDataIdInvalid    = NewNormalError(NormalSubcategoryMarketData, 0, http.StatusBadRequest, "market data id is invalid")
	ErrMarketDataNotFound     = NewNormalError(NormalSubcategoryMarketData, 1, http.StatusBadRequest, "market data not found")
	ErrMarketDataAssetIdInvalid = NewNormalError(NormalSubcategoryMarketData, 2, http.StatusBadRequest, "market data asset id is invalid")
	ErrMarketDataDateInvalid  = NewNormalError(NormalSubcategoryMarketData, 3, http.StatusBadRequest, "market data date is invalid")
	ErrMarketDataPriceInvalid = NewNormalError(NormalSubcategoryMarketData, 4, http.StatusBadRequest, "market data price is invalid")
)
