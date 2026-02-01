package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("       Host:        localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 44, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nSomeOtherHeader: does exist\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)
	assert.Equal(t, 0, len(headers))

	// Test: Valid single header with capital letters
	headers = NewHeaders()
	data = []byte("HoSt: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid Missing colon
	headers = NewHeaders()
	data = []byte("HostWithoutColon\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err, "invalid header: missing colon")
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid Empty key (colon at start)
	headers = NewHeaders()
	data = []byte(": value\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err, "invalid header: empty key")
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Only whitespace before colon
	headers = NewHeaders()
	data = []byte("   : value\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err, "invalid header: empty key")
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid char in key
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err, "invalid header: invalid characters in field name")
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
