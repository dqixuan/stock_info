import akshare as ak
import json
import sys
import pandas as pd


def detect_market(symbol: str) -> str:
    """
    根据股票代码判断市场：沪市 (sse) 或深市 (szse)
    """
    if symbol.startswith("6"):
        return "sse"
    elif symbol.startswith(("0", "3")):
        return "szse"
    else:
        return "unknown"


def get_margin_data(symbol: str, start_date: str = None, end_date: str = None) -> dict:
    """
    获取指定股票的融资融券数据
    :param symbol:     股票代码，例如 '600519'
    :param start_date: 开始日期，格式 'YYYYMMDD'，默认最近一条
    :param end_date:   结束日期，格式 'YYYYMMDD'，默认最近一条
    """
    market = detect_market(symbol)

    if market == "unknown":
        return {"error": f"无法识别股票代码 {symbol} 的市场，请检查代码是否正确"}

    try:
        # 根据市场选择对应接口，直接传入日期参数
        if market == "sse":
            df = ak.stock_margin_detail_sse(date=end_date or "")
        else:
            df = ak.stock_margin_detail_szse(date=end_date or "")

        print(f"[DEBUG] 接口返回列名: {df.columns.tolist()}", file=sys.stderr)
        print(f"[DEBUG] 总行数: {len(df)}", file=sys.stderr)

    except Exception as e:
        return {"error": f"接口调用失败: {str(e)}"}

    # --- 过滤目标股票 ---
    code_col = "标的证券代码"
    if code_col not in df.columns:
        return {"error": f"未找到列 '{code_col}'，实际列名: {df.columns.tolist()}"}

    # 统一股票代码格式（去除前缀零等）
    df[code_col] = df[code_col].astype(str).str.strip()
    target = df[df[code_col] == symbol].copy()

    if target.empty:
        return {
            "error": f"未找到股票 {symbol} 的融资融券数据",
            "hint": f"请确认该股票代码属于{'沪市' if market == 'sse' else '深市'}，当前查询日期: {end_date or '最新'}"
        }

    # --- 日期过滤（如果接口不支持日期参数，则手动过滤）---
    date_col = "信用交易日期"
    if date_col in target.columns:
        target[date_col] = target[date_col].astype(str).str.strip()
        if start_date:
            target = target[target[date_col] >= start_date]
        if end_date:
            target = target[target[date_col] <= end_date]

    if target.empty:
        return {"error": f"股票 {symbol} 在 {start_date} ~ {end_date} 范围内无融资融券数据"}

    # --- 序列化 ---
    records = []
    for _, row in target.iterrows():
        record = {}
        for key, value in row.items():
            if isinstance(value, float):
                record[key] = value
            elif isinstance(value, int):
                record[key] = int(value)
            elif hasattr(value, 'strftime'):
                record[key] = value.strftime('%Y-%m-%d')
            else:
                record[key] = str(value)
        records.append(record)

    return {
        "symbol": symbol,
        "market": market,
        "count": len(records),
        "data": records if len(records) > 1 else records[0]
    }


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(json.dumps(
            {"error": "用法: python get_stock_margin_info.py <股票代码> [开始日期] [结束日期]"},
            ensure_ascii=False
        ))
        sys.exit(1)

    symbol_arg     = sys.argv[1]
    start_date_arg = sys.argv[2] if len(sys.argv) > 2 else None
    end_date_arg   = sys.argv[3] if len(sys.argv) > 3 else None

    result = get_margin_data(symbol_arg, start_date_arg, end_date_arg)
    print(json.dumps(result, ensure_ascii=False, indent=2))
