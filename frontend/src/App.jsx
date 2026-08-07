import { useState } from 'react'

function App() {
  const [name, setName] = useState('')
  const [helloResult, setHelloResult] = useState('')
  const [stockResult, setStockResult] = useState('')
  const [loading, setLoading] = useState(false)

  // 调用 GET /helloworld/{name}
  const handleSayHello = async () => {
    if (!name.trim()) {
      setHelloResult('请输入名称')
      return
    }
    setLoading(true)
    try {
      const res = await fetch(`/helloworld/${encodeURIComponent(name)}`)
      const data = await res.json()
      setHelloResult(data.message || JSON.stringify(data))
    } catch (err) {
      setHelloResult(`请求失败: ${err.message}`)
    } finally {
      setLoading(false)
    }
  }

  // 调用 POST /api/stock
  const handleSaveStock = async () => {
    setLoading(true)
    try {
      const res = await fetch('/api/stock', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: '{}',
      })
      const data = await res.json()
      setStockResult(JSON.stringify(data, null, 2))
    } catch (err) {
      setStockResult(`请求失败: ${err.message}`)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container">
      <h1>股票信息服务</h1>

      <section className="card">
        <h2>1. SayHello 接口</h2>
        <p className="desc">调用 <code>GET /helloworld/{'{name}'}</code></p>
        <div className="row">
          <input
            type="text"
            placeholder="请输入名称"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <button onClick={handleSayHello} disabled={loading}>
            {loading ? '请求中...' : '调用 SayHello'}
          </button>
        </div>
        {helloResult && <pre className="result">{helloResult}</pre>}
      </section>

      <section className="card">
        <h2>2. 初始化股票信息</h2>
        <p className="desc">调用 <code>POST /api/stock</code>，异步拉取并保存 A 股股票数据</p>
        <button onClick={handleSaveStock} disabled={loading}>
          {loading ? '请求中...' : '初始化股票数据'}
        </button>
        {stockResult && <pre className="result">{stockResult}</pre>}
      </section>
    </div>
  )
}

export default App