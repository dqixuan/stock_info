import requests
import json
import sys


def get_ashare_price(symbol):
    """东财 push2 单股接口：一次只查一只，避免全市场拉取被限流"""
    try:
        code = symbol.strip()
        # 确定 secid 前缀：1=沪, 0=深/北
        if code.startswith("6"):
            secid = "1." + code
        elif code.startswith(("0", "3", "4", "8", "9")):
            secid = "0." + code
        else:
            return {"error": f"无效的股票代码格式: {symbol}"}

        url = "https://push2.eastmoney.com/api/qt/stock/get"
        params = {
            "secid": secid,
            "fltt": "2",
            "invt": "2",
            "fields": "f43,f44,f45,f46,f47,f48,f50,f51,f52,f57,f58,f60,f86,f170",
        }
        headers = {
            "User-Agent": (
                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "
                "AppleWebKit/537.36 (KHTML, like Gecko) "
                "Chrome/124.0.0.0 Safari/537.36"
            ),
            "Referer": "https://quote.eastmoney.com/",
        }
        r = session.get(url, params=params, timeout=5, headers=headers)

#         r = requests.get(url, params=params, timeout=5)
        r.raise_for_status()
        d = r.json().get("data")
        if not d:
            return {"error": f"未能找到股票代码: {symbol}"}

        return {
            "symbol": code,
            "market": secid,
            "name": d.get("f58"),
            "latest_price": d.get("f43"),        # 最新价
            "change_percentage": d.get("f170"),  # 涨跌幅(%)
            "open_price": d.get("f46"),          # 今开
            "high_price": d.get("f44"),          # 最高
            "low_price": d.get("f45"),           # 最低
            "last_close": d.get("f60"),          # 昨收
            "buy_price": d.get("f19"),           # 买一价
            "sell_price": d.get("f20"),          # 卖一价
            "volume": d.get("f47"),              # 成交量
            "turnover": d.get("f48"),            # 成交额
            "timestamp": d.get("f86"),
            "currency": "CNY",
        }
    except requests.exceptions.RequestException as e:
        return {"error": f"网络/接口错误: {e}"}
    except Exception as e:
        return {"error": str(e)}


if __name__ == "__main__":
    if len(sys.argv) > 1:
        result = get_ashare_price(sys.argv[1])
        print(json.dumps(result, ensure_ascii=False, indent=2))
    else:
        print(json.dumps({"error": "请提供一个股票代码 (例如: python get_ashare_price.py 600519)"},
                         ensure_ascii=False))
