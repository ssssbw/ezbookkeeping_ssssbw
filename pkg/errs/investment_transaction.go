package errs

import "net/http"

// Error codes related to investment transactions
var (
	ErrInvestmentTransactionIdInvalid      = NewNormalError(NormalSubcategoryInvestmentTransaction, 0, http.StatusBadRequest, "investment transaction id is invalid")
	ErrInvestmentTransactionNotFound       = NewNormalError(NormalSubcategoryInvestmentTransaction, 1, http.StatusBadRequest, "investment transaction not found")
	ErrInvestmentTransactionTypeInvalid    = NewNormalError(NormalSubcategoryInvestmentTransaction, 2, http.StatusBadRequest, "investment transaction type is invalid")
	ErrInvestmentTransactionAssetIdInvalid = NewNormalError(NormalSubcategoryInvestmentTransaction, 3, http.StatusBadRequest, "investment transaction asset id is invalid")
	ErrInvestmentTransactionAccountIdInvalid = NewNormalError(NormalSubcategoryInvestmentTransaction, 4, http.StatusBadRequest, "investment transaction account id is invalid")
	ErrInvestmentTransactionAmountInvalid  = NewNormalError(NormalSubcategoryInvestmentTransaction, 5, http.StatusBadRequest, "investment transaction amount is invalid")
	ErrInvestmentTransactionQuantityInvalid = NewNormalError(NormalSubcategoryInvestmentTransaction, 6, http.StatusBadRequest, "investment transaction quantity is invalid")
)
