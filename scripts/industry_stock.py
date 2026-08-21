import akshare as ak
import pandas as pd

def get_ths_concept_kline(concept_name="创新药", target_date="2026-08-20"):
    """
    通过 AkShare 获取同花顺 App 概念板块历史数据
    """
    df_ths = None

    # 尝试不同版本的同花顺概念板块 API 名称
    try:
        # 新版本 / 标准接口：获取同花顺概念板块指数行情
        if hasattr(ak, "stock_board_concept_hist_ths"):
            df_ths = ak.stock_board_concept_hist_ths(start_date="20260101", end_date="20261231", symbol=concept_name)
        elif hasattr(ak, "stock_board_concept_index_ths"):
            df_ths = ak.stock_board_concept_index_ths(symbol=concept_name, start_date="20260101", end_date="20261231")
        else:
            # 兜底方案：升到最新版或改用东方财富概念板块（东财与同花顺概念涨跌幅极度接近）
            print("当前 akshare 版本较低，建议执行: pip install --upgrade akshare")
            df_ths = ak.stock_board_concept_hist_em(symbol=concept_name, start_date="20260101", end_date="20261231")

    except Exception as e:
        print(f"AKShare 接口调用失败: {e}")
        return None

    if df_ths is not None and not df_ths.empty:
        # 统一字段重命名
        column_map = {
            "日期": "日期",
            "开盘": "开盘", "开盘价": "开盘",
            "收盘": "收盘", "收盘价": "收盘",
            "最高": "最高", "最高价": "最高",
            "最低": "最低", "最低价": "最低",
            "成交量": "成交量"
        }
        df_ths.rename(columns=column_map, inplace=True)

        # 统一格式化日期为 YYYY-MM-DD
        df_ths["日期"] = df_ths["日期"].astype(str)

        day_df = df_ths[df_ths["日期"] == target_date]
        if not day_df.empty:
            return day_df[["日期", "开盘", "收盘", "最高", "最低", "成交量"]]
        else:
            latest_date = df_ths["日期"].max()
            print(f"未找到【{target_date}】的数据，最新可用交易日为：{latest_date}")
            return df_ths.tail(1)[["日期", "开盘", "收盘", "最高", "最低", "成交量"]]

    return None

if __name__ == "__main__":
    target_date = "2026-08-20"
    df = get_ths_concept_kline(concept_name="创新药", target_date=target_date)

    if df is not None:
        print("\n【行情获取成功】:")
        print(df)