package server

import (
	"github.com/ZhengHe-MD/network-examples/tcp/chat/protocol"
	"log"
)

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

func (s *TcpChatServer) handleBroadcast(cc *clientConn, cmd *protocol.BroadCastCommand) (err error) {
	userNames, ok := s.groupToMembers[cmd.GroupName]
	if !ok {
		log.Printf("group:%s doesn't exist", cmd.GroupName)
		return
	}

	userNameSet := make(map[string]interface{})
	for _, userName := range userNames {
		userNameSet[userName] = struct{}{}
	}

	for scc, _ := range s.clientConnSet {
		if scc == cc {
			continue
		}

		if _, ok := userNameSet[scc.name]; ok {
			err = scc.writer.Write(&protocol.ReceiveCommand{
				BaseCommand: cmd.BaseCommand,
				From:        cc.name,
				Data:        cmd.Data,
			})
		}
	}
	return
}

func (s *TcpChatServer) handleLogin(cc *clientConn, cmd *protocol.LoginCommand) (err error) {
	cc.name = cmd.Username
	log.Printf("set username:%s", cc.name)
	return
}

func (s *TcpChatServer) handleLogout(cc *clientConn, cmd *protocol.LogoutCommand) (err error) {
	delete(s.clientConnSet, cc)
	err = cc.conn.Close()
	log.Printf("user:%s logged out", cc.name)
	return
}

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
