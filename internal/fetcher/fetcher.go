package fetcher

import (
	"encoding/json"
	"math/rand"
	"time"

	"lp_fetcher_golang/internal/models"
)

const (
	pageSize       = 20
	searchMaxPages = 3
)

type TooManyResultsError struct {
	Message string
}

func (e *TooManyResultsError) Error() string {
	return e.Message
}

func fetchAllPages(fetchPage func(offset int) ([]models.ItemBrief, error), maxPages int) ([]models.ItemBrief, error) {
	var items []models.ItemBrief
	offset := 0
	callCount := 0

	for {
		if maxPages > 0 && callCount >= maxPages {
			return nil, &TooManyResultsError{Message: "查询结果过多、请精确查询关键词"}
		}
		callCount++

		itemsOnce, err := fetchPage(offset)
		if err != nil {
			return nil, err
		}
		if len(itemsOnce) == 0 {
			break
		}
		items = append(items, itemsOnce...)
		offset += pageSize
		time.Sleep(time.Duration(rand.Intn(2000)+1000) * time.Millisecond)
	}

	return items, nil
}

func buildFetchResult(items []models.ItemBrief, err error) models.FetchResult {
	if err != nil {
		if _, ok := err.(*TooManyResultsError); ok {
			return models.FetchResult{
				Count:  0,
				Msg:    err.Error(),
				Status: "failed",
				Items:  nil,
			}
		}
		return models.FetchResult{
			Count:  0,
			Msg:    "未知错误",
			Status: "failed",
			Items:  nil,
		}
	}

	return models.FetchResult{
		Count:  len(items),
		Msg:    "",
		Status: "success",
		Items:  items,
	}
}

func getItemsOnce(keyword, wdtoken, dash string, offset int) ([]models.ItemBrief, error) {
	param := map[string]interface{}{
		"shopId":    shopID,
		"key":       keyword,
		"offset":    offset,
		"limit":     pageSize,
		"sortId":    0,
		"sortOrder": "desc",
		"from":      "h5",
	}

	reqURL, err := buildThorURL(
		"decorate/search.itemList/1.0",
		param,
		wdtoken,
		map[string]string{"_": dash},
	)
	if err != nil {
		return nil, err
	}

	body, err := doGET(reqURL, "https://h5.weidian.com/")
	if err != nil {
		return nil, err
	}

	var searchResp models.SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, err
	}

	return searchResp.BriefItems(), nil
}

func getItemsAll(keyword, wdtoken, dash string) ([]models.ItemBrief, error) {
	return fetchAllPages(func(offset int) ([]models.ItemBrief, error) {
		return getItemsOnce(keyword, wdtoken, dash, offset)
	}, searchMaxPages)
}

func getCateItemsOnce(cateID, wdtoken, dash string, offset int) ([]models.ItemBrief, error) {
	param := map[string]interface{}{
		"cateId":                 cateID,
		"shopId":                 shopID,
		"offset":                 offset,
		"limit":                  pageSize,
		"sortField":              "all",
		"sortType":               "desc",
		"isQdFx":                 false,
		"isHideSold":             true,
		"hideItemRealAmount":     true,
		"from":                   "h5",
		"fanSpreadMode":          1,
		"isShopItemListConfOpen": true,
		"attrQuery":              []interface{}{},
		"isStockDown":            0,
		"isConsumerProtect":      true,
		"hideItemComment":        false,
	}

	reqURL, err := buildThorURL(
		"decorate/itemCate.getCateItemList/1.0",
		param,
		wdtoken,
		map[string]string{"_": dash},
	)
	if err != nil {
		return nil, err
	}

	body, err := doGET(reqURL, shopReferer())
	if err != nil {
		return nil, err
	}

	var searchResp models.SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, err
	}

	return searchResp.BriefItems(), nil
}

func getCateItemsAll(cateID, wdtoken, dash string) ([]models.ItemBrief, error) {
	return fetchAllPages(func(offset int) ([]models.ItemBrief, error) {
		return getCateItemsOnce(cateID, wdtoken, dash, offset)
	}, 0)
}

func FetchItems(keyword, wdtoken, dash string) models.FetchResult {
	items, err := getItemsAll(keyword, wdtoken, dash)
	return buildFetchResult(items, err)
}

func FetchItemsByCategory(cateID, wdtoken, dash string) models.FetchResult {
	items, err := getCateItemsAll(cateID, wdtoken, dash)
	return buildFetchResult(items, err)
}
