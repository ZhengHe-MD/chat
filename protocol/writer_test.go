package protocol

import (
	"bytes"
	"testing"
)

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

func TestWriteBroadCastMessage(t *testing.T) {
	cases := []struct {
		cmd             *BroadCastCommand
		expectedMessage string
	}{
		{
			&BroadCastCommand{
				BaseCommand: BaseCommand{ProtocolName, ProtocolVersion},
				GroupName:   "g1",
				Data:        []byte("hello world"),
			},
			"CHAT/1.0 BROADCAST g1 hello world\n",
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

func TestWriteLoginMessage(t *testing.T) {
	cases := []struct {
		cmd             *LoginCommand
		expectedMessage string
	}{
		{
			&LoginCommand{
				BaseCommand: BaseCommand{ProtocolName, ProtocolVersion},
				Username:    "zhenghe",
			},
			"CHAT/1.0 LOGIN zhenghe\n",
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

func TestWriteLogoutMessage(t *testing.T) {
	cases := []struct {
		cmd             *LogoutCommand
		expectedMessage string
	}{
		{
			&LogoutCommand{
				BaseCommand: BaseCommand{ProtocolName, ProtocolVersion},
			},
			"CHAT/1.0 LOGOUT\n",
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

func TestWriteReceiveMessage(t *testing.T) {
	cases := []struct {
		cmd             *ReceiveCommand
		expectedMessage string
	}{
		{
			&ReceiveCommand{
				BaseCommand: BaseCommand{ProtocolName, ProtocolVersion},
				From:        "zhenghe",
				Data:        []byte("hello world"),
			},
			"CHAT/1.0 RECEIVE zhenghe hello world\n",
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

func TestWriteGroupMessage(t *testing.T) {
	cases := []struct {
		cmd             *GroupCommand
		expectedMessage string
	}{
		{
			&GroupCommand{
				BaseCommand: BaseCommand{ProtocolName, ProtocolVersion},
				GroupName:   "g1",
				UserNames:   []string{"zhenghe", "xixi"},
			},
			"CHAT/1.0 GROUP g1 zhenghe xixi\n",
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

func TestWriteLeaveMessage(t *testing.T) {
	cases := []struct {
		cmd             *LeaveCommand
		expectedMessage string
	}{
		{
			&LeaveCommand{
				BaseCommand: BaseCommand{ProtocolName, ProtocolVersion},
				GroupName:   "g1",
			},
			"CHAT/1.0 LEAVE g1\n",
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
