package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"lp_fetcher_golang/internal/models"
)

const testTemplate = `# 总数量

{{count}}

# 黑胶

{{vinyl_items}}

# CD

{{cd_items}}
`

func writeTemplate(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "template.md")
	if err := os.WriteFile(path, []byte(testTemplate), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadSearchResult(t *testing.T) {
	t.Run("valid json", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "result.json")
		data := `{"count":1,"msg":"","status":"success","materials":[{"material":"黑胶","items":[{"itemId":"1","itemName":"LP","price":"100"}]}]}`
		if err := os.WriteFile(path, []byte(data), 0644); err != nil {
			t.Fatal(err)
		}

		got, err := LoadSearchResult(path)
		if err != nil {
			t.Fatal(err)
		}
		if got.Count != 1 || len(got.Materials) != 1 {
			t.Fatalf("unexpected result: %+v", got)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		_, err := LoadSearchResult(filepath.Join(t.TempDir(), "missing.json"))
		if err == nil {
			t.Fatal("expected error for missing file")
		}
	})
}

func TestRenderReport(t *testing.T) {
	t.Run("vinyl and cd", func(t *testing.T) {
		dir := t.TempDir()
		template := writeTemplate(t, dir)
		output := filepath.Join(dir, "report.md")

		result := models.MaterialFetchResult{
			Count:  2,
			Status: "success",
			Materials: []models.MaterialGroup{
				{
					Material: "黑胶",
					Items:    []models.MaterialItem{{ItemName: "LP Album", Price: "289"}},
				},
				{
					Material: "CD",
					Items:    []models.MaterialItem{{ItemName: "CD Album", Price: "99"}},
				},
			},
		}

		if err := RenderReport(result, template, output); err != nil {
			t.Fatal(err)
		}

		content, err := os.ReadFile(output)
		if err != nil {
			t.Fatal(err)
		}
		text := string(content)
		if !strings.Contains(text, "# 总数量") {
			t.Fatalf("missing count heading: %s", text)
		}
		if !strings.Contains(text, "- **价格**：289，**标题**：LP Album") || !strings.Contains(text, "- **价格**：99，**标题**：CD Album") {
			t.Fatalf("missing list items: %s", text)
		}
	})

	t.Run("vinyl only", func(t *testing.T) {
		dir := t.TempDir()
		template := writeTemplate(t, dir)
		output := filepath.Join(dir, "report.md")

		result := models.MaterialFetchResult{
			Count:  1,
			Status: "success",
			Materials: []models.MaterialGroup{
				{
					Material: "黑胶",
					Items:    []models.MaterialItem{{ItemName: "LP Album", Price: "289"}},
				},
			},
		}

		if err := RenderReport(result, template, output); err != nil {
			t.Fatal(err)
		}

		content, err := os.ReadFile(output)
		if err != nil {
			t.Fatal(err)
		}
		text := string(content)
		if !strings.Contains(text, "- **价格**：289，**标题**：LP Album") {
			t.Fatalf("missing vinyl item: %s", text)
		}
		cdSection := text[strings.Index(text, "# CD"):]
		if !strings.Contains(cdSection, notFoundText) {
			t.Fatalf("CD section should contain not-found text: %s", cdSection)
		}
	})

	t.Run("no vinyl or cd", func(t *testing.T) {
		dir := t.TempDir()
		template := writeTemplate(t, dir)
		output := filepath.Join(dir, "report.md")

		result := models.MaterialFetchResult{
			Count:  1,
			Status: "success",
			Materials: []models.MaterialGroup{
				{
					Material: "其他",
					Items:    []models.MaterialItem{{ItemName: "Poster", Price: "50"}},
				},
			},
		}

		if err := RenderReport(result, template, output); err != nil {
			t.Fatal(err)
		}

		content, err := os.ReadFile(output)
		if err != nil {
			t.Fatal(err)
		}
		text := string(content)
		if strings.Count(text, notFoundText) != 2 {
			t.Fatalf("expected two not-found sections, got: %s", text)
		}
	})
}
