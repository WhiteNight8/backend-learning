## struct

```go
type person struct {
	name string
	age int
}
```

### struct的匿名字段

最外层优先访问

```go
package main

import "fmt"

type Human struct {
    name string
    age int
    phone string
}

type Employee struct {
    Human
    speciality string
    phone string
}

func main() {
    Bob := Employee{Human { "Bob",34,"777-444-XXX"}, "Designer", "333-222"}
    fmt.Println("Bob's work phone is:,",Bob.phone)
}
```



### 面向对象

method

 a methd is a function with an implicit first argument, called a receiver

```go
func (r ReceiverType) funcName(parameters) (results)
```

- metho名字一样，但是接收者不一样，那么method就不一样
- method里面可以访问接收者字段
- 调用method通过.访问

指针作为receiver

如果一个method的receiver是*T，可以在一个T类型的实例变量V上面调用这个method，而不需要&V去调用这个method

method继承

method重写



如此美妙的go的面对对象



### 处理表单的输入

### 验证表单的输入

### 预防跨站脚本

### 防止多次递交表单

### 处理文件上传

