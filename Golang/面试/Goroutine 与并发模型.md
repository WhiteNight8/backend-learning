# Goroutine 与并发模型

##  goroutine 的概念与使用方法

Goroutine 是 Go 语言中的轻量级线程，由 Go 运行时(runtime)管理。使用方法非常简单，只需在函数调用前添加 `go` 关键字：

```go
func hello() {
    fmt.Println("Hello, world!")
}

func main() {
    go hello() // 启动一个goroutine执行hello函数
    
    // 使用匿名函数启动goroutine
    go func() {
        fmt.Println("Hello from anonymous function")
    }()
    
    time.Sleep(1 * time.Second) // 等待goroutine执行完毕
}
```



## goroutine 与线程的区别

Goroutine 与传统操作系统线程有以下区别：

- **资源占用**：goroutine 起始栈大小仅有 2KB，可按需增长；而线程通常有较大的固定栈大小(如 2MB)
- **调度方式**：goroutine 由 Go 运行时调度，实现了协作式用户态调度；而线程由操作系统调度
- **创建和销毁成本**：goroutine 创建和销毁开销极小，可轻松创建上万个；线程创建和销毁成本高昂
- **上下文切换**：goroutine 上下文切换成本远低于线程
- **通信方式**：goroutine 推荐使用 channel 通信；线程间通常使用共享内存和锁



## 启动 goroutine 的注意事项与最佳实践

- **确保主程序不会过早退出**：主函数退出时所有 goroutine 会被强制终止
- **使用 sync.WaitGroup 等待所有 goroutine 完成**
- **使用 channel 进行通信和同步**
- **避免过多 goroutine**：虽然 goroutine 很轻量，但创建过多会增加调度开销
- **处理 panic**：单个 goroutine 的 panic 会导致整个程序崩溃，应在 goroutine 中使用 recover

```go
func main() {
    var wg sync.WaitGroup
    
    for i := 0; i < 5; i++ {
        wg.Add(1)
        i := i // 创建变量副本
        
        go func() {
            defer wg.Done()
            defer func() {
                if r := recover(); r != nil {
                    fmt.Println("Recovered:", r)
                }
            }()
            
            fmt.Println("Working on:", i)
        }()
    }
    
    wg.Wait() // 等待所有goroutine完成
}
```





## 匿名函数在 goroutine 中的常见陷阱

最常见的陷阱是循环变量捕获：

```go
// 错误示例
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i) // 大多数情况下会打印 5 个 5
    }()
}

// 正确示例
for i := 0; i < 5; i++ {
    i := i // 创建变量副本
    go func() {
        fmt.Println(i) // 打印 0, 1, 2, 3, 4
    }()
}

// 或者通过参数传递
for i := 0; i < 5; i++ {
    go func(i int) {
        fmt.Println(i) // 打印 0, 1, 2, 3, 4
    }(i)
}
```



##  GOMAXPROCS 的设置与影响

GOMAXPROCS 控制可同时执行 Go 代码的操作系统线程数量：

- Go 1.5 后默认设为 CPU 核心数
- 可通过 `runtime.GOMAXPROCS(n)` 设置
- 也可通过环境变量 `GOMAXPROCS` 设置

影响：

- 过小会限制并行性能
- 过大会增加线程切换开销
- 对 I/O 密集型任务影响较小，对 CPU 密集型任务影响较大



##  goroutine 泄漏的常见原因与预防

泄漏原因：

- **阻塞的 channel 操作**：无缓冲 channel 的发送或接收操作在没有对应接收者或发送者时会永久阻塞
- **无限循环**：goroutine 内的无限循环不退出
- **等待永远不会结束的同步操作**：如等待已满的 WaitGroup
- **互斥锁忘记解锁**

预防措施：

- 为 channel 操作设置超时（使用 `select` 和 `time.After`）
- 提供取消机制（使用 context）
- 使用 defer 确保资源释放
- 定期检查和监控 goroutine 数量

```go
// 使用context提供取消机制
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Cancelled")
            return
        default:
            // 执行工作
            time.Sleep(100 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    go worker(ctx)
    
    time.Sleep(3 * time.Second)
    fmt.Println("Main exiting")
}
```



## 并发与并行的概念区分

- 并发(Concurrency)：
  - 结构性概念，表示程序设计能够同时处理多个任务
  - 不一定同时执行，而是有能力处理多个任务
  - 强调任务的结构和组织方式
  - Go 的设计理念："不要通过共享内存来通信，而是通过通信来共享内存"
- 并行(Parallelism)：
  - 执行性概念，表示多个任务真正同时执行
  - 需要多核心或多处理器硬件支持
  - 强调任务的执行方式
  - 通过 GOMAXPROCS 设置来影响并行度

Go 语言提供并发编程模型，而并行执行则基于具体硬件和 GOMAXPROCS 设置



# Channel 与通信

##  channel 的基本概念与声明方式

Channel 是 Go 语言中用于 goroutine 之间通信的管道，它遵循 CSP (Communicating Sequential Processes) 并发模型。Channel 提供了一种机制，使得数据可以在不同 goroutine 间安全地传递。

声明方式:

```go
// 声明一个传递int类型数据的channel
var ch1 chan int

// 使用make函数初始化无缓冲channel
ch2 := make(chan string)

// 使用make函数初始化带缓冲的channel，缓冲区大小为10
ch3 := make(chan float64, 10)

// 声明一个只读channel
var ch4 <-chan bool

// 声明一个只写channel
var ch5 chan<- []byte
```

## 缓冲通道与无缓冲通道的区别

**无缓冲通道 (Unbuffered Channel)**:

- 使用 `make(chan T)` 创建
- 同步特性: 发送操作会阻塞，直到有接收者接收数据
- 接收操作会阻塞，直到有发送者发送数据
- 用于保证同步，发送和接收是一个"握手"过程

**缓冲通道 (Buffered Channel)**:

- 使用 `make(chan T, capacity)` 创建
- 异步特性: 只有当缓冲区满时，发送操作才会阻塞
- 只有当缓冲区空时，接收操作才会阻塞
- 可以用作临时队列，缓解生产和消费速率不一致的问题

```go
// 无缓冲通道示例
ch := make(chan int)
go func() {
    ch <- 42 // 会阻塞，直到主goroutine接收数据
}()
value := <-ch // 两个goroutine同步点

// 缓冲通道示例
bufCh := make(chan int, 2)
bufCh <- 1 // 不会阻塞
bufCh <- 2 // 不会阻塞
// bufCh <- 3 // 会阻塞，因为缓冲区已满
```



## channel 的发送与接收操作

**发送操作**:

```go
ch <- value // 将value发送到channel ch
```

**接收操作**:

```go
value := <-ch       // 从channel接收数据并赋值给value
value, ok := <-ch   // ok为true表示接收成功，false表示channel已关闭
<-ch                // 接收数据但丢弃，仅用于同步
```

**阻塞特性**:

- 向已满的缓冲通道发送数据会阻塞
- 从空的通道接收数据会阻塞
- 向已关闭的通道发送数据会触发panic
- 从已关闭的通道接收数据会立即返回零值，ok为false



## channel 的关闭与遍历方法

**关闭通道**:

```go
close(ch) // 关闭通道
```

关闭通道的重要规则:

- 只有发送方应该关闭通道，接收方不应该关闭通道
- 通道关闭后不能再发送数据，但可以继续接收已有数据
- 重复关闭通道会触发panic

**遍历通道**:

```go
// 方法1: 使用for循环和range关键字
for value := range ch {
    fmt.Println(value) // 自动处理通道关闭
}

// 方法2: 使用for循环和接收操作
for {
    value, ok := <-ch
    if !ok {
        break // 通道已关闭
    }
    fmt.Println(value)
}
```



## select 语句的使用场景与超时处理

select 语句类似于 switch，但专用于通道操作，可以同时等待多个通道操作。

**基本用法**:

```go
select {
case v1 := <-ch1:
    fmt.Println("Received from ch1:", v1)
case ch2 <- v2:
    fmt.Println("Sent to ch2")
case <-ch3:
    fmt.Println("Received from ch3, but value discarded")
default:
    fmt.Println("No channel operations ready")
}
```

**主要使用场景**:

1. 非阻塞通道操作
2. 多通道等待
3. 超时处理
4. 取消操作

**超时处理示例**:

```go
select {
case data := <-ch:
    processData(data)
case <-time.After(2 * time.Second):
    fmt.Println("Operation timed out")
}
```

**上下文取消示例**:

```go
select {
case data := <-dataCh:
    processData(data)
case <-ctx.Done():
    fmt.Println("Operation cancelled")
    return
}
```



## 单向 channel 的应用

单向通道限制了通道的操作方向，提高了类型安全性和代码可读性。

**声明方式**:

```go
var sendCh chan<- int  // 只发送通道
var recvCh <-chan int  // 只接收通道
```

**转换规则**:

- 双向通道可以转换为单向通道，但单向通道不能转换为双向通道
- 单向通道的方向不能改变

**实际应用**:

```go
func producer(out chan<- int) {
    // 只能向out发送数据，不能从out接收数据
    for i := 0; i < 5; i++ {
        out <- i
    }
    close(out)
}

func consumer(in <-chan int) {
    // 只能从in接收数据，不能向in发送数据
    for num := range in {
        fmt.Println(num)
    }
}

func main() {
    ch := make(chan int)
    go producer(ch)  // 传递双向通道，自动转换为单向发送通道
    consumer(ch)     // 传递双向通道，自动转换为单向接收通道
}
```



## 使用 channel 实现常见的并发模式

### Fan-out 模式（分发任务）

```go
func fanOut(tasks []Task, workers int) {
    taskCh := make(chan Task, len(tasks))
    
    // 分发任务
    for _, task := range tasks {
        taskCh <- task
    }
    close(taskCh)
    
    // 启动workers
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for task := range taskCh {
                process(task)
            }
        }()
    }
    
    wg.Wait()
}
```

### Fan-in 模式（合并结果）

```go
func fanIn(channels ...<-chan int) <-chan int {
    result := make(chan int)
    var wg sync.WaitGroup
    
    // 为每个输入通道启动一个goroutine
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for n := range c {
                result <- n
            }
        }(ch)
    }
    
    // 当所有输入通道都关闭时，关闭结果通道
    go func() {
        wg.Wait()
        close(result)
    }()
    
    return result
}
```

### Pipeline 模式（数据流水线）

```go
func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

// 使用管道
func main() {
    // 构建流水线: 生成数字 -> 计算平方
    numbers := generator(1, 2, 3, 4, 5)
    squares := square(numbers)
    
    // 消费结果
    for sq := range squares {
        fmt.Println(sq)
    }
}
```

### Worker Pool 模式（工作池）

```go
func workerPool(numWorkers int, tasks <-chan Task, results chan<- Result) {
    var wg sync.WaitGroup
    
    // 启动固定数量的worker
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for task := range tasks {
                result := processTask(task)
                results <- result
            }
        }(i)
    }
    
    // 等待所有worker完成并关闭结果通道
    go func() {
        wg.Wait()
        close(results)
    }()
}
```

### 带取消的并发操作

```go
func processWithCancellation(ctx context.Context, data []int) <-chan int {
    results := make(chan int)
    go func() {
        defer close(results)
        for _, v := range data {
            select {
            case <-ctx.Done():
                fmt.Println("Operation cancelled")
                return
            case results <- process(v):
                // 继续处理
            }
        }
    }()
    return results
}

// 使用方式
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel() // 确保取消函数被调用
    
    data := []int{1, 2, 3, 4, 5}
    resultCh := processWithCancellation(ctx, data)
    
    for result := range resultCh {
        fmt.Println(result)
    }
}
```



# 同步与互斥

##  sync.WaitGroup 的使用方法

`sync.WaitGroup` 用于等待一组 goroutine 完成执行，工作原理类似于计数器：



```go
func main() {
    var wg sync.WaitGroup
    
    // 启动5个goroutine
    for i := 1; i <= 5; i++ {
        wg.Add(1) // 增加计数器
        
        i := i // 创建变量副本，避免闭包陷阱
        go func() {
            defer wg.Done() // 完成时减少计数器
            
            fmt.Printf("Worker %d starting\n", i)
            time.Sleep(time.Second)
            fmt.Printf("Worker %d done\n", i)
        }()
    }
    
    wg.Wait() // 阻塞直到计数器变为0
    fmt.Println("All workers done")
}
```

关键方法：

- `Add(delta int)`: 增加计数器，通常在启动 goroutine 之前调用
- `Done()`: 减少计数器，通常通过 defer 在 goroutine 结束时调用
- `Wait()`: 阻塞直到计数器变为 0

注意事项：

- `Add()` 应在 goroutine 外调用，避免竞态条件
- 确保每次 `Add()` 都有对应的 `Done()`
- 计数器不能为负，否则会触发 panic



## Mutex 与 RWMutex 的区别与使用场景

### Mutex (互斥锁)

互斥锁用于保护共享资源，同一时刻只允许一个 goroutine 访问资源：

```go
var (
    mu      sync.Mutex
    balance int
)

func deposit(amount int) {
    mu.Lock()
    defer mu.Unlock()
    balance += amount
}

func withdraw(amount int) bool {
    mu.Lock()
    defer mu.Unlock()
    if balance < amount {
        return false
    }
    balance -= amount
    return true
}
```

### RWMutex (读写互斥锁)

读写锁允许多个读操作并发，但写操作是互斥的：

```go
var (
    rwMu    sync.RWMutex
    data    map[string]string = make(map[string]string)
)

// 写操作需要写锁
func store(key, value string) {
    rwMu.Lock()
    defer rwMu.Unlock()
    data[key] = value
}

// 读操作只需要读锁
func lookup(key string) (string, bool) {
    rwMu.RLock()
    defer rwMu.RUnlock()
    value, ok := data[key]
    return value, ok
}
```

区别与使用场景：

- Mutex：适用于读写频率相近或写操作较多的场景
- RWMutex：适用于读操作远多于写操作的场景
- RWMutex 比 Mutex 有更多开销，对于简单操作可能不值得



## sync.Once 的应用：单例模式实现

`sync.Once` 保证函数仅执行一次，常用于单例模式、初始化操作等：

```go
type Database struct {
    connection string
}

var (
    instance *Database
    once     sync.Once
)

func GetDatabase() *Database {
    once.Do(func() {
        fmt.Println("Creating database connection...")
        instance = &Database{connection: "connected"}
    })
    return instance
}

func main() {
    // 多次调用，初始化只发生一次
    db1 := GetDatabase()
    db2 := GetDatabase()
    
    fmt.Println(db1 == db2) // true，是同一个实例
}
```

特点：

- 即使在并发环境下也能保证函数只执行一次
- 支持延迟初始化（懒加载）
- 线程安全、简洁高效
- 无需担心双重检查锁定的问题

适用场景：

- 单例模式实现
- 配置加载
- 日志系统初始化
- 资源池创建



## sync.Pool 的基本使用与性能优化

`sync.Pool` 用于存储和复用临时对象，减少 GC 压力：

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func processRequest(data []byte) {
    // 从池中获取对象
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset() // 重置缓冲区
        bufferPool.Put(buf) // 使用后放回池中
    }()
    
    // 使用缓冲区
    buf.Write(data)
    // 处理数据...
}
```

关键特性：

- 临时对象池，允许跨 goroutine 复用
- 任何时候池中的对象可能被 GC 回收
- 无法指定最大容量，不能控制池中对象数量
- 底层使用线程本地缓存减少竞争

性能优化：

- 最适合用于大小相似的临时对象（如 buffers）
- 使用前重置对象状态，避免状态泄露
- 调整池大小与对象生命周期以优化性能
- 并非所有对象都适合池化，简单小对象可能得不偿失



## sync.Map 的特点与适用场景

`sync.Map` 是 Go 1.9 引入的并发安全的 map 实现：

```go
var userCache sync.Map

func getUser(id int) (*User, error) {
    // 尝试从缓存获取
    if value, ok := userCache.Load(id); ok {
        return value.(*User), nil
    }
    
    // 缓存未命中，从数据库获取
    user, err := fetchUserFromDB(id)
    if err != nil {
        return nil, err
    }
    
    // 存入缓存
    userCache.Store(id, user)
    return user, nil
}

// 删除缓存
func invalidateUser(id int) {
    userCache.Delete(id)
}

// 遍历所有用户
func processAllUsers(process func(*User)) {
    userCache.Range(func(key, value interface{}) bool {
        process(value.(*User))
        return true // 继续遍历
    })
}
```

特点：

- 读取不需要锁
- 针对特定访问模式优化（读多写少）
- 空间消耗比加锁 map 高
- 不保证遍历顺序
- 性能随着元素数量增加而降低

适用场景：

- 多读少写场景
- 键值缓存
- 只增不减的数据集
- 读多写少且有大量并发访问的场景

不适用场景：

- 大量写操作
- 频繁更新同一个键
- 需要确定大小的场景
- 需要有序遍历的场景



## 原子操作 (atomic 包) 的基本用法

原子操作提供了低级别的同步原语，常用于简单的计数器和标志位：

```go
import (
    "fmt"
    "sync"
    "sync/atomic"
)

func main() {
    // 原子计数器
    var counter int64 = 0
    var wg sync.WaitGroup
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            atomic.AddInt64(&counter, 1)
            wg.Done()
        }()
    }
    
    wg.Wait()
    fmt.Println("Counter:", counter)
    
    // 原子存储和加载
    var flag atomic.Bool
    flag.Store(true)
    
    if flag.Load() {
        fmt.Println("Flag is set")
    }
    
    // 比较并交换 (CAS)
    var value int32 = 100
    swapped := atomic.CompareAndSwapInt32(&value, 100, 200)
    fmt.Println("Swapped:", swapped, "Value:", value)
}
```

常用原子操作：

- `Add`：原子加法
- `Load`：原子读取
- `Store`：原子存储
- `Swap`：原子交换
- `CompareAndSwap`：比较并交换 (CAS)

Go 1.19 后新增类型：

- `atomic.Int64`
- `atomic.Uint64`
- `atomic.Pointer[T]`
- `atomic.Bool`

使用场景：

- 简单计数器
- 标志位
- 无锁数据结构
- 自旋锁实现



## 免并发陷阱的最佳实践

### 1. 避免数据竞争

```go
// 错误示例
var counter int
func increment() {
    counter++ // 数据竞争！
}

// 正确做法
var counter int
var mu sync.Mutex
func increment() {
    mu.Lock()
    counter++
    mu.Unlock()
}

// 或使用原子操作
var counter int64
func increment() {
    atomic.AddInt64(&counter, 1)
}
```

### 2. 正确使用 defer 解锁

```go
// 推荐做法
func doSomething() {
    mu.Lock()
    defer mu.Unlock() // 确保解锁，即使发生panic
    
    // 处理逻辑...
}
```

### 3. 避免嵌套锁，防止死锁

```go
// 可能导致死锁
func transferMoney(from, to *Account, amount int) {
    from.mu.Lock()
    defer from.mu.Unlock()
    
    to.mu.Lock()         // 可能死锁！
    defer to.mu.Unlock()
    
    // 转账逻辑...
}

// 正确做法：按固定顺序获取锁
func transferMoney(from, to *Account, amount int) {
    // 确保按ID顺序获取锁
    if from.id < to.id {
        from.mu.Lock()
        defer from.mu.Unlock()
        to.mu.Lock()
        defer to.mu.Unlock()
    } else {
        to.mu.Lock()
        defer to.mu.Unlock()
        from.mu.Lock()
        defer from.mu.Unlock()
    }
    
    // 转账逻辑...
}
```

### 4. 正确使用 channel

```go
// 避免向已关闭的channel发送数据
close(ch)
ch <- data // 会触发panic！

// 正确用法：由发送方关闭
func producer(ch chan<- int, done <-chan struct{}) {
    defer close(ch) // 生产者负责关闭
    
    for i := 0; ; i++ {
        select {
        case ch <- i:
            // 成功发送
        case <-done:
            return // 收到结束信号
        }
    }
}
```

### 5. 合理使用 context 取消操作

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            // 清理并退出
            return
        default:
            // 正常工作
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel() // 始终调用cancel
    
    go worker(ctx)
    // ...
}
```

### 6. 避免 goroutine 泄漏

```go
// 可能导致泄漏
func processRequest() {
    ch := make(chan Result)
    go func() {
        result := doWork()
        ch <- result // 如果没人接收，会永远阻塞
    }()
    
    // 其他逻辑，可能提前返回而没有接收结果
}

// 正确做法
func processRequest(ctx context.Context) error {
    ch := make(chan Result, 1) // 加缓冲或使用context
    
    go func() {
        result := doWork()
        select {
        case ch <- result:
        case <-ctx.Done():
            // 上下文已取消，避免阻塞
        }
    }()
    
    select {
    case result := <-ch:
        return processResult(result)
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(2 * time.Second):
        return errors.New("timeout")
    }
}
```

### 7. 正确处理共享变量

```go
// 闭包陷阱
for i := 0; i < 10; i++ {
    go func() {
        fmt.Println(i) // 大多数情况下会打印10个 "10"
    }()
}

// 正确做法
for i := 0; i < 10; i++ {
    i := i // 创建局部副本
    go func() {
        fmt.Println(i)
    }()
    
    // 或通过参数传递
    go func(val int) {
        fmt.Println(val)
    }(i)
}
```

### 8. 使用静态分析工具

- 使用 `go vet` 检测常见错误
- 使用 `-race` 标志检测数据竞争
- 考虑使用第三方静态分析工具如 `golangci-lint`

```bash
go vet ./...
go test -race ./...
go build -race
```

### 9. 避免过度使用锁

```go
// 避免锁粒度过大
func processItems(items []Item) {
    var mu sync.Mutex
    
    for _, item := range items {
        mu.Lock()
        process(item) // 可能是耗时操作
        mu.Unlock()
    }
}

// 优化锁粒度
func processItems(items []Item) {
    results := make([]Result, len(items))
    
    var wg sync.WaitGroup
    for i, item := range items {
        wg.Add(1)
        i, item := i, item // 本地副本
        go func() {
            defer wg.Done()
            results[i] = process(item) // 并行处理
        }()
    }
    wg.Wait()
    
    // 最后合并结果...
}
```

### 10. 优先使用通信而非共享内存

```go
// Go 并发哲学：不要通过共享内存来通信，而是通过通信来共享内存

// 使用通道代替锁
func worker(tasks <-chan Task, results chan<- Result) {
    for task := range tasks {
        results <- process(task)
    }
}

func main() {
    tasks := make(chan Task, 100)
    results := make(chan Result, 100)
    
    // 启动工作池
    for i := 0; i < numWorkers; i++ {
        go worker(tasks, results)
    }
    
    // 分发任务和收集结果
    // ...
}
```

