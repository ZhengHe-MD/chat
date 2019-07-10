package protocol

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestInvalidMessage(t *testing.T) {
	cases := []struct {
		message     string
		expectedErr error
	}{
		{
			"hello world\n",
			InvalidMessageErr,
		},
		{
			"HTTP/1.0 SEND hello\n",
			InvalidMessageErr,
		},
		{
			"CHAT/1.1 SEND hello\n",
			InvalidMessageErr,
		},
		{
			"CHAT/1.0 STAR hello\n",
			UnsupportedCmdErr,
		},
	}

	for i, c := range cases {
		mr := NewCommandReader(strings.NewReader(c.message))

		_, err := mr.Read()
		if err != c.expectedErr {
			t.Errorf("case %d: should have err:%v got:%v",
				i, c.expectedErr, err)
		}
	}
}

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

func TestBroadCastMessage(t *testing.T) {
	cases := []struct {
		message           string
		expectedErr       error
		expectedGroupName string
		expectedData      []byte
	}{
		{
			"CHAT/1.0 BROADCAST g1 hello world\n",
			nil,
			"g1",
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
			broadcastCmd := cmd.(*BroadCastCommand)

			if broadcastCmd.GroupName != c.expectedGroupName {
				t.Errorf("case %d: should have groupName:%s got:%s",
					i, c.expectedGroupName, broadcastCmd.GroupName)
			}

			if !bytes.Equal(broadcastCmd.Data, c.expectedData) {
				t.Errorf("case %d: should have data:%v got:%v",
					i, c.expectedData, broadcastCmd.Data)
			}
		}
	}
}

func TestLoginMessage(t *testing.T) {
	cases := []struct {
		message          string
		expectedErr      error
		expectedUsername string
	}{
		{
			"CHAT/1.0 LOGIN zhenghe\n",
			nil,
			"zhenghe",
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
			loginCmd := cmd.(*LoginCommand)
			if loginCmd.Username != c.expectedUsername {
				t.Errorf("case %d: should have name:%s got:%s",
					i, c.expectedUsername, loginCmd.Username)
			}
		}
	}
}

func TestLogoutMessage(t *testing.T) {
	cases := []struct {
		message     string
		expectedErr error
	}{
		{
			"CHAT/1.0 LOGOUT\n",
			nil,
		},
	}

	for i, c := range cases {
		mr := NewCommandReader(strings.NewReader(c.message))

		_, err := mr.Read()
		if err != c.expectedErr {
			t.Errorf("case %d: should have err:%v got:%v",
				i, c.expectedErr, err)
		}
	}
}

func TestReceiveMessage(t *testing.T) {
	cases := []struct {
		message      string
		expectedErr  error
		expectedFrom string
		expectedData []byte
	}{
		{
			"CHAT/1.0 RECEIVE zhenghe hello\n",
			nil,
			"zhenghe",
			[]byte("hello"),
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
			receiveCmd := cmd.(*ReceiveCommand)
			if receiveCmd.From != c.expectedFrom {
				t.Errorf("case %d: should have from:%s got:%s",
					i, c.expectedFrom, receiveCmd.From)
			}

			if !reflect.DeepEqual(receiveCmd.Data, c.expectedData) {
				t.Errorf("case %d: should have data:%v got:%v",
					i, c.expectedData, receiveCmd.Data)
			}
		}
	}
}

func TestGroupMessage(t *testing.T) {
	cases := []struct {
		message           string
		expectedErr       error
		expectedGroupName string
		expectedUserNames []string
	}{
		{
			"CHAT/1.0 GROUP g1 zhenghe xixi\n",
			nil,
			"g1",
			[]string{"zhenghe", "xixi"},
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
			groupCmd := cmd.(*GroupCommand)
			if groupCmd.GroupName != c.expectedGroupName {
				t.Errorf("case %d: should have groupName:%s got:%s",
					i, c.expectedGroupName, groupCmd.GroupName)
			}

			if !reflect.DeepEqual(groupCmd.UserNames, c.expectedUserNames) {
				t.Errorf("case %d: should have userNames:%v got:%v",
					i, c.expectedUserNames, groupCmd.UserNames)
			}
		}
	}
}

func TestLeaveMessage(t *testing.T) {
	cases := []struct {
		message           string
		expectedErr       error
		expectedGroupName string
	}{
		{
			"CHAT/1.0 LEAVE g1\n",
			nil,
			"g1",
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
			leaveCmd := cmd.(*LeaveCommand)

			if leaveCmd.GroupName != c.expectedGroupName {
				t.Errorf("case %d: should have groupName:%s got:%s",
					i, c.expectedGroupName, leaveCmd.GroupName)
			}
		}
	}
}
