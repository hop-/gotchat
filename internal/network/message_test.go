package network

import (
	"bytes"
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
	headers := map[string]string{"Authorization": "Bearer token"}
	message := NewMessage(headers, nil)

	message.DeleteHeader("Authorization")
	_, exists := message.GetHeader("Authorization")
	if exists {
		t.Errorf("expected header Authorization to be deleted")
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
	if !bytes.Equal(headersData, expectedHeadersData) {
		t.Errorf("expected headers data %s, got %s", expectedHeadersData, headersData)
	}

	if !bytes.Equal(bodyData, body) {
		t.Errorf("expected body data %s, got %s", body, bodyData)
	}
}
