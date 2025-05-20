# Golang HTTP 服务基础

## net/http 包创建基本 HTTP 服务器

Go 语言的标准库 `net/http` 提供了 HTTP 客户端和服务器的实现。创建一个基本的 HTTP 服务器非常简单：

```go
package main

import (
    "fmt"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func main() {
    // 注册处理函数
    http.HandleFunc("/hello", helloHandler)
    
    // 启动服务器
    fmt.Println("服务器启动在 :8080...")
    http.ListenAndServe(":8080", nil)
}
```

这个例子创建了一个监听在 8080 端口的 HTTP 服务器，对路径 `/hello` 的请求将由 `helloHandler` 函数处理。



## HTTP 请求处理流程

HTTP 请求处理流程如下：

1. 客户端发送 HTTP 请求到服务器
2. 服务器接收请求并将其解析为 `http.Request` 对象
3. 服务器根据请求的 URL 路径查找对应的处理器（Handler）
4. 处理器处理请求并生成响应
5. 服务器将响应发送回客户端

在 Go 中，HTTP 处理器需要实现 `http.Handler` 接口：

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

处理函数（如上例中的 `helloHandler`）是一个便捷的方式，它会被转换为 `http.HandlerFunc` 类型，该类型实现了 `http.Handler` 接口。



## ServeMux 与路由注册

`ServeMux` 是 Go 的 HTTP 请求多路复用器，它将收到的请求根据 URL 路径分发给相应的处理器。

```go
package main

import (
    "fmt"
    "net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "首页")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "关于我们")
}

func main() {
    // 创建一个新的 ServeMux
    mux := http.NewServeMux()
    
    // 注册处理器
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/about", aboutHandler)
    
    // 启动服务器，使用自定义的 ServeMux
    fmt.Println("服务器启动在 :8080...")
    http.ListenAndServe(":8080", mux)
}
```

使用自定义的 `ServeMux` 可以更好地控制路由，避免全局 `http` 包级别的处理器注册可能带来的冲突。

### 路由匹配规则

- 精确匹配：如 `/about` 只匹配 `/about`
- 前缀匹配：如 `/users/` 会匹配 `/users/` 和其子路径，如 `/users/123`
- 最长匹配原则：当多个模式匹配同一 URL 时，选择最长的那个
- 根路径 `/` 是默认处理器，当没有其他匹配时使用



## 处理 HTTP 请求参数与表单

### URL 查询参数

```go
func userHandler(w http.ResponseWriter, r *http.Request) {
    // 获取查询参数
    query := r.URL.Query()
    name := query.Get("name")
    age := query.Get("age")
    
    fmt.Fprintf(w, "用户名: %s, 年龄: %s", name, age)
}
```

### 表单数据

```go
func formHandler(w http.ResponseWriter, r *http.Request) {
    // 解析表单数据
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "解析表单失败", http.StatusBadRequest)
        return
    }
    
    // 获取表单数据
    username := r.FormValue("username")
    password := r.FormValue("password")
    
    fmt.Fprintf(w, "用户名: %s, 密码: %s", username, password)
}
```

### JSON 数据

```go
type User struct {
    Username string `json:"username"`
    Email    string `json:"email"`
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
    // 检查请求方法
    if r.Method != http.MethodPost {
        http.Error(w, "只支持 POST 方法", http.StatusMethodNotAllowed)
        return
    }
    
    // 解析 JSON 数据
    var user User
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&user)
    if err != nil {
        http.Error(w, "解析 JSON 失败", http.StatusBadRequest)
        return
    }
    
    fmt.Fprintf(w, "接收到用户: %s, 邮箱: %s", user.Username, user.Email)
}
```



## 设置 HTTP 响应头与状态码

### 设置状态码

```go
func statusHandler(w http.ResponseWriter, r *http.Request) {
    // 设置状态码
    w.WriteHeader(http.StatusCreated) // 201
    fmt.Fprintf(w, "资源已创建")
}
```

### 设置响应头

```go
func headerHandler(w http.ResponseWriter, r *http.Request) {
    // 设置响应头
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Custom-Header", "custom-value")
    
    // 写入 JSON 响应
    data := map[string]string{"message": "成功"}
    json.NewEncoder(w).Encode(data)
}
```

### 返回不同类型的响应

```go
func responseTypeHandler(w http.ResponseWriter, r *http.Request) {
    contentType := r.URL.Query().Get("type")
    
    switch contentType {
    case "json":
        w.Header().Set("Content-Type", "application/json")
        data := map[string]string{"message": "这是 JSON 响应"}
        json.NewEncoder(w).Encode(data)
    case "html":
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprintf(w, "<html><body><h1>这是 HTML 响应</h1></body></html>")
    default:
        w.Header().Set("Content-Type", "text/plain")
        fmt.Fprintf(w, "这是纯文本响应")
    }
}
```



## 中间件的概念与实现

中间件是一种函数，它可以在 HTTP 请求到达最终处理器之前或之后执行代码。中间件通常用于：

- 日志记录
- 身份验证和授权
- 请求数据的预处理
- 响应数据的后处理
- 错误处理
- 请求计时和性能监控

### 基本中间件示例

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 处理前的代码
        start := time.Now()
        fmt.Printf("开始处理 %s 请求: %s\n", r.Method, r.URL.Path)
        
        // 调用下一个处理器
        next.ServeHTTP(w, r)
        
        // 处理后的代码
        fmt.Printf("完成处理 %s 请求: %s, 耗时: %v\n", r.Method, r.URL.Path, time.Since(start))
    })
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 获取认证令牌
        token := r.Header.Get("Authorization")
        
        // 简单的认证检查
        if token != "valid-token" {
            http.Error(w, "未授权", http.StatusUnauthorized)
            return
        }
        
        // 认证通过，继续处理
        next.ServeHTTP(w, r)
    })
}

func secureHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "这是受保护的内容")
}

func main() {
    // 创建处理器
    handler := http.HandlerFunc(secureHandler)
    
    // 应用中间件（从内到外执行）
    handler = authMiddleware(handler)
    handler = loggingMiddleware(handler)
    
    // 注册处理器
    http.Handle("/secure", handler)
    
    http.ListenAndServe(":8080", nil)
}
```

### 中间件链

更有组织的中间件链实现：

```go
type Middleware func(http.Handler) http.Handler

// 应用多个中间件
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
    for _, middleware := range middlewares {
        h = middleware(h)
    }
    return h
}

func main() {
    // 创建处理器
    handler := http.HandlerFunc(secureHandler)
    
    // 使用中间件链
    chainedHandler := Chain(handler, authMiddleware, loggingMiddleware)
    
    // 注册处理器
    http.Handle("/secure", chainedHandler)
    
    http.ListenAndServe(":8080", nil)
}
```



## 静态文件服务的配置

Go 的 `http` 包提供了 `FileServer` 处理器，可以方便地提供静态文件服务。

### 基本静态文件服务

go

```go
func main() {
    // 创建文件服务处理器，提供 "./static" 目录的内容
    fs := http.FileServer(http.Dir("./static"))
    
    // 注册处理器在 "/static/" 路径下
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    
    fmt.Println("静态文件服务器启动在 :8080...")
    http.ListenAndServe(":8080", nil)
}
```

这段代码会将 `./static` 目录中的文件提供给 `/static/` 路径下的请求。`http.StripPrefix` 函数用于去除 URL 路径中的前缀，使得文件服务器可以正确地定位文件。

### 单文件服务

go

```go
func serveFile(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./static/index.html")
}

func main() {
    http.HandleFunc("/", serveFile)
    http.ListenAndServe(":8080", nil)
}
```

### 自定义文件服务器

go

```go
func customFileServer(dir string, prefix string) http.Handler {
    fs := http.FileServer(http.Dir(dir))
    handler := http.StripPrefix(prefix, fs)
    
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 自定义逻辑，例如日志记录
        fmt.Printf("请求静态文件: %s\n", r.URL.Path)
        
        // 缓存控制
        w.Header().Set("Cache-Control", "max-age=86400") // 24小时
        
        // 调用原始文件服务器
        handler.ServeHTTP(w, r)
    })
}

func main() {
    // 使用自定义文件服务器
    http.Handle("/static/", customFileServer("./static", "/static/"))
    http.ListenAndServe(":8080", nil)
}
```

## 完整示例：综合应用

以下是一个结合上述所有知识点的完整 HTTP 服务器示例：

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

// 中间件类型
type Middleware func(http.Handler) http.Handler

// 中间件链
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
    for _, middleware := range middlewares {
        h = middleware(h)
    }
    return h
}

// 日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("开始 %s %s", r.Method, r.URL.Path)
        
        next.ServeHTTP(w, r)
        
        log.Printf("完成 %s %s (%v)", r.Method, r.URL.Path, time.Since(start))
    })
}

// API 响应格式
type ApiResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// 发送 JSON 响应的工具函数
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
    response, err := json.Marshal(payload)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write(response)
}

// 主页处理器
func homeHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    
    fmt.Fprintf(w, "欢迎访问 Go HTTP 服务示例")
}

// API 处理器
func apiHandler(w http.ResponseWriter, r *http.Request) {
    // 处理查询参数
    query := r.URL.Query()
    name := query.Get("name")
    if name == "" {
        name = "访客"
    }
    
    // 返回 JSON 响应
    response := ApiResponse{
        Status:  "success",
        Message: "API 请求成功",
        Data: map[string]string{
            "name": name,
            "time": time.Now().Format(time.RFC3339),
        },
    }
    
    respondJSON(w, http.StatusOK, response)
}

// 表单处理器
func formHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "只支持 POST 方法", http.StatusMethodNotAllowed)
        return
    }
    
    // 解析表单
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "解析表单失败", http.StatusBadRequest)
        return
    }
    
    // 获取表单数据
    name := r.FormValue("name")
    email := r.FormValue("email")
    
    // 返回响应
    fmt.Fprintf(w, "表单提交成功！姓名: %s, 邮箱: %s", name, email)
}

func main() {
    // 创建路由器
    mux := http.NewServeMux()
    
    // 注册处理器
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/api", apiHandler)
    mux.HandleFunc("/form", formHandler)
    
    // 静态文件服务
    fileServer := http.FileServer(http.Dir("./static"))
    mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
    
    // 使用中间件
    handler := Chain(mux, loggingMiddleware)
    
    // 启动服务器
    addr := ":8080"
    log.Printf("服务器启动在 %s...", addr)
    log.Fatal(http.ListenAndServe(addr, handler))
}
```

在这个完整示例中，我们实现了：

1. 基本的 HTTP 服务器
2. 不同类型的请求处理器
3. 路由注册与分发
4. 请求参数处理
5. JSON 响应
6. 中间件链
7. 静态文件服务

要运行这个服务器，你需要创建一个 `static` 目录在程序的同级目录下，然后在该目录中放置一些静态文件进行测试。

# Web 安全基础

## 输入验证与防止注入攻击

输入验证是防止注入攻击的第一道防线：

- **服务器端验证**：所有输入必须在服务器端进行验证，不能只依赖客户端验证
- **参数化查询**：使用预处理语句和参数化查询而非直接拼接SQL字符串
- **输入净化**：对特殊字符进行转义或过滤
- **白名单验证**：只接受已知安全的输入格式
- **ORM框架**：使用成熟的ORM框架可以减少SQL注入风险

## HTTPS配置与TLS证书处理

安全传输层是网站安全的基础：

- **强制HTTPS**：重定向HTTP请求到HTTPS
- **TLS最低版本**：配置服务器只接受TLS 1.2或更高版本
- **证书管理**：使用可信CA签发的证书，设置自动更新
- **密码套件选择**：启用强密码套件，禁用弱加密算法
- **HSTS**：实现HTTP严格传输安全（HSTS）头部

## 跨站脚本(XSS)防护措施

防止恶意脚本注入执行：

- **内容安全策略(CSP)**：限制可执行脚本来源
- **输出编码**：根据输出上下文正确编码用户输入
- **XSS过滤器**：使用现代框架提供的XSS防护功能
- **HttpOnly标志**：设置cookie的HttpOnly标志防止JavaScript访问
- **DOM-based XSS保护**：安全地处理前端JavaScript中的数据



## 跨站请求伪造(CSRF)防护

防止未授权的操作执行：

- **CSRF令牌**：为表单添加唯一、随机的令牌
- **同源检查**：验证请求来源
- **SameSite Cookie属性**：设置SameSite=Strict或Lax属性
- **请求方法限制**：敏感操作只使用POST而非GET
- **自定义请求头**：添加自定义头部如X-Requested-With

## 安全HTTP头的设置

通过HTTP头部增强安全：

- **Content-Security-Policy**：定义允许加载的资源来源
- **X-Content-Type-Options: nosniff**：防止MIME类型嗅探
- **X-Frame-Options**：控制页面是否可被嵌入frame
- **Referrer-Policy**：控制发送的Referrer信息
- **Permissions-Policy**：限制浏览器功能的使用

## 处理敏感数据与密码存储

保护敏感信息不被泄露：

- **密码哈希**：使用bcrypt、Argon2等专用哈希算法
- **盐值与工作因子**：使用随机盐值并设置适当的工作因子
- **敏感数据加密**：使用强加密算法保护存储数据
- **传输加密**：敏感数据在传输过程中加密
- **最小权限原则**：限制对敏感数据的访问权限

## Rate Limiting的实现与配置

防止滥用和暴力攻击：

- **请求限流**：基于IP、用户ID或API密钥限制请求频率
- **渐进式延迟**：连续失败的尝试增加等待时间
- **验证码**：对可疑活动触发验证码挑战
- **IP黑名单**：暂时或永久封禁恶意IP
- **监控与告警**：设置异常活动监控和告警系统



# RESTful API 开发

## RESTful API 设计原则

REST (表述性状态转移) 设计原则包括：

- **以资源为中心**：API应围绕资源设计，而非操作
- **使用HTTP方法表示操作**：GET获取，POST创建，PUT更新，DELETE删除
- **无状态通信**：服务器不保存客户端状态
- **统一接口**：使用标准化的方法访问资源
- **可缓存性**：响应应明确标识是否可缓存
- **分层系统**：客户端不需要了解后端架构



## 处理不同HTTP方法

- GET：获取资源，幂等，不应修改数据

  ```
  GET /api/users      // 获取用户列表
  GET /api/users/123  // 获取特定用户
  ```

- POST：创建资源，非幂等

  ```
  POST /api/users     // 创建新用户
  ```

- PUT：全量更新资源，幂等

  ```
  PUT /api/users/123  // 更新指定用户的全部信息
  ```

- PATCH：部分更新资源

  ```
  PATCH /api/users/123 // 更新用户的部分信息
  ```

- DELETE：删除资源，幂等

  ```
  DELETE /api/users/123 // 删除指定用户
  ```

## JSON与XML的序列化与反序列化

### JSON

```javascript
// JSON序列化(对象转字符串)
const user = { id: 123, name: "张三" };
const jsonString = JSON.stringify(user);

// JSON反序列化(字符串转对象)
const parsedUser = JSON.parse(jsonString);
```

### XML

```javascript
// 使用库如xml2js进行处理
const xml2js = require('xml2js');

// 序列化
const builder = new xml2js.Builder();
const xml = builder.buildObject({ user: { id: 123, name: "张三" } });

// 反序列化
const parser = new xml2js.Parser();
parser.parseString(xml, (err, result) => {
  const user = result.user;
});
```



## 内容协商与版本控制

### 内容协商

通过HTTP头部实现：

```
Accept: application/json
Accept: application/xml
```

服务器响应：

```
Content-Type: application/json
Content-Type: application/xml
```

### 版本控制方法

- **URL路径**：`/api/v1/users`
- **请求头**：`Accept: application/vnd.company.v1+json`
- **查询参数**：`/api/users?version=1.0`

## API身份验证的基本实现

### JWT（JSON Web Token）

```javascript
const jwt = require('jsonwebtoken');
const secret = 'your-secret-key';

// 生成令牌
function generateToken(user) {
  return jwt.sign({ id: user.id, role: user.role }, secret, { expiresIn: '24h' });
}

// 验证中间件
function authenticateToken(req, res, next) {
  const token = req.headers['authorization']?.split(' ')[1];
  
  if (!token) return res.sendStatus(401);
  
  jwt.verify(token, secret, (err, user) => {
    if (err) return res.sendStatus(403);
    req.user = user;
    next();
  });
}
```

### 基本身份验证

```javascript
function basicAuth(req, res, next) {
  const authHeader = req.headers.authorization;
  
  if (!authHeader || !authHeader.startsWith('Basic ')) {
    return res.status(401).send('Authentication required');
  }
  
  const credentials = Buffer.from(authHeader.split(' ')[1], 'base64')
    .toString()
    .split(':');
    
  const username = credentials[0];
  const password = credentials[1];
  
  // 验证逻辑...
}
```

## API错误处理最佳实践

### 标准错误响应格式

json

```json
{
  "error": {
    "code": "INVALID_PARAMETER",
    "message": "提供的用户ID无效",
    "details": {
      "field": "user_id",
      "value": "abc",
      "reason": "用户ID必须是数字"
    },
    "requestId": "req-123456"
  }
}
```

### HTTP状态码使用

- **2xx**：成功
- **4xx**：客户端错误（如400请求错误，401未授权，404未找到）
- **5xx**：服务器错误（如500内部服务器错误）

### 全局错误处理中间件

```javascript
app.use((err, req, res, next) => {
  console.error(err.stack);
  
  res.status(err.statusCode || 500).json({
    error: {
      code: err.code || 'INTERNAL_ERROR',
      message: err.message || '服务器内部错误',
      details: err.details,
      requestId: req.id
    }
  });
});
```

## API文档生成与维护

### Swagger/OpenAPI

```javascript
// 使用swagger-jsdoc和swagger-ui-express
const swaggerJsDoc = require('swagger-jsdoc');
const swaggerUi = require('swagger-ui-express');

const swaggerOptions = {
  definition: {
    openapi: '3.0.0',
    info: {
      title: '用户API',
      version: '1.0.0',
      description: '用户管理API文档'
    },
    servers: [
      {
        url: 'http://localhost:3000/api'
      }
    ]
  },
  apis: ['./routes/*.js']
};

const swaggerDocs = swaggerJsDoc(swaggerOptions);
app.use('/api-docs', swaggerUi.serve, swaggerUi.setup(swaggerDocs));
```

### API注释示例

```javascript
/**
 * @swagger
 * /users:
 *   get:
 *     summary: 获取所有用户
 *     description: 返回系统中的所有用户列表
 *     responses:
 *       200:
 *         description: 成功返回用户列表
 *         content:
 *           application/json:
 *             schema:
 *               type: array
 *               items:
 *                 type: object
 *                 properties:
 *                   id:
 *                     type: integer
 *                   name:
 *                     type: string
 */
```

