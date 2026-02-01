// Package headers provides functionality to parse and store HTTP-like headers.
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

	if !isValidFieldName(key) {
		return 0, false, errors.New("invalid header: invalid characters in field name")
	}

	key = strings.ToLower(key)

	h.Add(key, value)
	return bytesConsumed, false, nil
}

func (h Headers) Add(k string, v string) {
	_, exists := h[k]
	if exists {
		h[k] = h[k] + ", " + v
	} else {
		h[k] = v
	}
}

func isValidFieldName(s string) bool {
	for _, r := range s {
		valid := (r >= 'A' && r <= 'Z') ||
			(r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9') ||
			r == '!' || r == '#' || r == '$' || r == '%' || r == '&' ||
			r == '\'' || r == '*' || r == '+' || r == '-' || r == '.' ||
			r == '^' || r == '_' || r == '`' || r == '|' || r == '~'

		if !valid {
			return false
		}
	}
	return true
}

func NewHeaders() Headers {
	return map[string]string{}
}
