package protocol

import (
	"errors"
	"fmt"
	"strings"
)

// examples
// CHAT/1.0 SEND Body[name data]\n
// CHAT/1.0 LOGIN Body[username]\n
// CHAT/1.0 LOGOUT\n
// CHAT/1.0 RECEIVE Body[from data]\n
// CHAT/1.0 GROUP Body[groupname username ...]\n
// CHAT/1.0 LEAVE Body[groupname]\n
// CHAT/1.0 BROADCAST Body[groupname data]\n

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

var (
	InvalidMessageErr = errors.New("invalid message")
	UnsupportedCmdErr = errors.New("unsupported cmd")
)

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

type BroadCastCommand struct {
	BaseCommand
	GroupName string
	Data      []byte
}

func (c *BroadCastCommand) String() string {
	return strings.Join([]string{
		c.BaseCommand.String(),
		CmdBroadCast,
		c.GroupName,
		string(c.Data),
	}, ProtocolSep) + "\n"
}

type LoginCommand struct {
	BaseCommand
	Username string
}

func (c *LoginCommand) String() string {
	return strings.Join([]string{
		c.BaseCommand.String(),
		CmdLogin,
		c.Username,
	}, ProtocolSep) + "\n"
}

type LogoutCommand struct {
	BaseCommand
}

func (c *LogoutCommand) String() string {
	return strings.Join([]string{
		c.BaseCommand.String(),
		CmdLogout,
	}, ProtocolSep) + "\n"
}

type ReceiveCommand struct {
	BaseCommand
	From string
	Data []byte
}

func (c *ReceiveCommand) String() string {
	return strings.Join([]string{
		c.BaseCommand.String(),
		CmdReceive,
		c.From,
		string(c.Data),
	}, ProtocolSep) + "\n"
}

type GroupCommand struct {
	BaseCommand
	GroupName string
	UserNames []string
}

func (c *GroupCommand) String() string {
	return strings.Join([]string{
		c.BaseCommand.String(),
		CmdGroup,
		c.GroupName,
		strings.Join(c.UserNames, ProtocolSep),
	}, ProtocolSep) + "\n"
}

type LeaveCommand struct {
	BaseCommand
	GroupName string
}

func (c *LeaveCommand) String() string {
	return strings.Join([]string{
		c.BaseCommand.String(),
		CmdLeave,
		c.GroupName,
	}, ProtocolSep) + "\n"
}
