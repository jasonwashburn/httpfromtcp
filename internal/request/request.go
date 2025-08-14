package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	ParserState ParserState
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

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	r := Request{ParserState: ParserInitialized}

	for r.ParserState != ParserDone {
		if len(buf) == cap(buf) {
			newBuff := make([]byte, 2*len(buf))
			_ = copy(newBuff, buf)
			buf = newBuff
		}

		bytesRead, err := reader.Read(buf[readToIndex:])
		if errors.Is(err, io.EOF) {
			break
		}
		readToIndex += bytesRead
		numParsed, err := r.parse(buf[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("error parsing request: %w", err)
		}
		newBuff := make([]byte, len(buf[numParsed:]))
		copy(newBuff, buf[numParsed:])
		buf = newBuff
		readToIndex -= numParsed
	}
	return &r, nil
}

func parseRequestLine(line string) (int, *RequestLine, error) {
	if !strings.Contains(line, "\r\n") {
		// needs more data before we can parse the request line
		return 0, &RequestLine{}, nil
	}

	line = strings.Split(line, "\r\n")[0]
	fmt.Printf("Parsing request line: %s\n", line)

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

	fmt.Printf("Parsed request line: method=%#v, target=%#v, version=%#v\n", method, parts[1], httpVersion)

	rl := RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: parts[1],
		Method:        method,
	}
	fmt.Printf("Returning request line: %#v\n", rl)
	return bytesRead, &rl, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.ParserState {
	case ParserInitialized:
		n, rl, err := parseRequestLine(string(data))
		if err != nil {
			return -1, err
		}
		if n == 0 {
			return 0, nil // need more data
		}
		r.RequestLine = *rl
		r.ParserState = ParserDone
		return n, nil
	case ParserDone:
		return -1, errors.New("trying to read data in a done state")
	default:
		return -1, errors.New("unknown state")
	}
}
