package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	HeaderContentLength = "Content-Length"
	HeaderContentType   = "Content-Type"
)

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s: %d\r\n\r\n%s", HeaderContentLength, len(content), content)
}

type BaseMessage struct {
	Method string `json:"method"`
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	_, contentLength, content, err := parseMessage(msg)
	if err != nil {
		return "", nil, err
	}

	var baseMessage BaseMessage
	if err := json.Unmarshal(content, &baseMessage); err != nil {
		return "", nil, err
	}

	return baseMessage.Method, content[:contentLength], nil
}

func parseMessage(msg []byte) (messageLength, contentLength int, content []byte, err error) {
	header, content, found := bytes.Cut(msg, []byte("\r\n\r\n"))
	if !found {
		return
	}

	headerFields := make(map[string]any)
	if err = decodeHeader(header[:len(header)+2], headerFields); err != nil {
		return
	}

	contentLength = headerFields[HeaderContentLength].(int)
	messageLength = len(header) + 4 + contentLength

	return
}

func decodeHeader(header []byte, dst map[string]any) error {
	field, remaining, found := bytes.Cut(header, []byte("\r\n"))
	if !found {
		return nil
	}

	name, value, found := bytes.Cut(field, []byte(": "))
	if !found {
		return nil
	}

	switch n := string(name); n {
	case HeaderContentLength:
		contentLength, err := strconv.Atoi(string(value))
		if err != nil {
			return err
		}
		dst[n] = contentLength
	case HeaderContentType:
		dst[n] = string(value)
	default:
		return fmt.Errorf("unsupported field name: %s", n)
	}

	return decodeHeader(remaining, dst)
}

// type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)
func Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	messageLength, contentLength, content, err := parseMessage(data)
	if err != nil {
		return 0, nil, err
	}

	if len(content) < contentLength { // not ready yet
		return 0, nil, nil
	}

	return messageLength, data[:messageLength], nil
}
