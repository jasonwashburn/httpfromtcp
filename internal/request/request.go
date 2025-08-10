package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("Error reading from reader:", err)
	}

	lines := strings.Split(string(r), "\r\n")
	requestLine, err := parseRequestLine(lines[0])
	return &Request{RequestLine: *requestLine}, err
}

func parseRequestLine(line string) (*RequestLine, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return &RequestLine{}, fmt.Errorf("rquest line must have exactly 3 parts")
	}

	method := parts[0]
	if method != strings.ToUpper(method) {
		return &RequestLine{}, fmt.Errorf("method must be uppercase")
	}

	httpVersion := strings.Split(parts[2], "/")[1]
	if httpVersion != "1.1" {
		return &RequestLine{}, fmt.Errorf("http version must be 1.1")
	}

	return &RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: parts[1],
		Method:        method,
	}, nil
}
