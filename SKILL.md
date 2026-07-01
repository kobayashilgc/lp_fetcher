---
name: raccoon-records-search
description: 根据关键字或分类/风格检索微店「浣熊唱片」商品并生成 report.md 报告。支持关键字检索与按分类检索两种模式。关键字模式需提供 keyword，分类模式需用户确认分类 ID。报告输出至 resources/report.md。当用户提到浣熊唱片、唱片检索、lp_fetcher、按分类/风格检索或需要查询该店商品时使用。
version: 20260629001
---

# 浣熊唱片商品检索

根据关键字或分类检索微店「浣熊唱片」（shopId: 1404154952）的商品情况。

## 全局约束

- 工作目录为技能目录
- **仅允许**调用以下可执行文件：
  - `./scripts/lp_fetcher <子命令> [参数]` — 微店 API 检索
  - `python3 scripts/match_category.py "<关键词>" -i resources/category_result.json` — 分类匹配（场景二 Step 2）
- **严禁**以下行为：
  - 生成任何脚本或临时可执行文件（如 `.py`、`.sh`、`.go` 等）；**例外**：允许调用已打包的 `scripts/match_category.py`
  - 执行任何 Go 相关命令（如 `go run`、`go build`、`go test` 等）
  - 通过其他方式间接调用 CLI（如 `bash -c`、管道、包装脚本等）
- 最终报告统一输出至 `resources/report.md`（通过 `render --output` 指定）
- **当且仅当**全流程执行成功时，读取 `resources/report.md` 并向用户**原样返回**其全部内容，**不得**追加任何额外文字或说明

---

## 场景一：关键字检索

### 触发条件

用户需要按关键字查询商品，且能提供以下参数（**缺一不可**）：

| 参数 | CLI 参数 | 说明 |
|------|----------|------|
| 关键字 | `--keyword` | 店铺内搜索词 |

缺少关键字时，先向用户索取，**不得**自行猜测或省略。

### 执行步骤

在技能目录依次执行：

```bash
./scripts/lp_fetcher search --keyword "<关键字>" --output resources/search_result.json
./scripts/lp_fetcher render \
  --input resources/search_result.json \
  --template references/report_template.md \
  --output resources/report.md
```

- **Step 1**：`search` 调用微店 API，将结果写入 `resources/search_result.json`
- **Step 2**：`render` 读取 JSON，生成报告至 `resources/report.md`
- **Step 3**：检查执行是否成功；成功则读取 `resources/report.md` 并**原样返回**其全部内容，**不得**追加任何额外文字或说明

### 示例

```bash
./scripts/lp_fetcher search --keyword "Highway To Hell" --output resources/search_result.json
./scripts/lp_fetcher render \
  --input resources/search_result.json \
  --template references/report_template.md \
  --output resources/report.md
```

---

## 场景二：按分类/风格检索

### 强约束

- **只能**使用 `search_by_category` 子命令执行商品检索
- **严禁**使用 `search` 关键字搜索子命令（即使分类名、风格名、艺人名等看起来像搜索词，也不得调用 `search`）
- 若用户意图无法通过分类匹配完成，**终止技能**并建议用户改用场景一（关键字检索），**不得**在场景二内回退到 `search`

### 触发条件

用户希望按分类、风格、材质、厂牌、艺人专区等店铺分类浏览商品（如「爵士」「摇滚」「全新黑胶」「Taylor Swift」等）。

### 执行步骤

#### Step 1：准备分类数据

检查技能目录是否存在 `resources/category_result.json`：

- **不存在**：执行以下命令获取分类配置：

```bash
./scripts/lp_fetcher category --output resources/category_result.json
```

- **已存在**：跳过此步，直接读取文件

#### Step 2：匹配候选分类

执行分类匹配脚本，根据用户描述的分类/风格意图获取候选列表：

```bash
python3 scripts/match_category.py "<用户意图>" -i resources/category_result.json
```

脚本 stdout 返回 JSON 数组，每项含 `cateId` 与 `cateName`（已归一化：去空格、英文小写）。直接使用该结果整理**候选分类列表**，并向用户提问确认要检索哪一个。

示例提问格式：

> 找到以下候选分类，请选择要检索的一项（回复序号或完整分类名）：
> 1. cooljazz冷爵士（cateId: 141851486）
> 2. 爵士（cateId: 127630967）
> 3. …

若脚本返回空数组 `[]`：**终止技能**，告知用户未找到匹配分类，请重新描述或改用关键字检索。

#### Step 3：用户确认

- 用户选择的分类**必须**属于 Step 2 给出的候选列表（按序号或分类名匹配）
- 若用户输入**不属于**任何候选项：**终止技能**，告知用户未找到匹配分类，请重新描述或改用关键字检索
- 若用户选择有效：记录对应 `cateId`，继续执行

#### Step 4：按分类检索并生成报告

**仅允许**使用 `search_by_category`，**禁止**使用 `search`。

```bash
./scripts/lp_fetcher search_by_category --cateId "<cateId>" --output resources/search_result.json
./scripts/lp_fetcher render \
  --input resources/search_result.json \
  --template references/report_template.md \
  --output resources/report.md
```

- **Step 4a**：`search_by_category` 将结果写入 `resources/search_result.json`（**不得**改用 `search`）
- **Step 4b**：`render` 生成报告至 `resources/report.md`
- **Step 4c**：检查执行是否成功；成功则读取 `resources/report.md` 并**原样返回**其全部内容，**不得**追加任何额外文字或说明

---

## 结果文件

| 文件 | 用途 |
|------|------|
| `resources/category_result.json` | 分类检索中间产物，含一维 `cateList`（`cateId`、`catName` 叶子节点列表） |
| `resources/search_result.json` | 检索中间产物，含 `count`、`status`、`msg`、`materials`；用于判断执行成败 |
| `resources/report.md` | 最终检索结果报告；执行成功时原样返回给用户 |

执行成功时，Agent 读取 `resources/report.md` 并原样返回其内容（Step 2 调用 `match_category.py` 用于分类匹配，不属于最终输出）。

## 异常处理

- `resources/search_result.json` 中 `status == "failed"`：向用户说明 `msg`
  - 常见：`查询结果过多、请精确查询关键词`（结果超过 2 页，需缩小范围或换用更具体的分类）
- `resources/category_result.json` 中 `status == "failed"`：向用户说明 `msg`
- 网络或未知错误：提示用户稍后重试或检查网络连接

## 反馈方式

- **执行成功**：读取 `resources/report.md`，**原样返回**文件全部内容，**不得**追加任何额外文字或说明
- **执行失败**：说明对应 JSON 文件中的 `msg` 及可能原因（见异常处理）
