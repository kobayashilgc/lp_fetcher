package models

import "strings"

type Status struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

type ItemExtend struct {
	HotItemRank        *int    `json:"hotItemRank"`
	InterestFreeText   *string `json:"interestFreeText"`
	Knowledge          *string `json:"knowledge"`
	KnowledgeItemType  *string `json:"knowledgeItemType"`
	KnowledgeReItemCnt *int    `json:"knowledgeReItemCnt"`
}

type ItemTag struct {
	TagTitle string `json:"tagTitle"`
	TagType  int    `json:"tagType"`
}

type Item struct {
	AddTime              int         `json:"addTime"`
	CreditJoin           bool        `json:"creditJoin"`
	Discount             *string     `json:"discount"`
	Extend               ItemExtend  `json:"extend"`
	HasSku               bool        `json:"hasSku"`
	ItemComment          string      `json:"itemComment"`
	ItemID               string      `json:"itemId"`
	ItemImg              string      `json:"itemImg"`
	ItemName             string      `json:"itemName"`
	ItemTag              []ItemTag   `json:"itemTag"`
	ItemURL              string      `json:"itemUrl"`
	OriginalPrice        *string     `json:"originalPrice"`
	PreSale              bool        `json:"preSale"`
	Price                string      `json:"price"`
	PriceTag             *string     `json:"priceTag"`
	PromotionCommission  *string     `json:"promotionCommission"`
	Sold                 string      `json:"sold"`
	StartSaleLimit       int         `json:"startSaleLimit"`
	Status               int         `json:"status"`
	Stock                int         `json:"stock"`
	TagDO                *string     `json:"tagDO"`
	TakePrice            *string     `json:"takePrice"`
}

type SearchResult struct {
	AllSuccess bool   `json:"allSuccess"`
	HasData    bool   `json:"hasData"`
	ItemList   []Item `json:"itemList"`
}

type ItemBrief struct {
	ItemID    string `json:"itemId"`
	ItemName  string `json:"itemName"`
	Price     string `json:"price"`
	Material  string `json:"material"`
}

type FetchResult struct {
	Count  int         `json:"count"`
	Msg    string      `json:"msg"`
	Status string      `json:"status"`
	Items  []ItemBrief `json:"items"`
}

type SearchResponse struct {
	Status Status       `json:"status"`
	Result SearchResult `json:"result"`
}

func InferMaterial(itemName string) string {
	if strings.Contains(itemName, "CD") {
		return "CD"
	}
	if strings.Contains(itemName, "LP") || strings.Contains(itemName, "胶") {
		return "黑胶"
	}
	return "其他"
}

func (r *SearchResponse) BriefItems() []ItemBrief {
	items := make([]ItemBrief, 0, len(r.Result.ItemList))
	for _, item := range r.Result.ItemList {
		items = append(items, ItemBrief{
			ItemID:   item.ItemID,
			ItemName: item.ItemName,
			Price:    item.Price,
			Material: InferMaterial(item.ItemName),
		})
	}
	return items
}

type CategoryRaw struct {
	CateID        int           `json:"cateId"`
	CateName      string        `json:"cateName"`
	ChildCateList []CategoryRaw `json:"childCateList"`
}

type CategoryBrief struct {
	CateID        int             `json:"cateId"`
	CateName      string          `json:"cateName"`
	ChildCateList []CategoryBrief `json:"childCateList"`
}

type CategoryFetchResult struct {
	Count    int             `json:"count"`
	Msg      string          `json:"msg"`
	Status   string          `json:"status"`
	CateList []CategoryBrief `json:"cateList"`
}

type CategoryResult struct {
	CateList []CategoryRaw `json:"cateList"`
	HasData  bool          `json:"hasData"`
}

type CategoryResponse struct {
	Status Status         `json:"status"`
	Result CategoryResult `json:"result"`
}

func briefCate(raw CategoryRaw) CategoryBrief {
	children := make([]CategoryBrief, 0, len(raw.ChildCateList))
	for _, child := range raw.ChildCateList {
		children = append(children, briefCate(child))
	}
	return CategoryBrief{
		CateID:        raw.CateID,
		CateName:      raw.CateName,
		ChildCateList: children,
	}
}

func (r *CategoryResponse) BriefCateList() []CategoryBrief {
	cates := make([]CategoryBrief, 0, len(r.Result.CateList))
	for _, raw := range r.Result.CateList {
		cates = append(cates, briefCate(raw))
	}
	return cates
}
