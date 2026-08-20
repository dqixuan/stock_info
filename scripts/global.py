#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Global Market Data Fetcher
==========================
Fetches global stock index data and industry/sector data from Sina Finance
and akshare, outputting structured JSON for consumption by a Go backend.

Output JSON shapes
------------------
--index mode:
    {
        "realtime": [
            {
                "symbol":      str,   // e.g. "s_sh000001"
                "name":        str,   // e.g. "上证指数"
                "last":        float|null,
                "change":      float|null,
                "change_pct":  float|null,
                "volume":      float|null,
                "amount":      float|null,
                "country":     str,   // e.g. "中国"
                "region":      str    // e.g. "亚太"
            }, ...
        ],
        "history": {
            "<symbol>": [
                {
                    "date":       str,   // "YYYY-MM-DD"
                    "open":       float|null,
                    "high":       float|null,
                    "low":        float|null,
                    "close":      float|null,
                    "volume":     float|null
                }, ...
            ], ...
        }
    }

--industry china mode:
    [
        {
            "name":        str,
            "last":        float|null,
            "change_pct":  float|null,
            "change":      float|null,
            "volume":      float|null,
            "amount":      float|null,
            "up_count":    int|null,
            "down_count":  int|null,
            "limit_up":    int|null,
            "lead_stock":  str|null
        }, ...
    ]

--industry global mode:
    [
        {
            "symbol":      str,
            "name":        str,
            "sector":      str,
            "region":      str,
            "last":        float|null,
            "open":        float|null,
            "prev_close":  float|null,
            "high":        float|null,
            "low":         float|null,
            "change":      float|null,
            "change_pct":  float|null,
            "volume":      float|null
        }, ...
    ]

--all mode:
    {
        "index":    { ...same as --index... },
        "china":    [ ...same as --industry china... ],
        "global":   [ ...same as --industry global... ]
    }
"""

import argparse
import json
import math
import os
import sys
import time
import warnings
from typing import Any, Dict, List, Optional, Tuple

import requests

# Suppress pandas / akshare FutureWarnings that pollute stdout
warnings.filterwarnings("ignore")

# ---------------------------------------------------------------------------
# Constants
# ---------------------------------------------------------------------------

USER_AGENT = (
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "
    "AppleWebKit/537.36 (KHTML, like Gecko) "
    "Chrome/124.0.0.0 Safari/537.36"
)
REFERER_SINA = "https://finance.sina.com.cn/"
SINA_HQ_URL = "https://hq.sinajs.cn/list="
DEFAULT_TIMEOUT = 15  # seconds

# ---------------------------------------------------------------------------
# Global sector ETF proxy list
# Each entry: (sina_symbol, display_name, sector, region)
# Sina Finance uses US stock symbols prefixed with "gb_" for US-listed ETFs.
# ---------------------------------------------------------------------------
GLOBAL_SECTOR_ETFS: List[Tuple[str, str, str, str]] = [
    # US Sectors (SPDR ETFs)
    ("gb_xlk",  "科技精选行业ETF",      "科技",     "美国"),
    ("gb_xlf",  "金融精选行业ETF",      "金融",     "美国"),
    ("gb_xle",  "能源精选行业ETF",      "能源",     "美国"),
    ("gb_xlv",  "医疗保健精选行业ETF",  "医疗健康", "美国"),
    ("gb_xlp",  "必需消费品精选行业ETF","必需消费", "美国"),
    ("gb_xly",  "非必需消费品精选行业ETF","可选消费","美国"),
    ("gb_xli",  "工业精选行业ETF",      "工业",     "美国"),
    ("gb_xlb",  "材料精选行业ETF",      "原材料",   "美国"),
    ("gb_xlre", "房地产精选行业ETF",    "房地产",   "美国"),
    ("gb_xlu",  "公用事业精选行业ETF",  "公用事业", "美国"),
    ("gb_xlc",  "通信服务精选行业ETF",  "通信服务", "美国"),
    # Broad market / thematic proxies
    ("gb_soxx", "费城半导体ETF",        "半导体",   "美国"),
    ("gb_arkk", "ARK创新ETF",          "创新科技", "美国"),
    ("gb_gdx",  "黄金矿业ETF",          "黄金矿业", "全球"),
    ("gb_iau",  "iShares黄金ETF",       "贵金属",   "全球"),
    # Asia / EM broad proxies available on Sina
    ("gb_eem",  "新兴市场ETF",          "综合",     "新兴市场"),
    ("gb_fxi",  "中国大盘股ETF",        "综合",     "中国"),
    ("gb_ewj",  "日本ETF",             "综合",     "日本"),
    ("gb_ewh",  "香港ETF",             "综合",     "香港"),
    ("gb_ewg",  "德国ETF",             "综合",     "德国"),
    ("gb_ewu",  "英国ETF",             "综合",     "英国"),
]

# Country / region mapping for global index symbols returned by akshare
# akshare stock_zh_index_global returns a '地区' or '国家' column — we map
# common values to a normalised region string.
REGION_MAP: Dict[str, str] = {
    "美国":   "北美",
    "加拿大": "北美",
    "英国":   "欧洲",
    "德国":   "欧洲",
    "法国":   "欧洲",
    "意大利": "欧洲",
    "西班牙": "欧洲",
    "荷兰":   "欧洲",
    "瑞士":   "欧洲",
    "瑞典":   "欧洲",
    "欧元区": "欧洲",
    "日本":   "亚太",
    "中国":   "亚太",
    "香港":   "亚太",
    "韩国":   "亚太",
    "台湾":   "亚太",
    "澳大利亚":"亚太",
    "新加坡": "亚太",
    "印度":   "亚太",
    "巴西":   "拉美",
    "墨西哥": "拉美",
    "阿根廷": "拉美",
    "南非":   "非洲",
    "俄罗斯": "欧洲/亚洲",
}

# ---------------------------------------------------------------------------
# Session factory — strips proxy env vars per project convention
# ---------------------------------------------------------------------------

def make_session() -> requests.Session:
    """
    Create a requests.Session with trust_env=False so that HTTP_PROXY /
    HTTPS_PROXY environment variables injected by the OS or Go launcher are
    ignored.  The Go side calls stripProxyEnv() before spawning the script,
    but we add a second layer of defence here.
    """
    session = requests.Session()
    session.trust_env = False  # <-- critical: ignore proxy env vars
    session.headers.update(
        {
            "User-Agent": USER_AGENT,
            "Referer": REFERER_SINA,
            "Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
        }
    )
    return session


# Module-level shared session
_SESSION: Optional[requests.Session] = None


def get_session() -> requests.Session:
    global _SESSION
    if _SESSION is None:
        _SESSION = make_session()
    return _SESSION


# ---------------------------------------------------------------------------
# JSON serialisation helpers
# ---------------------------------------------------------------------------

def _safe_float(value: Any) -> Optional[float]:
    """Convert a value to float, returning None for NaN / inf / non-numeric."""
    try:
        f = float(value)
        if math.isnan(f) or math.isinf(f):
            return None
        return f
    except (TypeError, ValueError):
        return None


def _safe_int(value: Any) -> Optional[int]:
    """Convert a value to int, returning None on failure."""
    try:
        f = float(value)
        if math.isnan(f) or math.isinf(f):
            return None
        return int(f)
    except (TypeError, ValueError):
        return None


class _NaNSafeEncoder(json.JSONEncoder):
    """JSON encoder that converts NaN/inf to null rather than raising."""

    def iterencode(self, o, _one_shot=False):
        # Replace Python floats that are NaN/inf before encoding
        return super().iterencode(self._sanitise(o), _one_shot)

    def _sanitise(self, obj):
        if isinstance(obj, float):
            if math.isnan(obj) or math.isinf(obj):
                return None
            return obj
        if isinstance(obj, dict):
            return {k: self._sanitise(v) for k, v in obj.items()}
        if isinstance(obj, (list, tuple)):
            return [self._sanitise(v) for v in obj]
        return obj


def to_json(obj: Any) -> str:
    return json.dumps(obj, ensure_ascii=False, cls=_NaNSafeEncoder, indent=2)


# ---------------------------------------------------------------------------
# DataFrame column finder — fuzzy substring match
# ---------------------------------------------------------------------------

def find_column(df, *candidates: str) -> Optional[str]:
    """
    Find the first column in *df* whose name (lowercased, stripped) contains
    any of the candidate substrings (also lowercased).

    This tolerates akshare renaming columns between releases, e.g.
    '涨跌额' vs '涨跌', '最新价' vs '现价', etc.

    Returns the actual column name or None if nothing matched.

    Example
    -------
    col = find_column(df, "涨跌幅", "change_pct", "pct")
    if col:
        df[col].apply(...)
    """
    cols_lower = {c: c.lower().strip() for c in df.columns}
    for candidate in candidates:
        needle = candidate.lower().strip()
        for actual, lower in cols_lower.items():
            if needle in lower:
                return actual
    return None


# ---------------------------------------------------------------------------
# Sina Finance HQ direct fetcher (GBK-encoded, comma-separated)
# ---------------------------------------------------------------------------

def _sina_hq_fetch(symbols: List[str]) -> Dict[str, Dict[str, Any]]:
    """
    Fetch real-time quotes for a list of Sina symbols via hq.sinajs.cn.

    The response line format for each symbol is:
        var hq_str_<symbol>="name,open,prev_close,last,high,low,bid,ask,
                              volume,amount,date,time,...";

    Returns a dict keyed by symbol with parsed fields.
    Gracefully returns empty dict on any network or parse failure.
    """
    if not symbols:
        return {}

    url = SINA_HQ_URL + ",".join(symbols)
    result: Dict[str, Dict[str, Any]] = {}
    try:
        resp = get_session().get(url, timeout=DEFAULT_TIMEOUT)
        resp.encoding = "gbk"  # Sina HQ is GBK-encoded
        text = resp.text
    except Exception:
        return result

    for line in text.splitlines():
        line = line.strip()
        if not line or "=" not in line:
            continue
        try:
            # var hq_str_gb_xlk="Technology Select...,123.45,..."
            lhs, rhs = line.split("=", 1)
            symbol = lhs.replace("var hq_str_", "").strip()
            data_str = rhs.strip().strip(";").strip('"')
            if not data_str:
                continue
            parts = data_str.split(",")

            # Field positions (Sina standard for US stocks via gb_ prefix):
            # 0  name
            # 1  open
            # 2  prev_close
            # 3  last (current price)
            # 4  high
            # 5  low
            # 6  bid
            # 7  ask
            # 8  volume (shares)
            # 9  amount (total value)
            # 10 date  YYYY-MM-DD
            # 11 time  HH:MM:SS
            def _p(idx: int) -> Optional[float]:
                try:
                    return _safe_float(parts[idx])
                except IndexError:
                    return None

            prev_close = _p(2)
            last = _p(3)

            if prev_close and last is not None and prev_close != 0:
                change = (last - prev_close) if last is not None else None
                change_pct = (
                    (change / prev_close * 100) if change is not None else None
                )
            else:
                change = None
                change_pct = None

            result[symbol] = {
                "name":       parts[0] if parts else None,
                "open":       _p(1),
                "prev_close": prev_close,
                "last":       last,
                "high":       _p(4),
                "low":        _p(5),
                "change":     change,
                "change_pct": change_pct,
                "volume":     _p(8),
                "date":       parts[10] if len(parts) > 10 else None,
                "time":       parts[11] if len(parts) > 11 else None,
            }
        except Exception:
            continue  # skip malformed lines; never crash

    return result


# ---------------------------------------------------------------------------
# Global index — realtime
# ---------------------------------------------------------------------------

def _fetch_global_index_realtime() -> List[Dict[str, Any]]:
    """
    Primary source: akshare stock_zh_index_global()
    Fallback:       Sina HQ direct for a curated list of major index symbols.

    Returns a list of index snapshot dicts.
    """
    rows: List[Dict[str, Any]] = []

    # --- Primary: akshare ---
    try:
        import akshare as ak  # type: ignore

        df = ak.stock_zh_index_global()
        if df is not None and not df.empty:
            sym_col     = find_column(df, "代码", "symbol", "code")
            name_col    = find_column(df, "名称", "name", "指数名")
            last_col    = find_column(df, "最新价", "last", "现价", "收盘")
            change_col  = find_column(df, "涨跌额", "change", "涨跌")
            pct_col     = find_column(df, "涨跌幅", "pct", "change_pct", "涨幅")
            vol_col     = find_column(df, "成交量", "volume", "vol")
            amt_col     = find_column(df, "成交额", "amount", "amt")
            country_col = find_column(df, "国家", "地区", "country", "region")

            for _, row in df.iterrows():
                country = str(row[country_col]).strip() if country_col else ""
                rows.append(
                    {
                        "symbol":     str(row[sym_col]).strip() if sym_col else None,
                        "name":       str(row[name_col]).strip() if name_col else None,
                        "last":       _safe_float(row[last_col]) if last_col else None,
                        "change":     _safe_float(row[change_col]) if change_col else None,
                        "change_pct": _safe_float(row[pct_col]) if pct_col else None,
                        "volume":     _safe_float(row[vol_col]) if vol_col else None,
                        "amount":     _safe_float(row[amt_col]) if amt_col else None,
                        "country":    country,
                        "region":     REGION_MAP.get(country, "其他"),
                    }
                )
            if rows:
                return rows
    except Exception:
        pass  # fall through to Sina direct fallback

    # --- Fallback: Sina HQ direct (major world indices) ---
    # Sina uses s_<code> for domestic indices and specific codes for overseas
    FALLBACK_INDICES: List[Tuple[str, str, str]] = [
        # (sina_symbol, display_name, country)
        ("s_sh000001",  "上证指数",      "中国"),
        ("s_sz399001",  "深证成指",      "中国"),
        ("s_sz399006",  "创业板指",      "中国"),
        ("s_sh000300",  "沪深300",       "中国"),
        ("gb_dji",      "道琼斯工业",    "美国"),
        ("gb_ixic",     "纳斯达克综合",  "美国"),
        ("gb_inx",      "标准普尔500",   "美国"),
        ("gb_n225",     "日经225",       "日本"),
        ("gb_hsi",      "恒生指数",      "香港"),
        ("gb_ftse",     "富时100",       "英国"),
        ("gb_gdaxi",    "德国DAX",       "德国"),
        ("gb_fchi",     "法国CAC40",     "法国"),
        ("gb_stoxx50e", "欧洲斯托克50",  "欧元区"),
        ("gb_twii",     "台湾加权",      "台湾"),
        ("gb_kospi",    "韩国综合",      "韩国"),
        ("gb_axjo",     "澳大利亚ASX200","澳大利亚"),
        ("gb_sensex",   "印度孟买SENSEX","印度"),
        ("gb_bvsp",     "巴西IBOVESPA",  "巴西"),
    ]

    symbols = [s for s, _, _ in FALLBACK_INDICES]
    name_map = {s: (n, c) for s, n, c in FALLBACK_INDICES}
    quotes = _sina_hq_fetch(symbols)

    for sym in symbols:
        name, country = name_map[sym]
        q = quotes.get(sym, {})
        rows.append(
            {
                "symbol":     sym,
                "name":       q.get("name") or name,
                "last":       q.get("last"),
                "change":     q.get("change"),
                "change_pct": q.get("change_pct"),
                "volume":     q.get("volume"),
                "amount":     None,
                "country":    country,
                "region":     REGION_MAP.get(country, "其他"),
            }
        )

    return rows


# ---------------------------------------------------------------------------
# Global index — daily history
# ---------------------------------------------------------------------------

def _fetch_global_index_history(
    symbols: List[str],
    period: str = "daily",
    limit: int = 90,
) -> Dict[str, List[Dict[str, Any]]]:
    """
    Fetch daily OHLCV history for a list of global index symbols via akshare
    stock_zh_index_global_hist().

    *limit* controls how many recent trading days to return (default 90).

    Returns a dict keyed by symbol; missing symbols return an empty list.
    """
    history: Dict[str, List[Dict[str, Any]]] = {}

    try:
        import akshare as ak  # type: ignore
    except ImportError:
        return history

    for sym in symbols:
        try:
            # akshare signature: stock_zh_index_global_hist(symbol, period)
            # period: "daily" | "weekly" | "monthly"
            df = ak.stock_zh_index_global_hist(symbol=sym, period=period)
            if df is None or df.empty:
                history[sym] = []
                continue

            date_col  = find_column(df, "日期", "date", "time", "Date")
            open_col  = find_column(df, "开盘", "open", "Open")
            high_col  = find_column(df, "最高", "high", "High")
            low_col   = find_column(df, "最低", "low", "Low")
            close_col = find_column(df, "收盘", "close", "Close", "最新")
            vol_col   = find_column(df, "成交量", "volume", "Volume", "vol")

            # Take last *limit* rows (most recent)
            df = df.tail(limit)

            bars: List[Dict[str, Any]] = []
            for _, row in df.iterrows():
                date_val = row[date_col] if date_col else None
                try:
                    date_str = str(date_val)[:10]  # truncate to YYYY-MM-DD
                except Exception:
                    date_str = None

                bars.append(
                    {
                        "date":   date_str,
                        "open":   _safe_float(row[open_col]) if open_col else None,
                        "high":   _safe_float(row[high_col]) if high_col else None,
                        "low":    _safe_float(row[low_col]) if low_col else None,
                        "close":  _safe_float(row[close_col]) if close_col else None,
                        "volume": _safe_float(row[vol_col]) if vol_col else None,
                    }
                )
            history[sym] = bars
            # Be polite to the upstream server
            time.sleep(0.2)
        except Exception:
            history[sym] = []

    return history


# ---------------------------------------------------------------------------
# China industry boards (东方财富)
# ---------------------------------------------------------------------------

def _fetch_china_industry() -> List[Dict[str, Any]]:
    """
    Fetch real-time industry board data from 东方财富 via akshare
    stock_board_industry_name_em().

    Column names are resolved at runtime via find_column() to tolerate
    akshare version differences.

    Returns a list of industry board dicts.
    """
    rows: List[Dict[str, Any]] = []

    try:
        import akshare as ak  # type: ignore

        df = ak.stock_board_industry_name_em()
        if df is None or df.empty:
            return rows

        name_col    = find_column(df, "行业名称", "名称", "板块名称", "name")
        last_col    = find_column(df, "最新价", "现价", "last", "收盘", "价格")
        pct_col     = find_column(df, "涨跌幅", "涨幅", "change_pct", "pct")
        change_col  = find_column(df, "涨跌额", "涨跌", "change")
        vol_col     = find_column(df, "成交量", "volume", "vol")
        amt_col     = find_column(df, "成交额", "amount", "amt")
        up_col      = find_column(df, "上涨家数", "上涨数", "涨家数", "up")
        down_col    = find_column(df, "下跌家数", "下跌数", "跌家数", "down")
        limit_col   = find_column(df, "涨停数", "涨停", "limit_up")
        lead_col    = find_column(df, "领涨股票", "领涨", "lead", "龙头")

        for _, row in df.iterrows():
            rows.append(
                {
                    "name":       str(row[name_col]).strip() if name_col else None,
                    "last":       _safe_float(row[last_col]) if last_col else None,
                    "change_pct": _safe_float(row[pct_col]) if pct_col else None,
                    "change":     _safe_float(row[change_col]) if change_col else None,
                    "volume":     _safe_float(row[vol_col]) if vol_col else None,
                    "amount":     _safe_float(row[amt_col]) if amt_col else None,
                    "up_count":   _safe_int(row[up_col]) if up_col else None,
                    "down_count": _safe_int(row[down_col]) if down_col else None,
                    "limit_up":   _safe_int(row[limit_col]) if limit_col else None,
                    "lead_stock": str(row[lead_col]).strip() if lead_col else None,
                }
            )
    except Exception:
        pass  # return whatever we have (may be empty list)

    return rows


# ---------------------------------------------------------------------------
# International sector data (ETF proxy via Sina HQ)
# ---------------------------------------------------------------------------

def _fetch_global_sector() -> List[Dict[str, Any]]:
    """
    Fetch real-time quotes for representative global sector ETFs from
    Sina Finance hq.sinajs.cn.

    Uses the curated GLOBAL_SECTOR_ETFS constant list; each ETF acts as a
    proxy for its sector and region.

    Returns a list of sector dicts.
    """
    symbols = [sym for sym, _, _, _ in GLOBAL_SECTOR_ETFS]
    meta = {sym: (name, sector, region) for sym, name, sector, region in GLOBAL_SECTOR_ETFS}

    quotes = _sina_hq_fetch(symbols)

    rows: List[Dict[str, Any]] = []
    for sym in symbols:
        display_name, sector, region = meta[sym]
        q = quotes.get(sym, {})
        rows.append(
            {
                "symbol":     sym,
                "name":       q.get("name") or display_name,
                "sector":     sector,
                "region":     region,
                "last":       q.get("last"),
                "open":       q.get("open"),
                "prev_close": q.get("prev_close"),
                "high":       q.get("high"),
                "low":        q.get("low"),
                "change":     q.get("change"),
                "change_pct": q.get("change_pct"),
                "volume":     q.get("volume"),
            }
        )

    return rows


# ---------------------------------------------------------------------------
# High-level mode handlers
# ---------------------------------------------------------------------------

def run_index_mode(include_history: bool = True) -> Dict[str, Any]:
    """
    Collect global index realtime snapshots and optional daily history.

    History is fetched only for a representative subset of symbols to
    avoid hammering the upstream API.
    """
    realtime = _fetch_global_index_realtime()

    # Select symbols with non-null last price for history fetching (up to 10)
    history_candidates = [
        r["symbol"]
        for r in realtime
        if r.get("symbol") and r.get("last") is not None
    ][:10]

    history: Dict[str, List[Dict[str, Any]]] = {}
    if include_history and history_candidates:
        history = _fetch_global_index_history(history_candidates)

    return {"realtime": realtime, "history": history}


def run_china_industry_mode() -> List[Dict[str, Any]]:
    """Return China domestic industry board data."""
    return _fetch_china_industry()


def run_global_sector_mode() -> List[Dict[str, Any]]:
    """Return international sector ETF proxy data."""
    return _fetch_global_sector()


def run_all_mode() -> Dict[str, Any]:
    """Collect all data sources and return combined payload."""
    return {
        "index":  run_index_mode(include_history=True),
        "china":  run_china_industry_mode(),
        "global": run_global_sector_mode(),
    }


# ---------------------------------------------------------------------------
# CLI entry point
# ---------------------------------------------------------------------------

def _build_arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        description=(
            "Fetch global stock index and industry/sector data, "
            "output JSON to stdout for Go backend consumption."
        )
    )
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument(
        "--index",
        action="store_true",
        help="Fetch global index realtime snapshots + daily history",
    )
    group.add_argument(
        "--industry",
        choices=["china", "global", "all"],
        metavar="MODE",
        help=(
            "Fetch industry/sector data. "
            "china = A股行业板块; global = 国际板块ETF代理; all = 两者"
        ),
    )
    group.add_argument(
        "--all",
        dest="fetch_all",
        action="store_true",
        help="Fetch all data (index + china industry + global sector)",
    )
    parser.add_argument(
        "--no-history",
        action="store_true",
        default=False,
        help="Skip daily history fetch in --index mode (faster)",
    )
    return parser


def main() -> None:
    # Strip proxy-related env vars as a Python-side safety net
    # (the Go launcher already calls stripProxyEnv, but belt-and-suspenders)
    for var in (
        "HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy",
        "ALL_PROXY", "all_proxy", "NO_PROXY", "no_proxy",
    ):
        os.environ.pop(var, None)

    parser = _build_arg_parser()
    args = parser.parse_args()

    try:
        if args.fetch_all:
            result = run_all_mode()

        elif args.index:
            result = run_index_mode(include_history=not args.no_history)

        elif args.industry == "china":
            result = run_china_industry_mode()

        elif args.industry == "global":
            result = run_global_sector_mode()

        elif args.industry == "all":
            result = {
                "china":  run_china_industry_mode(),
                "global": run_global_sector_mode(),
            }

        else:
            parser.print_help(sys.stderr)
            sys.exit(1)

        # Output JSON to stdout — Go reads this via os.Exec / NextResult
        print(to_json(result))

    except KeyboardInterrupt:
        sys.exit(0)
    except Exception as exc:
        # Never crash silently: emit a JSON error object so Go can detect it
        error_payload = {
            "error": True,
            "message": str(exc),
            "data": None,
        }
        print(to_json(error_payload))
        sys.exit(1)


if __name__ == "__main__":
    main()
