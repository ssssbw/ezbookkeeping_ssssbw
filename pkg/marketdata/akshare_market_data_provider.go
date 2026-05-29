package marketdata

import (
	"fmt"
	"strconv"
	"time"

	"github.com/BlakeLiAFK/akshare/fund"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

func mustAtoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

type AkshareMarketDataProvider struct {
}

func NewAkshareMarketDataProvider() *AkshareMarketDataProvider {
	return &AkshareMarketDataProvider{}
}

func (p *AkshareMarketDataProvider) GetLatestPrice(c core.Context, assetCode string, market string) (*models.MarketData, error) {
	records, err := fund.FundOpenFundInfoEm(assetCode, "单位净值走势", "2026")
	if err != nil {
		log.Errorf(c, "[akshare.GetLatestPrice] failed to get fund info for %s: %s", assetCode, err.Error())
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no data found for asset %s", assetCode)
	}

	latest := records[len(records)-1]

	dateStr, ok := latest["净值日期"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid date format for asset %s", assetCode)
	}

	priceStr, ok := latest["单位净值"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid price format for asset %s", assetCode)
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date %s for asset %s", dateStr, assetCode)
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid price %s for asset %s", priceStr, assetCode)
	}

	return &models.MarketData{
		AssetId: 0,
		Date:    date.Unix(),
		Price:   int64(price * 10000),
	}, nil
}

func (p *AkshareMarketDataProvider) GetHistoricalPrices(c core.Context, assetCode string, market string, startTime int64, endTime int64) ([]*models.MarketData, error) {
	startYear := time.Unix(startTime, 0).Format("2006")
	endYear := time.Unix(endTime, 0).Format("2006")

	var allRecords []map[string]interface{}

	for year := startYear; year <= endYear; year = fmt.Sprintf("%d", 1+mustAtoi(year)) {
		records, err := fund.FundOpenFundInfoEm(assetCode, "单位净值走势", year)
		if err != nil {
			log.Errorf(c, "[akshare.GetHistoricalPrices] failed to get fund info for %s year %s: %s", assetCode, year, err.Error())
			continue
		}
		allRecords = append(allRecords, records...)
	}

	var result []*models.MarketData
	for _, record := range allRecords {
		dateStr, ok := record["净值日期"].(string)
		if !ok {
			continue
		}

		priceStr, ok := record["单位净值"].(string)
		if !ok {
			continue
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		dateUnix := date.Unix()
		if dateUnix < startTime || dateUnix > endTime {
			continue
		}

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			continue
		}

		result = append(result, &models.MarketData{
			AssetId: 0,
			Date:    dateUnix,
			Price:   int64(price * 10000),
		})
	}

	return result, nil
}

func (p *AkshareMarketDataProvider) GetAllFundNames(c core.Context) (map[string]string, error) {
	records, err := fund.FundNameEm()
	if err != nil {
		log.Errorf(c, "[akshare.GetAllFundNames] failed to get fund names: %s", err.Error())
		return nil, err
	}

	result := make(map[string]string)
	for _, record := range records {
		code, ok := record["基金代码"].(string)
		if !ok {
			continue
		}

		name, ok := record["基金简称"].(string)
		if !ok {
			continue
		}

		result[code] = name
	}

	return result, nil
}
