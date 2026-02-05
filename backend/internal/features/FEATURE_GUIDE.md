# Feature 开发指南

本指南定义 Feature 模块的开发规范。

---

## 1. Feature 定位

Feature 是按业务功能垂直划分的独立模块，每个 Feature 包含处理该功能所需的完整逻辑。

| 特性 | 说明 |
|------|------|
| **高内聚** | 功能代码集中在一个目录 |
| **独立测试** | 可独立单元测试 |
| **依赖向下** | 只依赖 Core、Repository、Infra |

---

## 2. 目录结构

```
features/featurename/
├── README.md          # 模块文档（推荐）
├── interface.go       # 对外接口定义（必须）
├── deps.go            # 依赖接口定义（必须）
├── types.go           # 领域类型（可选）
├── errors.go          # 错误定义（推荐）
├── service.go         # 实现（必须）
└── service_test.go    # 单元测试（必须）
```

---

## 3. 文件规范

### interface.go

```go
// Service defines the interface for xxx operations.
type Service interface {
    // MethodName does something.
    // @Params:
    //   - ctx: context
    //   - param: description
    // @Returns:
    //   - result: description
    //   - error: when fails
    MethodName(ctx context.Context, param Type) (Result, error)
}
```

### deps.go

```go
// dependency defines the xxx capability required.
type dependency interface {
    Method(ctx context.Context) error
}
```

### service.go

```go
type Impl struct {
    dep dependency
}

var _ Service = (*Impl)(nil) // compile-time check

func New(dep dependency) *Impl {
    return &Impl{dep: dep}
}

func (s *Impl) MethodName(ctx context.Context, param Type) (Result, error) {
    // implementation
}
```

### service_test.go

```go
func TestImpl_MethodName(t *testing.T) {
    // Arrange
    mockDep := &mockDependency{}
    svc := New(mockDep)

    // Act
    result, err := svc.MethodName(ctx, param)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

---

## 4. 依赖约束

| 可以依赖 | 不能依赖 |
|---------|---------|
| Core Services | 其他 Features |
| Repository | API Layer |
| Infrastructure | Facade |

---

## 5. 开发流程

1. 创建目录 `internal/features/newfeature/`
2. 定义接口 `interface.go`
3. 定义依赖 `deps.go`
4. 编写测试 `service_test.go`（TDD）
5. 实现逻辑 `service.go`
6. 运行测试确认通过
7. 在 Facade 中注册

---

## 6. 现有 Features

| Feature | 职责 |
|---------|------|
| `paperfeed` | 论文推荐流 |
| `papersearch` | 论文搜索 |
