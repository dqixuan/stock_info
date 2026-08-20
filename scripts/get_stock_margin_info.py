#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import json
import sys
import os

import pandas as pd
import akshare as ak

for k in ["HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy",
          "ALL_PROXY", "all_proxy", "NO_PROXY", "no_proxy"]:
    os.environ.pop(k, None)


def detect_market(symbol: str) -> str:
    code = symbol.strip().split(".")[0]
    if code.startswith("6"):
        return "sse"
    elif code.startswith(("0", "3")):
        return "szse"
    elif code.startswith(("8", "4", "92")):
        return "bse"
    return "unknown"


def _pick_interface(market: str):
    candidates = {
        "sse": ["stock_margin_detail_sse"],
        "szse": ["stock_margin_detail_szse"],
        "bse": ["stock_margin_detail_bse", "stock_margin_detail_bj"],
    }[market]
    for name in candidates:
        if hasattr(ak, name):
            return getattr(ak, name)
    return None


def normalize_code_series(series: pd.Series) -> pd.Series:
    """规范化证券代码列：永远先生成字符串，再拆分/补零"""
    s = series.fillna("").astype(str).str.strip()      # 1) 先转 str
    s = s.str.replace(r"\s+", "", regex=True)          # 去空白
    s = s.str.split(".").str[0]                        # 去掉 .BJ/.SH 后缀
    s = s.str.zfill(6)                                # 补 6 位，仅对非空
    return s


def get_margin_data(symbol: str, trade_date: str = None) -> dict:
    market = detect_market(symbol)
    if market == "unknown":
        return {"error": f"无法识别股票代码 {symbol} 的市场"}

    fetch = _pick_interface(market)
    if fetch is None:
        return {"error": "akshare 版本过旧，缺少融资融券明细接口，请升级 akshare"}

    try:
        df = fetch(date=trade_date or "")
        print(f"[DEBUG] 列名: {df.columns.tolist()}", file=sys.stderr)
    except Exception as e:
        return {"error": f"接口调用失败: {e}"}

    code_col = "标的证券代码" if "标的证券代码" in df.columns else "证券代码"
    if code_col not in df.columns:
        return {"error": f"未找到代码列，实际列名: {df.columns.tolist()}"}

    # 关键：先转成字符串列，再做归一化
    df["__normalized_code__"] = normalize_code_series(df[code_col])

    sym = symbol.strip().split(".")[0].zfill(6)
    target = df[df["__normalized_code__"] == sym].copy()

    if target.empty:
        return {"error": f"未找到证券 {symbol} 的融资融券数据", "market": market}

    # 日期过滤：全员转为字符串后比较，NaN 兜底
    date_col = "信用交易日期"
    if date_col in target.columns and trade_date:
        target = target[target[date_col].fillna("").astype(str).str.strip() == trade_date]

    records = []
    for _, row in target.iterrows():
        rec = {}
        for k, v in row.items():
            if k == "__normalized_code__":
                continue
            if isinstance(v, (int, float)):
                rec[k] = round(v, 2) if isinstance(v, float) else v
            elif hasattr(v, "strftime"):
                rec[k] = v.strftime("%Y-%m-%d")
            else:
                rec[k] = str(v)
        records.append(rec)

    return {
        "symbol": symbol,
        "market": market,
        "count": len(records),
        "data": records[0] if len(records) == 1 else records,
    }


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(json.dumps({"usage": "python get_margin_info.py <代码> [日期YYYYMMDD]"},
                         ensure_ascii=False))
        sys.exit(1)
    print(json.dumps(get_margin_data(sys.argv[1], sys.argv[2] if len(sys.argv) > 2 else None),
                     ensure_ascii=False, indent=2))
