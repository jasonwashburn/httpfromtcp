package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestLine, err := parseRequestLine(rawBytes)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func requestLineFromString(rawRequestLine string) (*RequestLine, error) {
	parts := strings.Split(rawRequestLine, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line: %s", rawRequestLine)
	}
	method := parts[0]
	target := parts[1]
	versionString := parts[2]

	for _, r := range method {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	if versionString != "HTTP/1.1" {
		return nil, fmt.Errorf("invalid HTTP version: %s", versionString)
	}
	httpVersion := strings.Split(versionString, "/")[1]

	return &RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: target,
		Method:        method,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, fmt.Errorf("could not find CRLF in request line")
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, err
	}

	return requestLine, nil
}
