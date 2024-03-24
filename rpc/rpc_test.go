package rpc_test

import (
	"testing"

	"edulsp/rpc"
)

type EncodingExample struct {
	Method string
}

func TestEncodeMessage(t *testing.T) {
	want := "Content-Length: 18\r\n\r\n{\"Method\":\"hello\"}"
	got := rpc.EncodeMessage(EncodingExample{Method: "hello"})
	if got != want {
		t.Fatalf("want: %q, got: %q", want, got)
	}
}

func TestDecodeMessage(t *testing.T) {
	incomingMessage := "Content-Length: 18\r\nContent-Type: application/vscode-jsonrpc; charset=utf-8\r\n\r\n{\"method\":\"hello\"}"
	method, content, err := rpc.DecodeMessage([]byte(incomingMessage))
	if err != nil {
		t.Fatal(err)
	}

	if len(content) != 18 {
		t.Fatalf("want: 18, got: %d", len(content))
	}

	if method != "hello" {
		t.Fatalf(`want: "hello", got: %s`, method)
	}
}
