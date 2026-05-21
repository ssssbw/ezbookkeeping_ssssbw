package models

// InvestmentAssetType represents investment asset type
type InvestmentAssetType byte

// Investment asset types
const (
	INVESTMENT_ASSET_TYPE_FUND   InvestmentAssetType = 1
	INVESTMENT_ASSET_TYPE_STOCK  InvestmentAssetType = 2
	INVESTMENT_ASSET_TYPE_ETF    InvestmentAssetType = 3
	INVESTMENT_ASSET_TYPE_BOND   InvestmentAssetType = 4
	INVESTMENT_ASSET_TYPE_CRYPTO InvestmentAssetType = 5
)

// String returns a textual representation of the investment asset type
func (t InvestmentAssetType) String() string {
	switch t {
	case INVESTMENT_ASSET_TYPE_FUND:
		return "Fund"
	case INVESTMENT_ASSET_TYPE_STOCK:
		return "Stock"
	case INVESTMENT_ASSET_TYPE_ETF:
		return "ETF"
	case INVESTMENT_ASSET_TYPE_BOND:
		return "Bond"
	case INVESTMENT_ASSET_TYPE_CRYPTO:
		return "Crypto"
	default:
		return "Unknown"
	}
}

// IsValid returns whether the investment asset type is valid
func (t InvestmentAssetType) IsValid() bool {
	switch t {
	case INVESTMENT_ASSET_TYPE_FUND, INVESTMENT_ASSET_TYPE_STOCK, INVESTMENT_ASSET_TYPE_ETF,
		INVESTMENT_ASSET_TYPE_BOND, INVESTMENT_ASSET_TYPE_CRYPTO:
		return true
	default:
		return false
	}
}

// InvestmentTransactionType represents investment transaction type
type InvestmentTransactionType byte

// Investment transaction types
const (
	INVESTMENT_TRANSACTION_TYPE_BUY               InvestmentTransactionType = 1
	INVESTMENT_TRANSACTION_TYPE_SELL              InvestmentTransactionType = 2
	INVESTMENT_TRANSACTION_TYPE_DIVIDEND_CASH     InvestmentTransactionType = 3
	INVESTMENT_TRANSACTION_TYPE_DIVIDEND_REINVEST InvestmentTransactionType = 4
	INVESTMENT_TRANSACTION_TYPE_SPLIT             InvestmentTransactionType = 5
	INVESTMENT_TRANSACTION_TYPE_CONVERSION_OUT    InvestmentTransactionType = 6
	INVESTMENT_TRANSACTION_TYPE_CONVERSION_IN     InvestmentTransactionType = 7
)

// String returns a textual representation of the investment transaction type
func (t InvestmentTransactionType) String() string {
	switch t {
	case INVESTMENT_TRANSACTION_TYPE_BUY:
		return "Buy"
	case INVESTMENT_TRANSACTION_TYPE_SELL:
		return "Sell"
	case INVESTMENT_TRANSACTION_TYPE_DIVIDEND_CASH:
		return "Dividend (Cash)"
	case INVESTMENT_TRANSACTION_TYPE_DIVIDEND_REINVEST:
		return "Dividend (Reinvest)"
	case INVESTMENT_TRANSACTION_TYPE_SPLIT:
		return "Split"
	case INVESTMENT_TRANSACTION_TYPE_CONVERSION_OUT:
		return "Conversion (Out)"
	case INVESTMENT_TRANSACTION_TYPE_CONVERSION_IN:
		return "Conversion (In)"
	default:
		return "Unknown"
	}
}

// IsValid returns whether the investment transaction type is valid
func (t InvestmentTransactionType) IsValid() bool {
	switch t {
	case INVESTMENT_TRANSACTION_TYPE_BUY, INVESTMENT_TRANSACTION_TYPE_SELL,
		INVESTMENT_TRANSACTION_TYPE_DIVIDEND_CASH, INVESTMENT_TRANSACTION_TYPE_DIVIDEND_REINVEST,
		INVESTMENT_TRANSACTION_TYPE_SPLIT, INVESTMENT_TRANSACTION_TYPE_CONVERSION_OUT,
		INVESTMENT_TRANSACTION_TYPE_CONVERSION_IN:
		return true
	default:
		return false
	}
}

// IsBuyOrSell returns whether the transaction type is buy or sell
func (t InvestmentTransactionType) IsBuyOrSell() bool {
	return t == INVESTMENT_TRANSACTION_TYPE_BUY || t == INVESTMENT_TRANSACTION_TYPE_SELL
}

// IsDividend returns whether the transaction type is dividend
func (t InvestmentTransactionType) IsDividend() bool {
	return t == INVESTMENT_TRANSACTION_TYPE_DIVIDEND_CASH || t == INVESTMENT_TRANSACTION_TYPE_DIVIDEND_REINVEST
}

// InvestmentMarket represents investment market
type InvestmentMarket byte

// Investment markets
const (
	INVESTMENT_MARKET_CN InvestmentMarket = 1
	INVESTMENT_MARKET_HK InvestmentMarket = 2
	INVESTMENT_MARKET_US InvestmentMarket = 3
)

// String returns a textual representation of the investment market
func (m InvestmentMarket) String() string {
	switch m {
	case INVESTMENT_MARKET_CN:
		return "CN"
	case INVESTMENT_MARKET_HK:
		return "HK"
	case INVESTMENT_MARKET_US:
		return "US"
	default:
		return "Unknown"
	}
}

// IsValid returns whether the investment market is valid
func (m InvestmentMarket) IsValid() bool {
	switch m {
	case INVESTMENT_MARKET_CN, INVESTMENT_MARKET_HK, INVESTMENT_MARKET_US:
		return true
	default:
		return false
	}
}
