package network

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Message struct {
	headers map[string]string
	body    []byte
}

func NewMessage(headers map[string]string, body []byte) *Message {
	return &Message{
		headers: headers,
		body:    body,
	}
}

func (m *Message) Headers() map[string]string {
	return m.headers
}

func (m *Message) Body() []byte {
	return m.body
}

func (m *Message) BodyTo(v any) error {
	if len(m.body) == 0 {
		return fmt.Errorf("message body is empty")
	}

	err := json.Unmarshal(m.body, v)
	if err != nil {
		return fmt.Errorf("failed to unmarshal message body: %w", err)
	}

	return nil
}

func (m *Message) SetHeader(key, value string) {
	m.headers[key] = value
}

func (m *Message) GetHeader(key string) (string, bool) {
	value, exists := m.headers[key]
	return value, exists
}

func (m *Message) DeleteHeader(key string) {
	delete(m.headers, key)
}

func newMessageFromBytes(headersData []byte, bodyData []byte) (*Message, error) {
	headers, err := bytesToHeaders(headersData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse headers: %w", err)
	}

	return &Message{
		headers: headers,
		body:    bodyData,
	}, nil
}

func (m *Message) toBytes() ([]byte, []byte) {
	headersData := headersToBytes(m.headers)

	return headersData, m.body
}

func bytesToHeaders(headerData []byte) (map[string]string, error) {
	headers := make(map[string]string)
	lines := bytes.Split(headerData, []byte("\r\n"))

	for _, line := range lines {
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header line: %s", line)
		}
		headers[string(parts[0])] = string(parts[1])
	}

	return headers, nil
}

func headersToBytes(headers map[string]string) []byte {
	var buffer bytes.Buffer

	for key, value := range headers {
		buffer.WriteString(key)
		buffer.WriteString(":")
		buffer.WriteString(value)
		buffer.WriteString("\r\n")
	}

	return buffer.Bytes()
}
