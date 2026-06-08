package marketdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/BlakeLiAFK/akshare/fund"
	"github.com/BlakeLiAFK/akshare/stock"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/log"
)

type EastMoneyAssetListProvider struct {
	httpClient *http.Client
}

func NewEastMoneyAssetListProvider() *EastMoneyAssetListProvider {
	return &EastMoneyAssetListProvider{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func mapFundType(fundType string) string {
	if strings.Contains(fundType, "债") || strings.Contains(fundType, "货币") {
		return "fixed_income"
	}
	if strings.Contains(fundType, "商品") {
		return "commodity"
	}
	return "equity"
}

func (p *EastMoneyAssetListProvider) ListFunds(ctx core.Context) ([]*AssetMeta, error) {
	url := "http://fund.eastmoney.com/js/fundcode_search.js"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Referer", "http://fund.eastmoney.com/")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		log.Errorf(ctx, "[eastmoney.ListFunds] failed to request %s: %s", url, err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	content := string(body)
	idx := strings.Index(content, "[")
	if idx < 0 {
		return nil, fmt.Errorf("no JSON array found in response")
	}
	content = content[idx:]
	if strings.HasSuffix(content, ";") {
		content = content[:len(content)-1]
	}

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
			Market:      1,
			Category:    mapFundType(item[3]),
			SubCategory: item[3],
		})
	}

	return result, nil
}

func (p *EastMoneyAssetListProvider) ListStocks(ctx core.Context) ([]*AssetMeta, error) {
	stocks, err := stock.StockInfoACodeNameEm()
	if err != nil {
		log.Errorf(ctx, "[eastmoney.ListStocks] failed: %s", err.Error())
		return nil, err
	}

	var result []*AssetMeta
	for _, s := range stocks {
		result = append(result, &AssetMeta{
			Code:     s.Code,
			Name:     s.Name,
			Market:   1,
			Category: "equity",
		})
	}

	return result, nil
}

func (p *EastMoneyAssetListProvider) ListETFs(ctx core.Context) ([]*AssetMeta, error) {
	etfs, err := fund.FundEtfSpotEm()
	if err != nil {
		log.Errorf(ctx, "[eastmoney.ListETFs] failed: %s", err.Error())
		return nil, err
	}

	var result []*AssetMeta
	for _, e := range etfs {
		code, _ := e["代码"].(string)
		name, _ := e["名称"].(string)
		if code == "" || name == "" {
			continue
		}
		result = append(result, &AssetMeta{
			Code:     code,
			Name:     name,
			Market:   1,
			Category: "equity",
		})
	}

	return result, nil
}
