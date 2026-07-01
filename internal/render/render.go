package render

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"lp_fetcher_golang/internal/models"
)

const notFoundText = "未找到相关物品"

func LoadSearchResult(path string) (models.MaterialFetchResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return models.MaterialFetchResult{}, fmt.Errorf("未找到 JSON 文件: %s", path)
	}

	var result models.MaterialFetchResult
	if err := json.Unmarshal(data, &result); err != nil {
		return models.MaterialFetchResult{}, err
	}

	return result, nil
}

func findMaterialItems(groups []models.MaterialGroup, name string) []models.MaterialItem {
	for _, group := range groups {
		if group.Material == name {
			return group.Items
		}
	}
	return nil
}

func formatMaterialSection(items []models.MaterialItem) string {
	if len(items) == 0 {
		return notFoundText
	}

	var b strings.Builder
	for i, item := range items {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString("- **价格**：")
		b.WriteString(item.Price)
		b.WriteString("，**标题**：")
		b.WriteString(item.ItemName)
	}
	return b.String()
}

func RenderReport(result models.MaterialFetchResult, templatePath, outputPath string) error {
	template, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}

	vinylItems := findMaterialItems(result.Materials, "黑胶")
	cdItems := findMaterialItems(result.Materials, "CD")

	content := string(template)
	content = strings.ReplaceAll(content, "{{count}}", fmt.Sprintf("%d", result.Count))
	content = strings.ReplaceAll(content, "{{vinyl_items}}", formatMaterialSection(vinylItems))
	content = strings.ReplaceAll(content, "{{cd_items}}", formatMaterialSection(cdItems))

	return os.WriteFile(outputPath, []byte(content), 0644)
}
