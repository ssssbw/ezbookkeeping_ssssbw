package errs

import "net/http"

// Error codes related to investment assets
var (
	ErrInvestmentAssetIdInvalid    = NewNormalError(NormalSubcategoryInvestmentAsset, 0, http.StatusBadRequest, "investment asset id is invalid")
	ErrInvestmentAssetNotFound     = NewNormalError(NormalSubcategoryInvestmentAsset, 1, http.StatusBadRequest, "investment asset not found")
	ErrInvestmentAssetTypeInvalid  = NewNormalError(NormalSubcategoryInvestmentAsset, 2, http.StatusBadRequest, "investment asset type is invalid")
	ErrInvestmentAssetMarketInvalid = NewNormalError(NormalSubcategoryInvestmentAsset, 3, http.StatusBadRequest, "investment asset market is invalid")
	ErrInvestmentAssetInUseCannotBeDeleted = NewNormalError(NormalSubcategoryInvestmentAsset, 4, http.StatusBadRequest, "investment asset is in use and cannot be deleted")
	ErrInvestmentAssetCodeAlreadyExists = NewNormalError(NormalSubcategoryInvestmentAsset, 5, http.StatusBadRequest, "investment asset code already exists")
)
