## interface

什么是interface

interface是一组method签名的组合，通过interface来定义对象的一组行为

```go
type Human struct {
    name string
    age int
    phone string
}

type Student struct {
    Human
    school string
    loan float32
}

type Employee struct {
    Human
    company string
    money float 32
}

func (h *Human) SayHi() {
    fmt.Printf("Hi, I am %s you can call me %s\n", h,name, h.home)
}

func (h *Human) Sing(lyrics string) {
    fmt.Println("la la la...",lyrics)
}

func (h *Human) Guzzle(beeStein string) {
    fmt.Printlb("guzzle guzzle...", beeStein)
}

func (e *Employee) SayHi() {
    fmt.Printf("hi  I am %s, I work at %s, call me on %s", e.name, e.company, e.phone)
}

func (s *Student) BorrowMoney(amount float32) {
    s.loan += amount
}

func (e *Employee) SpendSalary(amout float32) {
    e.money -= amount
}

type Men interface {
    SayHi()
    Sing(lyrics string)
    Guzzle(beeStein string)
}

type YoungChap interface {
    SayHi()
    Sing(song string)
    BorrowMoeny(amount float32)
}

type ElderLyGent interface {
    SayHi()
    Sing(song string)
    SpendSalary(amount float32)
}
```

interface可以别任意的对象实现，任意的类型都实现了空interface



**interface值**

interface就是一组抽象方法的集合，必须由其他非iinterface类型实现，而不能自我实现

**空interface**

一个函数把interface{} 作为参数，那么就可以接受任意类型的值作为参数，如果返回interface{}，也可以返回任意类型

**interface函数参数**

任何实现String方法的类型都能作为参数被fmt.Println调用

**interface变量存储的类型**

- comma-ok断言
- switch测试



**嵌入interface**

内置的逻辑语法

**反射**

能够检查程序在运行时的状态



## 并发

并行设计

## **goroutine**

并发编程

**channels**

通信机制

**Buffered Challels**

缓存设置

**range和close**

生产者关闭channel

**select**

监听channel上的数据流动

**超时**

select设置

**runtime goroutine**

- Goexit
- Gosched
- NumCPU
- NumGoroutine
- GOMAXPROCS



## 访问数据库

### database/sql接口

### 使用mysql数据库

### 使用sqlite数据库

### 使用postgreSQL数据库

### 使用Beego orm库进行ORM开发

### NOSQL数据库操作





