package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// 	Mutate the Headers by adding newly parsed key-value pairs
	// Return n (the number of bytes consumed), done (whether or not it has finished parsing headers), and err (if it encountered an error)

	// Look for a CRLF, if it doesn't find one, assume you haven't been given enough data yet. Consume no data, return false for done, and nil for err.
	if !strings.Contains(string(data), "\r\n") {
		return 0, false, nil
	}
	// If you do find a CRLF, but it's at the start of the data, you've found the end of the headers, so return the proper values immediately.
	if bytes.HasPrefix(data, []byte("\r\n")) {
		return 2, true, nil
	}
	// Remove any extra whitespace from the key and value, but ensure there are no spaces between the colon and the key.
	workingString := strings.Split(string(data), "\r\n")[0]
	bytesConsumed := len(workingString) + 2
	workingString = strings.TrimSpace(workingString)

	splitString := strings.SplitN(workingString, ":", 2)
	key := splitString[0]
	value := strings.TrimSpace(splitString[1])

	if key[len(key)-1] == ' ' {
		return 0, false, errors.New("invalid header: space before colon")
	}

	// Assuming the format was valid (if it isn't return an error), add the key/value pair to the Headers map and return the number of bytes consumed. Note: The Parse function should only return done=true when the data starts with a CRLF, which can't happen when it finds a new key/value pair.

	h[key] = value
	return bytesConsumed, false, nil

	// It's important to understand that this function will be called over and over until all the headers are parsed, and it can only parse one key/value pair at a time.
}

func NewHeaders() Headers {
	return map[string]string{}
}
