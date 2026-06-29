# lp_fetcher

微店「某某唱片」商品检索 CLI 工具。支持按关键字搜索、按店铺分类浏览，并将结果渲染为 Markdown 报告。

## 功能

| 子命令 | 说明 |
|--------|------|
| `search` | 按关键字检索商品，结果写入 `search_result.json` |
| `search_by_category` | 按分类 ID 检索商品，结果写入 `search_result.json` |
| `category` | 获取店铺分类叶子列表，结果写入 `category_result.json` |
| `render` | 读取 `search_result.json`，生成 `report.md` |

## 环境要求

- Go 1.22+

## 快速开始

### 关键字检索

```bash
./lp_fetcher search --keyword "Highway To Hell"
./lp_fetcher render
```

打开 `report.md` 查看结果。

### 按分类检索

```bash
# 1. 获取分类叶子列表（首次或分类有变更时执行）
./lp_fetcher category

# 2. 用 match_category.py 匹配目标 cateId，再执行分类检索
python3 match_category.py "爵士"
./lp_fetcher search_by_category --cateId "<cateId>"
./lp_fetcher render
```

## 参数说明

### 子命令参数

**search**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--keyword` | — | 店铺内搜索词（必填） |
| `--output` | `search_result.json` | 输出 JSON 路径 |

**search_by_category**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--cateId` | — | 店铺分类 ID（必填） |
| `--output` | `search_result.json` | 输出 JSON 路径 |

**category**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--output` | `category_result.json` | 输出 JSON 路径 |

**render**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--input` | `search_result.json` | 输入 JSON 路径 |
| `--template` | `report_template.md` | Markdown 模板路径 |
| `--output` | `report.md` | 输出 Markdown 路径 |

## 输出文件

| 文件 | 说明 |
|------|------|
| `category_result.json` | 分类叶子列表，含一维 `cateList`（`cateId`、`catName`，已归一化） |
| `search_result.json` | 检索结果，含 `count`、`status`、`msg`、`materials`（按材质分组的数组） |
| `report.md` | 最终 Markdown 报告 |

`search_result.json` 示例：

```json
{
  "count": 2,
  "msg": "",
  "status": "success",
  "materials": [
    {
      "material": "黑胶",
      "items": [
        { "itemId": "123", "itemName": "Album LP", "price": "289" }
      ]
    },
    {
      "material": "CD",
      "items": [
        { "itemId": "456", "itemName": "Album CD", "price": "99" }
      ]
    }
  ]
}
```

`category_result.json` 示例：

```json
{
  "count": 2,
  "msg": "",
  "status": "success",
  "cateList": [
    { "cateId": 141851486, "catName": "cooljazz冷爵士" },
    { "cateId": 127630967, "catName": "爵士" }
  ]
}
```

### match_category.py

根据关键词在 `category_result.json` 中做子串匹配，返回 `{cateId, cateName}` 列表。

```bash
python3 match_category.py "cool jazz"
python3 match_category.py cooljazz -i category_result.json
```

输出示例：

```json
[
  { "cateId": 141851486, "cateName": "cooljazz冷爵士" }
]
```

## 限制与异常

- **关键字检索分页上限**：最多拉取 3 页（每页 20 条）。超出时 `status` 为 `failed`，`msg` 为「查询结果过多、请精确查询关键词」，需缩小搜索范围或改用分类检索。
- **分类检索**：无分页上限，自动翻页直至无更多数据。
- **请求间隔**：翻页时随机等待 1–3 秒，避免触发限流。
- **请求失败**：检查网络连接，或稍后重试。

## 项目结构

```text
lp_fetcher_golang/
├── main.go                    # 程序入口
├── report_template.md         # Markdown 报告模板
├── report_template.html       # 旧版 HTML 模板（保留，render 不再使用）
├── SKILL.md                   # Cursor Agent 技能说明
├── match_category.py          # 分类关键词匹配脚本
├── internal/
│   ├── cli/                   # Cobra 子命令（search、category、render 等）
│   ├── fetcher/               # 微店 API 请求与分页逻辑
│   ├── models/                # 数据模型
│   └── render/                # JSON → Markdown 渲染
└── bin/                       # 构建产物（gitignore）
```

## 技术栈

- Go 1.22
- [Cobra](https://github.com/spf13/cobra) — CLI 框架
- 微店 Thor API（`thor.weidian.com`）
