import akshare as ak
import json
import sys


def detect_market(symbol: str) -> str:
    """根据股票代码判断市场：沪市(sse)、深市(szse) 或 北交所(bse)"""
    code = symbol.split(".")[0].strip()
    if code.startswith("6"):
        return "sse"
    elif code.startswith(("0", "3")):
        return "szse"
    elif code.startswith(("8", "4", "92")):
        return "bse"
    else:
        return "unknown"


def get_margin_data(symbol: str, start_date: str = None, end_date: str = None) -> dict:
    market = detect_market(symbol)

    if market == "unknown":
        return {"error": f"无法识别股票代码 {symbol} 的市场，请检查代码是否正确"}

    try:
        if market == "sse":
            df = ak.stock_margin_detail_sse(date=end_date or "")
        elif market == "szse":
            df = ak.stock_margin_detail_szse(date=end_date or "")
        else:  # bse 北交所
            df = ak.stock_margin_detail_bse(date=end_date or "")

        print(f"[DEBUG] 接口返回列名: {df.columns.tolist()}", file=sys.stderr)
        print(f"[DEBUG] 总行数: {len(df)}", file=sys.stderr)
    except Exception as e:
        return {"error": f"接口调用失败: {str(e)}"}

    # 兼容沪深(标的证券代码) 与 北交所(证券代码) 的列名
    code_col = "标的证券代码" if "标的证券代码" in df.columns else "证券代码"
    if code_col not in df.columns:
        return {"error": f"未找到列 '{code_col}'，实际列名: {df.columns.tolist()}"}

    # 归一化：去掉空格和 .BJ/.SH/.SZ 后缀后匹配
    df[code_col] = df[code_col].astype(str).str.strip()
    sym = symbol.split(".")[0].strip()
    target = df[df[code_col].str.split(".").str[0] == sym].copy()

    if target.empty:
        return {
            "error": f"未找到股票 {symbol} 的融资融券数据",
            "hint": f"请确认该股票代码属于{'北交所' if market == 'bse' else ('沪市' if market == 'sse' else '深市')}，当前查询日期: {end_date or '最新'}"
        }

    # 日期过滤（北交所明细无此列，会自动跳过）
    date_col = "信用交易日期"
    if date_col in target.columns:
        target[date_col] = target[date_col].astype(str).str.strip()
        if start_date:
            target = target[target[date_col] >= start_date]
        if end_date:
            target = target[target[date_col] <= end_date]

    if target.empty:
        return {"error": f"股票 {symbol} 在 {start_date} ~ {end_date} 范围内无融资融券数据"}

    # 序列化（保持数字类型，时间转字符串）
    records = []
    for _, row in target.iterrows():
        record = {}
        for key, value in row.items():
            if isinstance(value, (float, int)):
                record[key] = value
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
            ensure_ascii=False))
        sys.exit(1)
    result = get_margin_data(
        sys.argv[1],
        sys.argv[2] if len(sys.argv) > 2 else None,
        sys.argv[3] if len(sys.argv) > 3 else None,
    )
    print(json.dumps(result, ensure_ascii=False, indent=2))
