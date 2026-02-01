package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// Look for CRLF
	crlfIdx := bytes.Index(data, []byte("\r\n"))
	if crlfIdx == -1 {
		return 0, false, nil
	}

	// If CRLF is at the start, we've reached the end of headers
	if crlfIdx == 0 {
		return 2, true, nil
	}

	// Extract the header line (everything before CRLF)
	headerLine := data[:crlfIdx]
	bytesConsumed := crlfIdx + 2

	// Find the colon separator
	colonIdx := bytes.IndexByte(headerLine, ':')
	if colonIdx == -1 {
		return 0, false, errors.New("invalid header: missing colon")
	}

	// Extract key and value
	key := string(bytes.TrimSpace(headerLine[:colonIdx]))
	value := string(bytes.TrimSpace(headerLine[colonIdx+1:]))

	if key == "" {
		return 0, false, errors.New("invalid header: empty key")
	}

	// Check for space before colon (key shouldn't have trailing space in original)
	if colonIdx > 0 && headerLine[colonIdx-1] == ' ' {
		return 0, false, errors.New("invalid header: space before colon")
	}

	key = strings.ToLower(key)

	h[key] = value
	return bytesConsumed, false, nil
}

func NewHeaders() Headers {
	return map[string]string{}
}
