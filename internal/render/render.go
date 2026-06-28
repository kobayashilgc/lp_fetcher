package render

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"lp_fetcher_golang/internal/models"
)

const (
	beginItemRow = "<!-- BEGIN_ITEM_ROW -->"
	endItemRow   = "<!-- END_ITEM_ROW -->"
)

func LoadSearchResult(path string) (models.FetchResult, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return models.FetchResult{}, "", fmt.Errorf("未找到 JSON 文件: %s", path)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return models.FetchResult{}, "", err
	}

	var result models.FetchResult
	if err := json.Unmarshal(data, &result); err != nil {
		return models.FetchResult{}, "", err
	}

	keyword := ""
	if kw, ok := raw["keyword"]; ok {
		json.Unmarshal(kw, &keyword)
	}

	return result, keyword, nil
}

func RenderReport(result models.FetchResult, keyword, templatePath, outputPath string) error {
	template, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}
	templateStr := string(template)

	begin := strings.Index(templateStr, beginItemRow)
	end := strings.Index(templateStr, endItemRow)
	if begin < 0 || end < 0 || end <= begin {
		return fmt.Errorf("模板中未找到行标记")
	}

	rowTemplate := templateStr[begin+len(beginItemRow) : end]

	var rowsHTML strings.Builder
	for _, item := range result.Items {
		row := rowTemplate
		row = strings.ReplaceAll(row, "{{material}}", item.Material)
		row = strings.ReplaceAll(row, "{{price}}", item.Price)
		row = strings.ReplaceAll(row, "{{item_name}}", item.ItemName)
		rowsHTML.WriteString(row)
	}

	html := templateStr[:begin] + rowsHTML.String() + templateStr[end+len(endItemRow):]
	html = strings.ReplaceAll(html, "{{keyword}}", keyword)
	html = strings.ReplaceAll(html, "{{count}}", fmt.Sprintf("%d", result.Count))

	return os.WriteFile(outputPath, []byte(html), 0644)
}
