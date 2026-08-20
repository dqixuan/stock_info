import akshare as ak
import pandas as pd

# 1. 获取所有行业板块列表 (东方财富)
industry_list = ak.stock_board_industry_name_em()
print(industry_list[['板块名称', '板块代码']])

# 2. 获取特定行业（例如：“半导体”）的成分股
# symbol 参数传入具体的行业名称
cons_df = ak.stock_board_industry_cons_em(symbol="半导体")

# 提取股票代码和名称
stocks = cons_df[['代码', '名称']]
print(stocks.head())