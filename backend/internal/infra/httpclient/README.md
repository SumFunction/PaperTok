# HTTP Client Infrastructure

> HTTP 客户端基础设施，封装 HTTP 请求能力

---

## 职责

- 提供 HTTP 客户端接口
- 支持超时配置
- 便于测试 Mock

---

## 接口

```go
type HTTPClient interface {
    Get(ctx context.Context, url string) (*http.Response, error)
    Do(req *http.Request) (*http.Response, error)
}
```

---

## 文件结构

| 文件 | 说明 |
|------|------|
| `client.go` | HTTP 客户端实现 |

---

## 使用示例

```go
client := httpclient.NewClient(httpclient.Config{
    Timeout: 10 * time.Second,
})

// GET 请求
resp, err := client.Get(ctx, "http://example.com/api")

// 读取响应
body, err := httpclient.ReadBody(resp)
```

---

## 测试 Mock

```go
type mockHTTPClient struct {
    response *http.Response
    err      error
}

func (m *mockHTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
    return m.response, m.err
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    return m.response, m.err
}
```
