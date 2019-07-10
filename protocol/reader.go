package protocol

import (
	"bufio"
	"io"
	"strings"
)

type CommandReader struct {
	reader *bufio.Reader
}

func NewCommandReader(reader io.Reader) *CommandReader {
	return &CommandReader{
		reader: bufio.NewReader(reader),
	}
}

func (r *CommandReader) Read() (cmd interface{}, err error) {
	line, _, err := r.reader.ReadLine()
	if err != nil {
		return
	}

	parts := strings.Split(string(line), ProtocolSep)

	if len(parts) < 2 {
		err = InvalidMessageErr
		return
	}

	proVer := parts[0]
	proVerParts := strings.Split(strings.TrimSpace(proVer), "/")
	if len(proVerParts) != 2 {
		err = InvalidMessageErr
		return
	}
	protocol, version := proVerParts[0], proVerParts[1]

	if protocol != ProtocolName || version != ProtocolVersion {
		err = InvalidMessageErr
		return
	}
	base := BaseCommand{protocol, version}

	cmdName := strings.TrimSpace(parts[1])

	switch cmdName {
	case CmdSend:
		if len(parts) < 4 {
			err = InvalidMessageErr
			return
		}

		name := strings.TrimSpace(parts[2])
		message := strings.Join(parts[3:], ProtocolSep)

		cmd = &SendCommand{base, name, []byte(message)}
		return
	case CmdBroadCast:
		if len(parts) < 4 {
			err = InvalidMessageErr
			return
		}

		groupName := strings.TrimSpace(parts[2])
		message := strings.Join(parts[3:], ProtocolSep)
		cmd = &BroadCastCommand{base, groupName, []byte(message)}
	case CmdLogin:
		if len(parts) != 3 {
			err = InvalidMessageErr
			return
		}

		username := strings.TrimSpace(parts[2])

		cmd = &LoginCommand{base, username}
	case CmdLogout:
		if len(parts) != 2 {
			err = InvalidMessageErr
			return
		}

		cmd = &LogoutCommand{base}
	case CmdReceive:
		if len(parts) < 4 {
			err = InvalidMessageErr
			return
		}

		f := strings.TrimSpace(parts[2])
		message := strings.Join(parts[3:], ProtocolSep)

		cmd = &ReceiveCommand{base, f, []byte(message)}
	case CmdGroup:
		if len(parts) < 5 {
			err = InvalidMessageErr
			return
		}

		groupName := strings.TrimSpace(parts[2])
		userNames := parts[3:]

		cmd = &GroupCommand{base, groupName, userNames}
	case CmdLeave:
		if len(parts) < 3 {
			err = InvalidMessageErr
			return
		}

		groupName := strings.TrimSpace(parts[2])
		cmd = &LeaveCommand{base, groupName}
	default:
		err = UnsupportedCmdErr
	}

	return
}
