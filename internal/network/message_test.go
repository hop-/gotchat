package network

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

func TestNewMessage(t *testing.T) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer token",
	}
	body := []byte(`{"key":"value"}`)

	message := NewMessage(headers, body)

	if !reflect.DeepEqual(message.Headers(), headers) {
		t.Errorf("expected headers %v, got %v", headers, message.Headers())
	}

	if !reflect.DeepEqual(message.Body(), body) {
		t.Errorf("expected body %s, got %s", body, message.Body())
	}
}

func TestSetHeader(t *testing.T) {
	message := NewMessage(map[string]string{}, nil)
	message.SetHeader("Content-Type", "application/json")

	value, exists := message.GetHeader("Content-Type")
	if !exists || value != "application/json" {
		t.Errorf("expected header Content-Type to be application/json, got %v", value)
	}
}

func TestGetHeader(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}
	message := NewMessage(headers, nil)

	value, exists := message.GetHeader("Authorization")
	if !exists || value != "Bearer token" {
		t.Errorf("expected header Authorization to be Bearer token, got %v", value)
	}

	_, exists = message.GetHeader("Non-Existent")
	if exists {
		t.Errorf("expected header Non-Existent to not exist")
	}
}

func TestDeleteHeader(t *testing.T) {
	// Test deleting an existing header
	headers := map[string]string{"Authorization": "Bearer token"}
	message := NewMessage(headers, nil)

	message.DeleteHeader("Authorization")
	_, exists := message.GetHeader("Authorization")
	if exists {
		t.Errorf("expected header Authorization to be deleted")
	}

	// Test deleting a non-existent header
	message.DeleteHeader("Non-Existent")
	_, exists = message.GetHeader("Non-Existent")
	if exists {
		t.Errorf("expected header Non-Existent to not exist")
	}

	// Test deleting from an empty headers map
	message = NewMessage(map[string]string{}, nil)
	message.DeleteHeader("Authorization")
	_, exists = message.GetHeader("Authorization")
	if exists {
		t.Errorf("expected header Authorization to not exist in an empty headers map")
	}
}

func TestBodyTo(t *testing.T) {
	body := []byte(`{"key":"value"}`)
	message := NewMessage(nil, body)

	var result map[string]string
	err := message.BodyTo(&result)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := map[string]string{"key": "value"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestBodyTo_EmptyBody(t *testing.T) {
	message := NewMessage(nil, nil)

	var result map[string]string
	err := message.BodyTo(&result)
	if err == nil {
		t.Errorf("expected error for empty body, got nil")
	}
}

func TestNewMessageFromBytes(t *testing.T) {
	headersData := []byte("Content-Type:application/json\r\nAuthorization:Bearer token\r\n")
	bodyData := []byte(`{"key":"value"}`)

	message, err := newMessageFromBytes(headersData, bodyData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedHeaders := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer token",
	}
	if !reflect.DeepEqual(message.Headers(), expectedHeaders) {
		t.Errorf("expected headers %v, got %v", expectedHeaders, message.Headers())
	}

	if !reflect.DeepEqual(message.Body(), bodyData) {
		t.Errorf("expected body %s, got %s", bodyData, message.Body())
	}
}

func TestToBytes(t *testing.T) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer token",
	}
	body := []byte(`{"key":"value"}`)

	message := NewMessage(headers, body)
	headersData, bodyData := message.toBytes()

	expectedHeadersData := []byte("Content-Type:application/json\r\nAuthorization:Bearer token\r\n")
	expectedHeadersDataWithAnotherOrder := []byte("Authorization:Bearer token\r\nContent-Type:application/json\r\n")
	if !bytes.Equal(headersData, expectedHeadersData) && !bytes.Equal(headersData, expectedHeadersDataWithAnotherOrder) {
		t.Errorf("expected headers data %s, got %s", expectedHeadersData, headersData)
	}

	if !bytes.Equal(bodyData, body) {
		t.Errorf("expected body data %s, got %s", body, bodyData)
	}
}

func TestDeserializeMessage(t *testing.T) {
	// Test with valid data
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer token",
	}
	body := []byte(`{"key":"value"}`)

	headersData := headersToBytes(headers)
	headersSize := uint64(len(headersData))
	bodySize := uint64(len(body))

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, headersSize)
	binary.Write(&buf, binary.LittleEndian, bodySize)
	buf.Write(headersData)
	buf.Write(body)

	data := buf.Bytes()

	message, err := DeserializeMessage(data)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(message.Headers(), headers) {
		t.Errorf("expected headers %v, got %v", headers, message.Headers())
	}

	if !reflect.DeepEqual(message.Body(), body) {
		t.Errorf("expected body %s, got %s", body, message.Body())
	}

	// Test with truncated data
	truncatedData := data[:len(data)-5]
	_, err = DeserializeMessage(truncatedData)
	if err == nil {
		t.Errorf("expected error for truncated data, got nil")
	}

	// Test with invalid header size
	invalidHeaderSizeData := make([]byte, len(data))
	copy(invalidHeaderSizeData, data)
	binary.LittleEndian.PutUint64(invalidHeaderSizeData[:8], uint64(len(data)+100)) // Invalid size

	_, err = DeserializeMessage(invalidHeaderSizeData)
	if err == nil {
		t.Errorf("expected error for invalid header size, got nil")
	}

	// Test with invalid body size
	invalidBodySizeData := make([]byte, len(data))
	copy(invalidBodySizeData, data)
	binary.LittleEndian.PutUint64(invalidBodySizeData[8:16], uint64(len(data)+100)) // Invalid size

	_, err = DeserializeMessage(invalidBodySizeData)
	if err == nil {
		t.Errorf("expected error for invalid body size, got nil")
	}

	// Test with empty data
	_, err = DeserializeMessage([]byte{})
	if err == nil {
		t.Errorf("expected error for empty data, got nil")
	}

	// Test with only header size and body size but no actual data
	var emptyMessageBuff bytes.Buffer
	binary.Write(&emptyMessageBuff, binary.LittleEndian, uint64(0))
	binary.Write(&emptyMessageBuff, binary.LittleEndian, uint64(0))
	_, err = DeserializeMessage(emptyMessageBuff.Bytes())
	if err != nil {
		t.Errorf("expected no error for empty message, got %v", err)
	}
}

func TestSerializeMessage(t *testing.T) {
	// Test with valid headers and body
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer token",
	}
	body := []byte(`{"key":"value"}`)

	message := NewMessage(headers, body)
	serializedData, err := SerializeMessage(message)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Deserialize the serialized data to verify correctness
	deserializedMessage, err := DeserializeMessage(serializedData)
	if err != nil {
		t.Errorf("unexpected error during deserialization: %v", err)
	}

	if !reflect.DeepEqual(deserializedMessage.Headers(), headers) {
		t.Errorf("expected headers %v, got %v", headers, deserializedMessage.Headers())
	}

	if !reflect.DeepEqual(deserializedMessage.Body(), body) {
		t.Errorf("expected body %s, got %s", body, deserializedMessage.Body())
	}

	// Test with empty headers and body
	message = NewMessage(map[string]string{}, []byte{})
	serializedData, err = SerializeMessage(message)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	deserializedMessage, err = DeserializeMessage(serializedData)
	if err != nil {
		t.Errorf("unexpected error during deserialization: %v", err)
	}

	if len(deserializedMessage.Headers()) != 0 {
		t.Errorf("expected empty headers, got %v", deserializedMessage.Headers())
	}

	if len(deserializedMessage.Body()) != 0 {
		t.Errorf("expected empty body, got %s", deserializedMessage.Body())
	}
}
