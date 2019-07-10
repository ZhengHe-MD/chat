package main

import (
	"context"
	"github.com/ZhengHe-MD/network-examples/tcp/chat/server"
)

func main() {
	ctx := context.TODO()
	var s server.ChatServer
	s = server.NewTcpChatServer()
	s.Start(ctx, ":3333")
}