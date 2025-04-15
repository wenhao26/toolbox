# 项目说明
> 文件结构
```text
httpx/
├── client.go         // HTTP客户端配置与构造
├── request.go        // 请求构建器（With系列方法）
├── executor.go       // 请求执行器（重试、上下文、日志）
├── parser.go         // 响应解析封装
├── backoff.go        // 重试退避策略
├── middleware.go     // 请求中间件封装
├── trace.go          // Trace ID 工具
└── example_test.go   // 示例测试代码
```
> 为什么要分多个文件
* client.go：封装 HTTP 客户端配置，支持注入全局行为（如超时、重试策略等）
* request.go：请求构建器，支持构造 GET、POST 请求、参数注入等
* executor.go：请求执行器，真正发起请求，内含重试机制、错误日志记录等
* parser.go：响应解析器，将响应体转换为结构体、字符串或字节
* middleware.go：中间件定义，支持注入多个中间件处理请求头、认证、加密等
* trace.go：TraceID 生成工具，默认生成唯一请求链路 ID
* backoff.go：封装退避策略，控制重试时间间隔逻辑
* example_test.go：用于演示与测试，帮助新手快速掌握用法

# 特性说明
* ✅ 模块化：不同功能职责独立，便于阅读和修改
* ✅ 中间件支持：灵活注入通用逻辑，如认证、日志、追踪等
* ✅ 重试机制：高可用关键特性，生产稳定性保障
* ✅ 动态解析：可支持结构化解析（JSON）或原始解析
* ✅ 易扩展：能无缝集成到任何项目，不影响其他组件
* ✅ 资源安全：无泄漏、正确释放连接、健壮性高

> 模块化的好处
* 易测试：每个模块都可以独立单元测试；
* 易扩展：增加某一功能，只需在对应模块中加逻辑；
* 可重用：哪怕你以后换成 gRPC、WebSocket，也能保留中间件链逻辑；
* 符合开闭原则：对扩展开放，对修改封闭。

# 开发逻辑与流程
```text
请求构造（request.go）
   ↓
执行请求（executor.go）
   ↓
触发中间件（middleware.go）
   ↓
发起 HTTP 请求（client.go）
   ↓
自动重试（backoff.go）
   ↓
读取响应 & Trace 日志记录
   ↓
数据解析（parser.go）
```

# 核心组件设计思想详解
1、 Client 是「总控器」
* 它封装了超时、重试策略、是否启用 Trace 等。
* 所有请求都通过它来构造（client.Get(...)、client.Post(...)）

2、Request 是「构造者」
* 提供链式 API（如 .WithHeader()、.WithJSON()）；
* 每个请求都可以独立挂载自己的中间件和参数。

3、Executor 是「执行器」
* 集成日志、重试、trace、连接关闭、安全读取等逻辑；
* 是稳定性保证核心，避免资源泄漏和错误未处理。

4、Parser 是「结果分析器」
* 用户不需要操作 io.ReadAll()，直接调用 .ParseJSON() 即可；
* 可选结构体解析、纯文本、[]byte 响应体。

5、中间件机制
* 类似于 Gin、Echo 的思想；
* 你可以注入认证 token、trace-id、请求签名等操作，不影响主流程。

# 测试用例使用
```shell
go test -v 
go test -v -run {用例函数名称} 
go test -cover
```