package cron

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/services"
)

// RemoveExpiredTokensJob represents the cron job which periodically remove expired user tokens from the database
var RemoveExpiredTokensJob = &CronJob{
	Name:        "RemoveExpiredTokens",
	Description: "Periodically remove expired user tokens from the database.",
	Period: CronJobFixedHourPeriod{
		Hour: 0,
	},
	Run: func(c *core.CronContext) error {
		return services.Tokens.DeleteAllExpiredTokens(c)
	},
}

// CreateScheduledTransactionJob represents the cron job which periodically create transaction by scheduled transaction template
var CreateScheduledTransactionJob = &CronJob{
	Name:        "CreateScheduledTransaction",
	Description: "Periodically create transaction by scheduled transaction template.",
	Period: CronJobEvery15MinutesPeriod{
		Second: 0,
	},
	Run: func(c *core.CronContext) error {
		return services.Transactions.CreateScheduledTransactions(c, time.Now().Unix(), c.GetInterval())
	},
}

// FetchMarketDataJob represents the cron job which periodically fetch market data for active investment assets
var FetchMarketDataJob = &CronJob{
	Name:        "FetchMarketData",
	Description: "Daily fetch market data for active investment assets.",
	Period: CronJobFixedHourPeriod{
		Hour: 21,
	},
	Run: func(c *core.CronContext) error {
		return services.MarketData.FetchAllActiveAssetsMarketData(c)
	},
}

// SyncAssetsJob represents the cron job which periodically sync global assets from external data source
var SyncAssetsJob = &CronJob{
	Name:        "SyncAssets",
	Description: "Quarterly sync global assets from external data source.",
	Period: CronJobCronPeriod{
		Expression: "0 2 1 1,4,7,10 *",
	},
	Run: func(c *core.CronContext) error {
		_, err := services.AssetSync.SyncAllAssets(c)
		return err
	},
}
