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

type ParserState int

const (
	ParserInitialized ParserState = iota
	ParserDone
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	r, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("Error reading from reader:", err)
	}

	lines := strings.Split(string(r), "\r\n")
	_, requestLine, err := parseRequestLine(lines[0])
	return &Request{RequestLine: *requestLine}, err
}

func parseRequestLine(line string) (int, *RequestLine, error) {
	if !strings.Contains(line, "\r\n") {
		// needs more data before we can parse the request line
		return 0, &RequestLine{}, nil
	}

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return -1, &RequestLine{}, fmt.Errorf("request line must have exactly 3 parts")
	}

	method := parts[0]
	if method != strings.ToUpper(method) {
		return -1, &RequestLine{}, fmt.Errorf("method must be uppercase")
	}

	httpVersion := strings.Split(parts[2], "/")[1]
	if httpVersion != "1.1" {
		return -1, &RequestLine{}, fmt.Errorf("http version must be 1.1")
	}

	bytesRead := len(line) + 2 // +2 for the \r\n at the end

	return bytesRead, &RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: parts[1],
		Method:        method,
	}, nil
}
