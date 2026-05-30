package marketdata

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/BlakeLiAFK/akshare/fund"
	"github.com/dop251/goja"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

type AkshareMarketDataProvider struct {
	httpClient *http.Client
}

func (p *AkshareMarketDataProvider) GetDataSourceName() string {
	return "akshare"
}

func NewAkshareMarketDataProvider() *AkshareMarketDataProvider {
	return &AkshareMarketDataProvider{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *AkshareMarketDataProvider) GetRealtimeEstimate(c core.Context, assetCode string, market string) (*models.MarketData, error) {
	return nil, nil
}

func (p *AkshareMarketDataProvider) GetLatestPrice(c core.Context, assetCode string, market string) (*models.MarketData, error) {
	url := fmt.Sprintf("https://fund.eastmoney.com/pingzhongdata/%s.js", assetCode)
	resp, err := p.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	vm := goja.New()
	_, err = vm.RunString(string(body))
	if err != nil {
		return nil, err
	}

	dataValue := vm.Get("Data_netWorthTrend")
	if dataValue == nil || goja.IsUndefined(dataValue) {
		return nil, fmt.Errorf("no data for asset %s", assetCode)
	}

	dataExport := dataValue.Export()
	dataArray, ok := dataExport.([]interface{})
	if !ok || len(dataArray) == 0 {
		return nil, fmt.Errorf("invalid data for asset %s", assetCode)
	}

	latest := dataArray[len(dataArray)-1]
	itemMap, ok := latest.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid record for asset %s", assetCode)
	}

	// x 是 int64 时间戳（毫秒）
	xVal, ok := itemMap["x"].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid timestamp for asset %s", assetCode)
	}

	yVal, ok := itemMap["y"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid price for asset %s", assetCode)
	}

	date := time.UnixMilli(xVal).UTC()
	shanghaiLoc, _ := time.LoadLocation("Asia/Shanghai")
	date = date.In(shanghaiLoc)

	return &models.MarketData{
		AssetId: 0,
		Date:    date.Unix(),
		Price:   int64(yVal * 10000),
	}, nil
}

func (p *AkshareMarketDataProvider) GetHistoricalPrices(c core.Context, assetCode string, market string, startTime int64, endTime int64) ([]*models.MarketData, error) {
	records, err := fund.FundOpenFundInfoEm(assetCode, "单位净值走势", "成立来")
	if err != nil {
		return nil, err
	}

	var result []*models.MarketData
	for _, record := range records {
		priceVal, ok := record["单位净值"].(float64)
		if !ok {
			continue
		}

		result = append(result, &models.MarketData{
			AssetId: 0,
			Date:    time.Now().Unix(),
			Price:   int64(priceVal * 10000),
		})
	}

	return result, nil
}

func (p *AkshareMarketDataProvider) GetAllFundNames(c core.Context) (map[string]string, error) {
	records, err := fund.FundNameEm()
	if err != nil {
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
