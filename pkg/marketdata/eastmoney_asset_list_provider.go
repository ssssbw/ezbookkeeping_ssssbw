package marketdata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

// EastMoneyAssetListProvider provides asset list data from eastmoney
type EastMoneyAssetListProvider struct {
	httpClient *http.Client
}

// NewEastMoneyAssetListProvider returns a new EastMoneyAssetListProvider
func NewEastMoneyAssetListProvider() *EastMoneyAssetListProvider {
	return &EastMoneyAssetListProvider{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// fundTypeToCategory maps eastmoney fund types to asset categories
var fundTypeToCategory = map[string]string{
	"股票型":   "equity",
	"混合型":   "equity",
	"债券型":   "fixed_income",
	"指数型":   "equity",
	"货币型":   "fixed_income",
	"QDII":   "equity",
	"FOF":    "equity",
	"ETF联接基金": "equity",
	"LOF":    "equity",
}

// mapFundType returns the asset category for a given fund type
func mapFundType(fundType string) string {
	if category, ok := fundTypeToCategory[fundType]; ok {
		return category
	}
	return "equity"
}

// ListFunds fetches all fund names from eastmoney
func (p *EastMoneyAssetListProvider) ListFunds(ctx core.Context) ([]*AssetMeta, error) {
	url := "http://fund.eastmoney.com/js/fundcode_search.js"

	resp, err := p.httpClient.Get(url)
	if err != nil {
		log.Errorf(ctx, "[eastmoney.ListFunds] failed to request %s: %s", url, err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
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

	var result []*AssetMeta
	for _, item := range data {
		if len(item) < 4 {
			continue
		}
		result = append(result, &AssetMeta{
			Code:        item[0],
			Name:        item[2],
			Market:      1, // CN market
			Category:    mapFundType(item[3]),
			SubCategory: item[3], // 如"混合型-灵活"
		})
	}

	return result, nil
}

// eastmoneyListResponse represents the JSON response from eastmoney stock/ETF list API
type eastmoneyListResponse struct {
	Data struct {
		Total int `json:"total"`
		List  []struct {
			F12 string `json:"f12"` // code
			F13 string `json:"f13"` // market
			F14 string `json:"f14"` // name
		} `json:"list"`
	} `json:"data"`
}

// ListStocks fetches all Chinese stock codes from eastmoney
func (p *EastMoneyAssetListProvider) ListStocks(ctx core.Context) ([]*AssetMeta, error) {
	url := "https://82.push2.eastmoney.com/api/qt/clist/get?pn=1&pz=5000&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&fid=f12&fs=m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23,m:0+t:81+s:2048&fields=f12,f13,f14"

	resp, err := p.httpClient.Get(url)
	if err != nil {
		log.Errorf(ctx, "[eastmoney.ListStocks] failed to request %s: %s", url, err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var listResp eastmoneyListResponse
	if err := json.Unmarshal(body, &listResp); err != nil {
		return nil, err
	}

	var result []*AssetMeta
	for _, item := range listResp.Data.List {
		result = append(result, &AssetMeta{
			Code:     item.F12,
			Name:     item.F14,
			Market:   1, // CN market
			Category: "equity",
		})
	}

	return result, nil
}

// ListETFs fetches all Chinese ETF codes from eastmoney
func (p *EastMoneyAssetListProvider) ListETFs(ctx core.Context) ([]*AssetMeta, error) {
	url := "https://88.push2.eastmoney.com/api/qt/clist/get?pn=1&pz=5000&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&fid=f3&fs=b:MK0021,b:MK0022,b:MK0023,b:MK0024&fields=f12,f13,f14"

	resp, err := p.httpClient.Get(url)
	if err != nil {
		log.Errorf(ctx, "[eastmoney.ListETFs] failed to request %s: %s", url, err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var listResp eastmoneyListResponse
	if err := json.Unmarshal(body, &listResp); err != nil {
		return nil, err
	}

	var result []*AssetMeta
	for _, item := range listResp.Data.List {
		result = append(result, &AssetMeta{
			Code:     item.F12,
			Name:     item.F14,
			Market:   1, // CN market
			Category: "equity",
		})
	}

	return result, nil
}
