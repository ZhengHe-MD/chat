package protocol

import (
	"bufio"
	"fmt"
	"io"
)

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
