package server

import (
	"context"
	"github.com/ZhengHe-MD/network-examples/tcp/chat/protocol"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type clientConn struct {
	conn   net.Conn
	name   string
	writer *protocol.CommandWriter
}

type TcpChatServer struct {
	listener       net.Listener
	clientConnSet  map[*clientConn]interface{}
	groupToMembers map[string][]string
	mu             *sync.RWMutex
}

func NewTcpChatServer() *TcpChatServer {
	return &TcpChatServer{
		mu:             &sync.RWMutex{},
		clientConnSet:  make(map[*clientConn]interface{}),
		groupToMembers: make(map[string][]string),
	}
}

func (s *TcpChatServer) listen(ctx context.Context, address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.listener = l

	log.Printf("Listening on %s", address)

	return err
}

func (s *TcpChatServer) Close(ctx context.Context) error {
	return s.listener.Close()
}

func (s *TcpChatServer) Start(ctx context.Context, address string) error {
	if err := s.listen(ctx, address); err != nil {
		return err
	}

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

func (s *TcpChatServer) remove(cc *clientConn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clientConnSet, cc)
}

const ClosedConnectionMsg = "use of closed network connection"

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
