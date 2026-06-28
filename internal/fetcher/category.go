package fetcher

import (
	"encoding/json"

	"lp_fetcher_golang/internal/models"
)

func FetchCategories(wdtoken string) models.CategoryFetchResult {
	param := map[string]interface{}{
		"shopId":    shopID,
		"attrQuery": []interface{}{},
		"from":      "h5",
	}

	reqURL, err := buildThorURL("decorate/itemCate.getCateTree/1.0", param, wdtoken, nil)
	if err != nil {
		return models.CategoryFetchResult{
			Status: "failed",
			Msg:    "未知错误",
		}
	}

	body, err := doGET(reqURL, shopReferer())
	if err != nil {
		return models.CategoryFetchResult{
			Status: "failed",
			Msg:    "未知错误",
		}
	}

	var cateResp models.CategoryResponse
	if err := json.Unmarshal(body, &cateResp); err != nil {
		return models.CategoryFetchResult{
			Status: "failed",
			Msg:    "未知错误",
		}
	}

	if cateResp.Status.Code != 0 {
		msg := cateResp.Status.Message
		if msg == "" {
			msg = "未知错误"
		}
		return models.CategoryFetchResult{
			Status: "failed",
			Msg:    msg,
		}
	}

	cateList := cateResp.BriefCateList()
	return models.CategoryFetchResult{
		Count:    len(cateList),
		Msg:      "",
		Status:   "success",
		CateList: cateList,
	}
}
