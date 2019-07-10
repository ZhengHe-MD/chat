# 用 Go 语言搭建聊天服务器 — 从协议设计到服务实现

本文介绍使用 Go 语言搭建聊天服务，从协议设计到服务实现的过程，并提供源码，你也许可以从这个过程体会到：

* 应用层通信协议的设计和实现
* 服务器运作的基本模型和实现
* 写代码的乐趣

本项目源码地址：https://github.com/ZhengHe-MD/chat

```sh
$ git clone https://github.com/ZhengHe-MD/chat.git
```

本文阅读说明：

* 约 5-10min 阅读时间
* 每段代码引用第一行中注明了出处

## 功能

在搭建真正的聊天服务器之前，我们明确这个聊天服务器需要支持的功能：

* 用户登录、登出
* 用户 A 向用户 B 发送消息
* 建立群聊天
* 用户离开群聊天
* 用户在群里广播消息

其它功能暂时不在本聊天服务考虑范围之内。

## 协议

为保证数据完整、有序地传输，服务在传输层采用 tcp 协议。在应用层，服务采用一个自定义的 (有自主知识产权的)、简陋的、基本能够满足上述功能的协议，CHAT 协议，它的基本格式如下所示：

```sh
[ProtocolName]/[ProtocolVersion] [Command] [Body]\n
```

因为是 CHAT 的第一个正式版本，我们给它一个版本号 1.0。CHAT/1.0 支持的命令如下所示：

```sh
# 用户登录、登出
CHAT/1.0 LOGIN userName\n
CHAT/1.0 LOGOUT\n
# 用户 A 发送消息，用户 B 接收消息
CHAT/1.0 SEND userName data\n
CHAT/1.0 RECEIVE userName data\n
# 建立群聊天、离开群聊天、群广播消息
CHAT/1.0 GROUP groupName [userName]...\n
CHAT/1.0 LEAVE groupName\n
CHAT/1.0 BROADCAST groupName data\n
```

### 协议实现

先定义一些常量：

```go
// chat/protocol/command.go
const (
	ProtocolName    = "CHAT"
	ProtocolVersion = "1.0"
	ProtocolSep     = " "

	CmdSend      = "SEND"
	CmdBroadCast = "BROADCAST"
	CmdLogin     = "LOGIN"
	CmdLogout    = "LOGOUT"
	CmdReceive   = "RECEIVE"
	CmdGroup     = "GROUP"
	CmdLeave     = "LEAVE"
)
```

概念上，我们可以称每条消息为一条命令 (command)，接着我们将定义这些命令结构体及其字符串方法：

```go
// tcp/chat/protocol/command.go
type BaseCommand struct {
  Protocol string
  Version  string
}

func (c *BaseCommand) String() string {
	return fmt.Sprintf("%s/%s", c.Protocol, c.Version)
}

type SendCommand struct {
	BaseCommand
	Name string
	Data []byte
}

func (c *SendCommand) String() string {
	return strings.Join([]string{
		c.BaseCommand.String(),
		CmdSend,
		c.Name,
		string(c.Data),
	}, ProtocolSep) + "\n"
}

// ...check other commands in the source code
```

紧接着，我们需要给出这些命令的 reader 和 writer

* reader：负责利用 io.Reader 接口读取 CHAT 消息并反序列化成相应的命令结构体
* writer：负责利用 io.Writer 接口将命令序列化成相应的 CHAT 消息

writer 比较简单：

```go
// chat/protocol/writer.go
type CommandWriter struct {
  writer *bufio.Writer
}

func NewCommandWriter(writer io.Writer) *CommandWriter {
	return &CommandWriter{
		writer: bufio.NewWriter(writer),
	}
}

func (w *CommandWriter) Write(cmd interface{}) (err error) {
	_, err = w.writer.WriteString(fmt.Sprintf("%v", cmd))
	return w.writer.Flush()
}
```

reader 逻辑比较复杂，大致的步骤就是：

* 从 io.Reader 接口中读取一行
* 解析协议的名称和版本是否符合定义
* 解析命令名称，根据名称来决定如何最后解析成哪个命令结构体，利用 type switch 实现
* 过程中发生任何错误返回 InvalidMessageErr/UnsupportedCmdErr，一切顺利则返回命令本身

举例可以到 tcp/chat/protocol/writer.go 中查到

**至此，我们完成了协议部分的实现**

### 协议实现测试

我不是 TDD 的拥趸，比较习惯先写代码后写测试，必要的测试覆盖能够让我们在重构时充满信心。

reader 测试举例：

```go
// chat/protocol/reader_test.go
func TestSendMessage(t *testing.T) {
	cases := []struct {
		message      string
		expectedErr  error
		expectedName string
		expectedData []byte
	}{
		{
			"CHAT/1.0 SEND zhenghe hello world\n",
			nil,
			"zhenghe",
			[]byte("hello world"),
		},
	}

	for i, c := range cases {
		mr := NewCommandReader(strings.NewReader(c.message))

		cmd, err := mr.Read()
		if err != c.expectedErr {
			t.Errorf("case %d: should have err:%v got:%v",
				i, c.expectedErr, err)
		}

		if err == nil {
			sendCmd := cmd.(*SendCommand)
			if sendCmd.Name != c.expectedName {
				t.Errorf("case %d: should have name:%s got:%s",
					i, c.expectedName, sendCmd.Name)
			}

			if !reflect.DeepEqual(sendCmd.Data, c.expectedData) {
				t.Errorf("case %d: should have data:%v got:%v",
					i, c.expectedData, sendCmd.Data)
			}
		}
	}
}
```

writer 测试举例：

```go
func TestWriteSendMessage(t *testing.T) {
	cases := []struct {
		cmd             *SendCommand
		expectedMessage string
	}{
		{
			&SendCommand{
				BaseCommand: BaseCommand{ProtocolName, ProtocolVersion},
				Name:        "zhenghe",
				Data:        []byte("hello world"),
			},
			"CHAT/1.0 SEND zhenghe hello world\n",
		},
	}

	for i, c := range cases {
		buf := bytes.NewBuffer([]byte{})
		mw := NewCommandWriter(buf)

		_ = mw.Write(c.cmd)

		if buf.String() != c.expectedMessage {
			t.Errorf("Case %d: expect message:%s got:%s",
				i, c.expectedMessage, buf.String())
		}
	}
}
```

测试写完后，执行一下：

```sh
$ go test -v
=== RUN   TestInvalidMessage
--- PASS: TestInvalidMessage (0.00s)
=== RUN   TestSendMessage
--- PASS: TestSendMessage (0.00s)
=== RUN   TestBroadCastMessage
--- PASS: TestBroadCastMessage (0.00s)
=== RUN   TestLoginMessage
--- PASS: TestLoginMessage (0.00s)
=== RUN   TestLogoutMessage
--- PASS: TestLogoutMessage (0.00s)
=== RUN   TestReceiveMessage
--- PASS: TestReceiveMessage (0.00s)
=== RUN   TestGroupMessage
--- PASS: TestGroupMessage (0.00s)
=== RUN   TestLeaveMessage
--- PASS: TestLeaveMessage (0.00s)
=== RUN   TestWriteSendMessage
--- PASS: TestWriteSendMessage (0.00s)
=== RUN   TestWriteBroadCastMessage
--- PASS: TestWriteBroadCastMessage (0.00s)
=== RUN   TestWriteLoginMessage
--- PASS: TestWriteLoginMessage (0.00s)
=== RUN   TestWriteLogoutMessage
--- PASS: TestWriteLogoutMessage (0.00s)
=== RUN   TestWriteReceiveMessage
--- PASS: TestWriteReceiveMessage (0.00s)
=== RUN   TestWriteGroupMessage
--- PASS: TestWriteGroupMessage (0.00s)
=== RUN   TestWriteLeaveMessage
--- PASS: TestWriteLeaveMessage (0.00s)
PASS
ok      github.com/ZhengHe-MD/network-examples/tcp/chat/protocol        0.007s
```

心满意足。（完整测试请在源码中查看）

## 服务

服务需要实现的功能很简单，启动和关闭。我们可以用一个借口来定义它：

```go
// chat/serverserver.go
type ChatServer interface {
	Start(ctx context.Context, address string) error
	Close(ctx context.Context) error
}
```

以防万一未来有一天我抽风想更换传输层协议。

### TCP 服务器实现

先理清楚服务器本身需要保留哪些信息：

```go
// chat/server/tcp.go
type TcpChatServer struct {
  listener       net.Listener                   
  clientConnSet  map[*clientConn]interface{} // 保存 tcp 长连接，以及每个连接的标识符，即用户名
	groupToMembers map[string][]string         // 保存聊天组及其组内用户名称
	mu             *sync.Mutex 		             // 保护 clientConnSet 和 groupSet 的访问线程安全
}

type clientConn struct {
  conn   net.Conn
	name   string
	writer *protocol.CommandWriter
}
```

服务器运作的基本原理就是：

1. 主线程监听某个地址 (ip:port)，循环地等待连接请求
2. 主线程接到连接请求，建立连接，并创建一个新的线程负责处理连接断开之前的所有通信活动，主进程继续循环等待连接请求

我们可以这样实现它：

```go
// chat/server/tcp.go
func (s *TcpChatServer) Start(ctx context.Context, address string) error {
	if err := s.listen(ctx, address); err != nil {
		return err
	}

  // 主线程循环等待
Loop:
	for {
		select {
		case <-ctx.Done():
			log.Println("chat server is shutting down...")
			// TODO: implement elegant shutdown logic
			log.Println("shutdown successfully")
			break Loop
		default:
		}
		// 收到新的连接请求
		conn, err := s.listener.Accept()

		if err != nil {
			log.Print(err)
			continue
		}
    
		clientConn := s.accept(conn)
		go s.serve(clientConn)
	}

	return nil
}
```

其中 accept 函数负责建立 clientConn 结构体并存入 clientConnSet 中：

```go
// chat/server/tcp.go
func (s *TcpChatServer) accept(conn net.Conn) *clientConn {
	log.Printf("Accepting connection from %s", conn.RemoteAddr().String())

	s.mu.Lock()
	defer s.mu.Unlock()

	cc := &clientConn{
		conn:   conn,
		writer: protocol.NewCommandWriter(conn),
	}

	s.clientConnSet[cc] = struct{}{}
	return cc
}
```

serve 函数内部利用一个循环，在连接断开之前，负责处理与该客户端连接的后续通信工作：

```go
// chat/server/tcp.go
func (s *TcpChatServer) serve(cc *clientConn) {
	mr := protocol.NewCommandReader(cc.conn)
	defer s.remove(cc)

	for {
		var err error

		cmd, err := mr.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("read message err:%v", err)
			// https://github.com/golang/go/blob/f686a2890b34996455c7d7aba9a0efba74b613f5/src/net/error_test.go#L506
			if strings.Contains(err.Error(), ClosedConnectionMsg) {
				break
			}
			continue
		}
    
    // 这里甚至可以实现服务的中间件逻辑

		if cmd != nil {
			switch v := cmd.(type) {
			case *protocol.SendCommand:
				err = s.handleSend(cc, cmd.(*protocol.SendCommand))
			case *protocol.BroadCastCommand:
				err = s.handleBroadcast(cc, cmd.(*protocol.BroadCastCommand))
			case *protocol.LoginCommand:
				err = s.handleLogin(cc, cmd.(*protocol.LoginCommand))
			case *protocol.LogoutCommand:
				err = s.handleLogout(cc, cmd.(*protocol.LogoutCommand))
			case *protocol.GroupCommand:
				err = s.handleGroup(cc, cmd.(*protocol.GroupCommand))
			case *protocol.LeaveCommand:
				err = s.handleLeave(cc, cmd.(*protocol.LeaveCommand))
			default:
				log.Printf("cmd:%T %v not supported", v, v)
			}

			if err != nil {
				log.Printf("handle cmd err:%v", err)
			}
		}
	}
}
```

这里我们看到了类似 HTTP 服务器里常见的 controller/handler，负责处理使用不同 HTTP 方法访问不同 url 的请求，在 CHAT 协议中，controller/handler 负责处理不同命令请求，接下来就分别看一下这些 handler：

**handleLogin**

handleLogin 负责登录用户，记录用户名信息：

```go
// chat/server/handler.go
func (s *TcpChatServer) handleLogin(cc *clientConn, cmd *protocol.LoginCommand) (err error) {
	cc.name = cmd.Username
	log.Printf("set username:%s", cc.name)
	return
}
```

**handleLogout**

handleLogout 负责登出用户，删除连接：

```go
// chat/server/handler.go
func (s *TcpChatServer) handleLogout(cc *clientConn, cmd *protocol.LogoutCommand) (err error) {
	delete(s.clientConnSet, cc)
	err = cc.conn.Close()
	log.Printf("user:%s logged out", cc.name)
	return
}
```

**handleSend**

handleSend 负责向指定用户发送消息：

```go
// chat/server/handler.go
func (s *TcpChatServer) handleSend(cc *clientConn, cmd *protocol.SendCommand) (err error) {
	for scc, _ := range s.clientConnSet {
		if scc.name == cmd.Name {
			// fail-fast
			if err = scc.writer.Write(&protocol.ReceiveCommand{
				BaseCommand: cmd.BaseCommand,
				From:        cc.name,
				Data:        cmd.Data,
			}); err != nil {
				return
			}
		}
	}
	return
}
```

**handleGroup**

handleGroup 负责建立群聊天：

```go
// chat/server/handler.go
func (s *TcpChatServer) handleGroup(cc *clientConn, cmd *protocol.GroupCommand) (err error) {
	s.mu.RLock()
	if _, ok := s.groupToMembers[cmd.GroupName]; ok {
		s.mu.RUnlock()
		log.Printf("group:%s exists", cmd.GroupName)
		return
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.groupToMembers[cmd.GroupName] = cmd.UserNames
	log.Printf("create group:%s", cmd.GroupName)
	return
}
```

**handleLeave**

handleLeave 帮助用户离开群聊天

```go
// chat/server/handler.go
func (s *TcpChatServer) handleLeave(cc *clientConn, cmd *protocol.LeaveCommand) (err error) {
	s.mu.RLock()
	if _, ok := s.groupToMembers[cmd.GroupName]; !ok {
		s.mu.RUnlock()
		log.Printf("group:%s doesn't exist", cmd.GroupName)
		return
	}

	var userNames []string
	for _, userName := range s.groupToMembers[cmd.GroupName] {
		if userName != cc.name {
			userNames = append(userNames, userName)
		}
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.groupToMembers[cmd.GroupName] = userNames
	log.Printf("%s leave group:%s", cc.name, cmd.GroupName)
	return
}
```

## 使用聊天服务

### 启动服务

```go
// chat/server/cmd/main.go
func main() {
  ctx := context.TODO()
  var s server.ChatServer
  s = server.NewTcpChatServer()
  s.Start(ctx, ":3333")
}
```

执行：

```sh
$ go run main.go
> 2019/07/10 09:08:04 Listening on :3333
```

服务启动完毕

### 访问服务

由于我们只实现了聊天服务端逻辑，要访问它还需要一个懂 tcp 的客户端，我们可以使用 nc 命令，即 netcat：

```sh
$ man nc
> ...
> NAME
>      nc -- arbitrary TCP and UDP connections and listens
> ...
```

整个演示过程如下所示：

![demo](http://g.recordit.co/IQvQfEDkPJ.gif)

看一下服务的输出：

```
2019/07/10 09:15:13 Listening on :3333
2019/07/10 09:15:35 Accepting connection from 127.0.0.1:64610
2019/07/10 09:15:39 set username:zhangsan
2019/07/10 09:15:43 Accepting connection from 127.0.0.1:64614
2019/07/10 09:15:47 set username:lisi
2019/07/10 09:15:53 Accepting connection from 127.0.0.1:64619
2019/07/10 09:15:57 set username:wangwu
2019/07/10 09:16:01 Accepting connection from 127.0.0.1:64620
2019/07/10 09:16:05 set username:zhaoliu
2019/07/10 09:16:26 create group:g1
2019/07/10 09:16:39 user:zhangsan logged out
2019/07/10 09:16:39 read message err:read tcp 127.0.0.1:3333->127.0.0.1:64610: use of closed network connection
2019/07/10 09:16:44 user:lisi logged out
2019/07/10 09:16:44 read message err:read tcp 127.0.0.1:3333->127.0.0.1:64614: use of closed network connection
2019/07/10 09:16:50 user:wangwu logged out
2019/07/10 09:16:50 read message err:read tcp 127.0.0.1:3333->127.0.0.1:64619: use of closed network connection
2019/07/10 09:16:53 user:zhaoliu logged out
2019/07/10 09:16:53 read message err:read tcp 127.0.0.1:3333->127.0.0.1:64620: use of closed network connection
```

## 小结

本文介绍如何使用 Go 语言，基于自定义的 CHAT 协议，从头开始搭建一个聊天服务。
希望你能有所收获。

## 参考

* [github: nqbao/learn-go/chatserver](https://github.com/nqbao/learn-go/tree/chat/0.0.1/chatserver)






























