package marketdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

type EastMoneyMarketDataProvider struct {
	httpClient *http.Client
}

func (p *EastMoneyMarketDataProvider) GetDataSourceName() string {
	return "eastmoney"
}

func NewEastMoneyMarketDataProvider() *EastMoneyMarketDataProvider {
	return &EastMoneyMarketDataProvider{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *EastMoneyMarketDataProvider) GetRealtimeEstimate(c core.Context, assetCode string, market string) (*models.MarketData, error) {
	url := fmt.Sprintf("https://fundgz.1234567.com.cn/js/%s.js", assetCode)

	resp, err := p.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jsonStr := strings.TrimPrefix(string(body), "jsonpgz(")
	jsonStr = strings.TrimSuffix(jsonStr, ");")

	if jsonStr == "" || jsonStr == "null" {
		return nil, nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, nil
	}

	gsz, ok := data["gsz"].(string)
	if !ok || gsz == "" {
		return nil, nil
	}

	price, err := strconv.ParseFloat(gsz, 64)
	if err != nil {
		return nil, nil
	}

	gztime, _ := data["gztime"].(string)
	var estimateTime int64
	if gztime != "" {
		if t, err := time.ParseInLocation("2006-01-02 15:04", gztime, time.Local); err == nil {
			estimateTime = t.Unix()
		}
	}

	if estimateTime == 0 {
		estimateTime = time.Now().Unix()
	}

	return &models.MarketData{
		AssetId: 0,
		Date:    estimateTime,
		Price:   int64(price * 10000),
	}, nil
}

func (p *EastMoneyMarketDataProvider) GetLatestPrice(c core.Context, assetCode string, market string) (*models.MarketData, error) {
	url := fmt.Sprintf("https://fundgz.1234567.com.cn/js/%s.js", assetCode)

	resp, err := p.httpClient.Get(url)
	if err != nil {
		log.Errorf(c, "[eastmoney.GetLatestPrice] failed to request %s: %s", url, err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jsonStr := strings.TrimPrefix(string(body), "jsonpgz(")
	jsonStr = strings.TrimSuffix(jsonStr, ");")

	if jsonStr == "" || jsonStr == "null" {
		return nil, fmt.Errorf("no data available for asset %s (possibly QDII or special fund)", assetCode)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("failed to parse response for asset %s: %s", assetCode, err.Error())
	}

	dateStr, ok := data["jzrq"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid date for asset %s", assetCode)
	}

	priceStr, ok := data["dwjz"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid price for asset %s", assetCode)
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return nil, err
	}

	return &models.MarketData{
		AssetId: 0,
		Date:    date.Unix(),
		Price:   int64(price * 10000),
	}, nil
}

func (p *EastMoneyMarketDataProvider) GetHistoricalPrices(c core.Context, assetCode string, market string, startTime int64, endTime int64) ([]*models.MarketData, error) {
	startDate := time.Unix(startTime, 0).Format("2006-01-02")
	endDate := time.Unix(endTime, 0).Format("2006-01-02")

	log.Debugf(c, "[eastmoney.GetHistoricalPrices] fetching prices for asset %s from %s to %s", assetCode, startDate, endDate)

	var allResults []*models.MarketData
	page := 1

	for {
		url := fmt.Sprintf("http://fund.eastmoney.com/f10/F10DataApi.aspx?type=lsjz&code=%s&page=%d&sdate=%s&edate=%s&per=49", assetCode, page, startDate, endDate)

		log.Debugf(c, "[eastmoney.GetHistoricalPrices] requesting URL: %s", url)

		resp, err := p.httpClient.Get(url)
		if err != nil {
			log.Errorf(c, "[eastmoney.GetHistoricalPrices] failed to request %s: %s", url, err.Error())
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		html := string(body)

		re := regexp.MustCompile(`<td[^>]*>([^<]*)</td>`)
		matches := re.FindAllStringSubmatch(html, -1)

		if len(matches) < 7 {
			break
		}

		pageCount := 0
		for i := 0; i < len(matches)-6; i += 7 {
			dateStr := strings.TrimSpace(matches[i][1])
			priceStr := strings.TrimSpace(matches[i+1][1])

			// Parse date as UTC to avoid timezone issues
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

			allResults = append(allResults, &models.MarketData{
				AssetId: 0,
				Date:    dateUnix,
				Price:   int64(price * 10000),
			})
			pageCount++
		}

		log.Debugf(c, "[eastmoney.GetHistoricalPrices] page %d: got %d records for asset %s", page, pageCount, assetCode)

		if pageCount == 0 {
			break
		}
		page++
	}

	log.Debugf(c, "[eastmoney.GetHistoricalPrices] total %d records for asset %s", len(allResults), assetCode)
	return allResults, nil
}

func (p *EastMoneyMarketDataProvider) GetAllFundNames(c core.Context) (map[string]string, error) {
	url := "http://fund.eastmoney.com/js/fundcode_search.js"

	resp, err := p.httpClient.Get(url)
	if err != nil {
		log.Errorf(c, "[eastmoney.GetAllFundNames] failed to request %s: %s", url, err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	content := string(body)
	content = strings.TrimPrefix(content, "var r = ")
	content = strings.TrimSuffix(content, ";")

	var data [][]string
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, item := range data {
		if len(item) >= 3 {
			result[item[0]] = item[2]
		}
	}

	return result, nil
}
