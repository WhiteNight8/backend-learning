### 定义变量

```go
// define single varible
var varibelName type

// define mutiple varible
var vname1, vname2, vname3 type

//initialize
var variableName type = value
var vname1,vname2,vname3 type = v1, v2, v3

var vanme1,vname2,vname3 = v1, v2,v3
vname,vname2,vname3 := v1, v2,v3

//简洁声明只能用于函数内部，var定义全局变量
```



### 常量

```go
const Pi= 3.1415926
// 常量，也就是在程序编译阶段就确定下来的值，程序运行时无法改变
```

### 内置基础类型

- Boolean
- 数值类型： int unit
- 字符串： 字符串不可变，修改的话，需要改为[]byte类型，再转为string
- 错误类型



### Go数据底层的存储

![image-20250331215401355](https://raw.githubusercontent.com/JoeyXXia/MyPictureData/main/image-20250331215401355.png)

### 一些技巧

- 分组声明
- iota枚举

### Go程序设计的规则

- 大写字母开头的变量是可导出的，是公有变量，小写字母开头的是不可导出的，是私有变量
- 大写字母开头的函数也是



### array，slice，map

长度也是数组类型的一部分

```go
a := [3]int{1,2,3}

b := [10]int{1,2,3}

c := [...]int{4,5,6}

```

数组之间的赋值时值的赋值，即当把一个数组作为参数传入函数时，传入的其实时该数组的副本，而不是它的指针

slice并不是真正意义上的动态数组，而是一个引用

slice在声明数组时，方括号内无任何字符

slice是引用类型，当引用修改了其中元素的值时，其他所有引用都会改变该值

- len
- cap
- append
- copy



### map

字典python

```go
var numbers map[string]int

numbers := make(map[string]int)
numbers['one'] = 1
numbers['ten'] = 10
```

- map是无序的
- map的疮毒是不固定的，引用类型
- len返回key的数量
- map的值方便修改
- 不是thread-safe，在多个go-routine存取是，需要使用mutex lock机制

delete删除map的元素

### make，new操作

make用于内建类型的内存分配，new 用于各种类型的内存分配

new返回指针



### 零值

零值，变量未填充前的默认值，通常为0





## session和数据存储

### session和cookie

### Go如何使用session

### session存储

### 预防session劫持





## 如何设计一个web框架

### 项目规划

### 自定义路由器设计

### controller设计

### 日志和配置设计

## 

## 扩展web框架

### 静态文件支持

### session支持

### 表单及其验证支持

### 用户认证

### 多语言支持

### pprof支持



























