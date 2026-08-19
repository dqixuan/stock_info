#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""获取同花顺行业/概念板块清单，供 Go 调用：stdout 输出 JSON"""
import json
import os
import sys
import warnings

import pandas as pd
import akshare as ak

# --- 忽略警告与代理清理（Go 调用时也保证不跑代理） ---
warnings.filterwarnings("ignore")
for key in ["HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy",
            "ALL_PROXY", "all_proxy", "NO_PROXY", "no_proxy"]:
    os.environ.pop(key, None)


def find_column(df, candidate_keywords):
    """从 DataFrame 中自动匹配符合关键词的列名"""
    for col in df.columns:
        if any(kw.lower() in str(col).lower() for kw in candidate_keywords):
            return col
    return None


def _fetch_single(fetch_fn, code_kws, name_kws, board_type):
    """拉取一类板块，返回统一结构的 DataFrame 或 None"""
    try:
        df = fetch_fn()
        code_col = find_column(df, code_kws)
        name_col = find_column(df, name_kws)
        if not code_col or not name_col:
            print(f"[WARN] {board_type} 列名匹配失败: {list(df.columns)}", file=sys.stderr)
            return None
        out = df[[code_col, name_col]].copy()
        out.columns = ["板块代码", "板块名称"]
        out["板块类型"] = board_type
        return out
    except Exception as e:
        print(f"[WARN] {board_type} 获取失败: {e}", file=sys.stderr)
        return None


def fetch_ths_boards() -> dict:
    """返回统一 JSON 结构，便于 Go 解析"""
    frames = []

    ind = _fetch_single(
        ak.stock_board_industry_name_ths,
        ["代码", "code"],
        ["行业", "板块", "名称", "name"],
        "同花顺行业",
    )
    if ind is not None:
        frames.append(ind)

    con = _fetch_single(
        ak.stock_board_concept_name_ths,
        ["代码", "code"],
        ["概念", "板块", "名称", "name"],
        "同花顺概念",
    )
    if con is not None:
        frames.append(con)

    if not frames:
        return {
            "success": False,
            "error": "未成功获取到任何板块数据，请升级 akshare: pip install --upgrade akshare",
            "data": [],
        }

    all_boards = pd.concat(frames, ignore_index=True)
    all_boards = all_boards[["板块类型", "板块代码", "板块名称"]]

    records = all_boards.to_dict(orient="records")
    # 统一转成纯字符串/JSON 安全类型
    for r in records:
        for k in r:
            if r[k] is not None and not isinstance(r[k], (str, int, float)):
                r[k] = str(r[k])

    # 可选：落一份 CSV 存档
    csv_file = "ths_stock_boards.csv"
    all_boards.to_csv(csv_file, index=False, encoding="utf-8-sig")
    print(f"[INFO] 已保存 CSV: {os.path.abspath(csv_file)}", file=sys.stderr)

    return {
        "success": True,
        "total": len(records),
        "industry_count": int((all_boards["板块类型"] == "同花顺行业").sum()),
        "concept_count": int((all_boards["板块类型"] == "同花顺概念").sum()),
        "csv": os.path.abspath(csv_file),
        "data": records,
    }


if __name__ == "__main__":
    # 支持可选参数： --type industry|concept|all
    wanted = sys.argv[1] if len(sys.argv) > 1 else "all"
    result = fetch_ths_boards()

    if wanted == "industry":
        result["data"] = [r for r in result["data"] if r["板块类型"] == "同花顺行业"]
    elif wanted == "concept":
        result["data"] = [r for r in result["data"] if r["板块类型"] == "同花顺概念"]
    result["total"] = len(result["data"])

    print(json.dumps(result, ensure_ascii=False))
