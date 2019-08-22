package main

import (
	"bufio"
	"errors"
	"io"
)

//Commands
const (
	DeleteKey byte = iota
	SetString
	GetString
)

// CommandExecutor ...
type CommandExecutor struct {
	input   *bufio.Reader
	storage Storage
}

// NewCommandExecutor ...
func NewCommandExecutor(input io.Reader, storage Storage) *CommandExecutor {
	return &CommandExecutor{input: bufio.NewReader(input), storage: storage}
}

// Execute ...
func (c *CommandExecutor) Execute() (output []byte, err error) {
	out := bufio.NewWriter(nil)
	switch c.readByte() {
	case DeleteKey:
		c.execDeleteKey(out)
	case SetString:
		c.execSetString(out)
	case GetString:
		c.execGetString(out)
	default:
		err = errors.New("Unknown command")
	}
	return
}

func (c *CommandExecutor) readByte() byte {
	b, err := c.input.ReadByte()
	if err != nil {
		panic(err)
	}
	return b
}

func (c *CommandExecutor) readString() string {
	sz := int(c.readByte())
	buf := make([]byte, sz)
	n, err := c.input.Read(buf)
	if n != sz || err != nil {
		panic(err)
	}
	return string(buf)
}

func (c *CommandExecutor) execDeleteKey(out *bufio.Writer) {
	key := c.readString()
	c.storage.DeleteKey(key)
	return
}

func (c *CommandExecutor) execSetString(out *bufio.Writer) {
	key := c.readString()
	value := c.readString()
	c.storage.SetString(key, value)
}

func (c *CommandExecutor) execGetString(out *bufio.Writer) {
	key := c.readString()
	value, ok := c.storage.GetString(key)
	if ok {
		out.WriteRune('"')
		out.WriteString(value)
		out.WriteRune('"')
	} else {
		out.WriteString("nil")
	}
}
