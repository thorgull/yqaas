package jq

import (
	"bytes"
	"context"
	"fmt"
	"gopkg.in/op/go-logging.v1"
	"os/exec"
	"strings"
)

var log = logging.MustGetLogger("yqaas-jq")

type JQ interface {
	Version() (string, bool)
	Evaluate(expression string, data []byte) ([]byte, error)
}

type jqCommand struct {
}

func (j *jqCommand) Version() (string, bool) {
	cmd := exec.CommandContext(context.Background(), "jq", "--version")
	var output = bytes.NewBufferString("")
	cmd.Stdout = output
	err := cmd.Run()
	if err != nil {
		log.Warningf("can not launch jq to get version... %s", err.Error())
		return "", false
	}

	return strings.TrimPrefix(strings.TrimSuffix(output.String(), "\n"), "jq-"), true
}

func (j *jqCommand) Evaluate(expression string, data []byte) ([]byte, error) {
	cmd := exec.CommandContext(context.Background(), "jq", expression)
	var output = bytes.NewBuffer(make([]byte, 0))
	cmd.Stdin = bytes.NewReader(data)
	cmd.Stdout = output

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error while running jq : %w", err)
	}

	return output.Bytes(), nil
}

func NewJQCommand() JQ {
	return &jqCommand{}
}
