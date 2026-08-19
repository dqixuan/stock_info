#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""通过新浪行情接口获取A股实时行情（方案一·新浪源）"""
import json
import sys
import requests

session = requests.Session()
session.trust_env = False   # 忽略代理，直连
session.headers.update({
    "User-Agent": (
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "
        "AppleWebKit/537.36 (KHTML, like Gecko) "
        "Chrome/124.0.0.0 Safari/537.36"
    ),
    "Referer": "https://finance.sina.com.cn/",   # 新浪必须带这个，否则 403
})


def get_sina_symbol(symbol):
    """返回新浪前缀代码：sh600519 / sz000001 / bj832566"""
    code = symbol.strip()
    if code.startswith("6"):
        return "sh" + code
    elif code.startswith(("0", "3")):
        return "sz" + code
    elif code.startswith(("4", "8", "9")):
        return "bj" + code
    return None


def get_ashare_price(symbol):
    try:
        sina_symbol = get_sina_symbol(symbol)
        if sina_symbol is None:
            return {"error": f"无效的股票代码格式: {symbol}"}

        url = f"https://hq.sinajs.cn/list={sina_symbol}"
        r = session.get(url, timeout=5)
        r.raise_for_status()

        # 关键：新浪返回 GBK，必须解码
        content = r.content.decode("gbk", errors="replace")

        # 取出引号里的部分
        start = content.find('"')
        end = content.rfind('"')
        if start == -1 or end <= start:
            return {"error": f"接口返回异常: {content[:200]}"}
        fields = content[start + 1:end].split(",")
        if len(fields) < 32 or not fields[0]:
            return {"error": f"未找到股票 {symbol} 的行情（可能代码错误或停牌）"}

        def f(i):
            return fields[i] if i < len(fields) else ""

        current = float(f(3))
        prev_close = float(f(2))
        change_amount = round(current - prev_close, 3)
        change_pct = round((current - prev_close) / prev_close * 100, 3) if prev_close else 0.0

        return {
            "symbol": symbol,
            "name": f(0),
            "latest_price": current,            # 当前价
            "change_amount": change_amount,     # 涨跌额
            "change_percentage": change_pct,    # 涨跌幅%
            "buy_price": float(f(6)),           # 买一
            "sell_price": float(f(7)),          # 卖一
            "last_close": prev_close,           # 昨收
            "open_price": float(f(1)),          # 今开
            "high_price": float(f(4)),          # 最高
            "low_price": float(f(5)),           # 最低
            "volume": float(f(8)),              # 成交量(股)
            "turnover": float(f(9)),            # 成交额(元)
            "timestamp": f"{f(30)} {f(31)}",    # 日期+时间
            "currency": "CNY",
        }
    except requests.exceptions.RequestException as e:
        return {"error": f"网络/接口错误: {e}"}
    except Exception as e:
        return {"error": str(e)}


if __name__ == "__main__":
    if len(sys.argv) > 1:
        print(json.dumps(get_ashare_price(sys.argv[1]), ensure_ascii=False))
    else:
        print(json.dumps({"error": "请提供一个股票代码 (例如: python get_stock_price.py 600519)"}, ensure_ascii=False))
