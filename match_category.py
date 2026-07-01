#!/usr/bin/env python3
import argparse
import json
import sys


def normalize(text: str) -> str:
    return "".join(c.lower() if "A" <= c <= "Z" else c for c in text if not c.isspace())


def main() -> int:
    parser = argparse.ArgumentParser(description="根据关键词匹配分类 cateId")
    parser.add_argument("query", help="分类匹配关键词")
    parser.add_argument(
        "-i",
        "--input",
        default="category_result.json",
        help="分类数据 JSON 路径（默认: category_result.json）",
    )
    args = parser.parse_args()

    query = normalize(args.query)
    if not query:
        print("query 不能为空", file=sys.stderr)
        return 1

    try:
        with open(args.input, encoding="utf-8") as f:
            data = json.load(f)
    except FileNotFoundError:
        print(f"文件不存在: {args.input}", file=sys.stderr)
        return 1
    except json.JSONDecodeError as e:
        print(f"JSON 解析失败: {e}", file=sys.stderr)
        return 1

    if data.get("status") == "failed":
        print(data.get("msg") or "分类数据获取失败", file=sys.stderr)
        return 1

    matches = []
    for item in data.get("cateList") or []:
        cat_name = item.get("catName", "")
        if query in cat_name:
            matches.append({
                "cateId": item["cateId"],
                "cateName": cat_name,
            })

    print(json.dumps(matches, ensure_ascii=False))
    return 0


if __name__ == "__main__":
    sys.exit(main())
