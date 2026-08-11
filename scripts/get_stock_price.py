import akshare as ak
import json
import sys

def get_ashare_price(symbol):
    try:
        # 添加交易所前缀
        if symbol.startswith("6"):
            symbol = "sh" + symbol
        elif symbol.startswith("0") or symbol.startswith("3"):
            symbol = "sz" + symbol
        else:
            return {"error": "无效的股票代码格式"}

        # 使用 akshare 获取 A 股实时行情数据
        stock_zh_a_spot_df = ak.stock_zh_a_spot()

        # 打印出 DataFrame 的头部，查看列名
        target_stock = stock_zh_a_spot_df[stock_zh_a_spot_df['代码'] == symbol]

        if not target_stock.empty:
            latest_price = target_stock['最新价'].iloc[0]  # 获取最新价
            stock_name = target_stock['名称'].iloc[0]  # 获取股票名称
            change_amount = target_stock['涨跌额'].iloc[0]  # 获取涨跌额
            change_percentage = target_stock['涨跌幅'].iloc[0]  # 获取涨跌幅
            buy_price = target_stock['买入'].iloc[0]  # 获取买入价
            sell_price = target_stock['卖出'].iloc[0]  # 获取卖出价
            last_close = target_stock['昨收'].iloc[0]  # 获取昨收价
            open_price = target_stock['今开'].iloc[0]  # 获取今开价
            high_price = target_stock['最高'].iloc[0]  # 获取最高价
            low_price = target_stock['最低'].iloc[0]  # 获取最低价
            volume = target_stock['成交量'].iloc[0]  # 获取成交量
            turnover = target_stock['成交额'].iloc[0]  # 获取成交额
            timestamp = target_stock['时间戳'].iloc[0]  # 获取时间戳

            data = {
                "symbol": symbol,
                "name": stock_name,
                "latest_price": float(latest_price),
                "change_amount": change_amount,
                "change_percentage": change_percentage,
                "buy_price": buy_price,
                "sell_price": sell_price,
                "last_close": last_close,
                "open_price": open_price,
                "high_price": high_price,
                "low_price": low_price,
                "volume": volume,
                "turnover": turnover,
                "timestamp": timestamp,
                "currency": "CNY"
            }
            print(json.dumps(data, ensure_ascii=False))
        else:
            print(json.dumps({"error": f"未能找到股票代码: {symbol}"}, ensure_ascii=False))

    except Exception as e:
        print(json.dumps({"error": str(e)}, ensure_ascii=False))

if __name__ == "__main__":
    if len(sys.argv) > 1:
        stock_symbol = sys.argv[1]
        get_ashare_price(stock_symbol)
    else:
        print(json.dumps({"error": "请提供一个股票代码 (例如: 600519)"}, ensure_ascii=False))
