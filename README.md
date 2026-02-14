# HomeNet - Go Server/Client 项目

## 项目说明

这是一个使用 Go 实现的 Server/Client 应用，功能如下：

- **Server 端**: 监听 `127.0.0.1:8080`，接收来自 Client 的请求，并返回请求方的 IP 地址
- **Client 端**: 每隔 1 秒向 Server 发送一次请求，获取并打印自己的 IP 地址

## 项目结构

```
homenet/
├── go.mod                # Go 模块文件
├── server/
│   └── main.go          # Server 程序
└── client/
    └── main.go          # Client 程序
```

## 运行方式

### 1. 启动 Server

在终端中运行：
```bash
cd /Users/shauntso/homenet
go run ./server/main.go
```

输出示例：
```
Server 启动在 :8080
来自 127.0.0.1 的请求
来自 127.0.0.1 的请求
```

### 2. 启动 Client（新的终端窗口）

在另一个终端中运行：
```bash
cd /Users/shauntso/homenet
go run ./client/main.go
```

输出示例：
```
Client 启动，每隔 1 秒向 http://127.0.0.1:8080 发送请求
[10:30:25] 我的 IP 地址是: 127.0.0.1
[10:30:26] 我的 IP 地址是: 127.0.0.1
[10:30:27] 我的 IP 地址是: 127.0.0.1
```

## 修改 Server 地址

如果需要修改 Server 的地址（例如连接到不同的 IP），编辑 [client/main.go](client/main.go)，修改以下行：

```go
const serverURL = "http://127.0.0.1:8080"
```

## 修改 Server 端口

如果需要修改 Server 的监听端口，编辑 [server/main.go](server/main.go)，修改以下行：

```go
const port = ":8080"
```

## 编译

如果需要编译为可执行文件：

```bash
# 编译 Server
go build -o server ./server/main.go

# 编译 Client
go build -o client ./client/main.go

# 运行编译后的程序
./server    # 终端 1
./client    # 终端 2
```
