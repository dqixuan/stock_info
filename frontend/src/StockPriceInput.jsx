import { useState } from 'react'

function StockPriceInput() {
  const [form, setForm] = useState({
    stockId: '',
    tradeDate: '',
    openPrice: '',
    closePrice: '',
    highPrice: '',
    lowPrice: '',
    volume: '',
    amount: '',
    changePercent: '',
    marginBalance: '',
  })
  const [result, setResult] = useState('')
  const [loading, setLoading] = useState(false)

  const handleChange = (field) => (e) => {
    setForm((prev) => ({ ...prev, [field]: e.target.value }))
  }

  const handleSubmit = async () => {
    if (!form.stockId.trim() || !form.tradeDate.trim()) {
      setResult('请输入股票代码和交易日期')
      return
    }

    const payload = {
      stock_id: form.stockId.trim(),
      trade_date: form.tradeDate.trim(),
      open_price: parseFloat(form.openPrice) || 0,
      close_price: parseFloat(form.closePrice) || 0,
      high_price: parseFloat(form.highPrice) || 0,
      low_price: parseFloat(form.lowPrice) || 0,
      volume: parseInt(form.volume, 10) || 0,
      amount: parseFloat(form.amount) || 0,
      change_percent: parseFloat(form.changePercent) || 0,
      margin_balance: parseFloat(form.marginBalance) || 0,
    }

    setLoading(true)
    try {
      const res = await fetch('/api/stock/price', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      const data = await res.json()
      setResult(JSON.stringify(data, null, 2))
    } catch (err) {
      setResult(`请求失败: ${err.message}`)
    } finally {
      setLoading(false)
    }
  }

  const handleClear = () => {
    setForm({
      stockId: '',
      tradeDate: '',
      openPrice: '',
      closePrice: '',
      highPrice: '',
      lowPrice: '',
      volume: '',
      amount: '',
      changePercent: '',
      marginBalance: '',
    })
    setResult('')
  }

  return (
    <section className="card">
      <h2>输入股票价格</h2>
      <p className="desc">调用 <code>POST /api/stock/price</code>，保存单日股票价格数据</p>

      <div className="form-grid">
        <div className="form-group">
          <label>股票代码</label>
          <input
            type="text"
            placeholder="如: 000001"
            value={form.stockId}
            onChange={handleChange('stockId')}
          />
        </div>
        <div className="form-group">
          <label>交易日期</label>
          <input
            type="date"
            value={form.tradeDate}
            onChange={handleChange('tradeDate')}
          />
        </div>
        <div className="form-group">
          <label>开盘价</label>
          <input
            type="number"
            step="0.01"
            placeholder="0.00"
            value={form.openPrice}
            onChange={handleChange('openPrice')}
          />
        </div>
        <div className="form-group">
          <label>收盘价</label>
          <input
            type="number"
            step="0.01"
            placeholder="0.00"
            value={form.closePrice}
            onChange={handleChange('closePrice')}
          />
        </div>
        <div className="form-group">
          <label>最高价</label>
          <input
            type="number"
            step="0.01"
            placeholder="0.00"
            value={form.highPrice}
            onChange={handleChange('highPrice')}
          />
        </div>
        <div className="form-group">
          <label>最低价</label>
          <input
            type="number"
            step="0.01"
            placeholder="0.00"
            value={form.lowPrice}
            onChange={handleChange('lowPrice')}
          />
        </div>
        <div className="form-group">
          <label>成交量</label>
          <input
            type="number"
            step="1"
            placeholder="0"
            value={form.volume}
            onChange={handleChange('volume')}
          />
        </div>
        <div className="form-group">
          <label>成交额</label>
          <input
            type="number"
            step="0.01"
            placeholder="0.00"
            value={form.amount}
            onChange={handleChange('amount')}
          />
        </div>
        <div className="form-group">
          <label>涨跌幅(%)</label>
          <input
            type="number"
            step="0.0001"
            placeholder="0.0000"
            value={form.changePercent}
            onChange={handleChange('changePercent')}
          />
        </div>
        <div className="form-group">
          <label>融资融券余额</label>
          <input
            type="number"
            step="0.01"
            placeholder="0.00"
            value={form.marginBalance}
            onChange={handleChange('marginBalance')}
          />
        </div>
      </div>

      <div className="row">
        <button onClick={handleSubmit} disabled={loading}>
          {loading ? '提交中...' : '提交'}
        </button>
        <button className="btn-secondary" onClick={handleClear} disabled={loading}>
          清空
        </button>
      </div>

      {result && <pre className="result">{result}</pre>}
    </section>
  )
}

export default StockPriceInput