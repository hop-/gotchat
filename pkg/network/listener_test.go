package network

import (
	"errors"
	"testing"
)

func TestListener_Accept_Success(t *testing.T) {
	mockListener := NewMockBasicListener(t)
	mockConn := NewMockBasicConn(t)

	mockListener.On("Accept").Return(mockConn, nil)

	listener := NewListener(mockListener)
	conn, err := listener.Accept()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if conn == nil {
		t.Errorf("expected connection, got nil")
	}
	mockListener.AssertExpectations(t)
}

func TestListener_Accept_Error(t *testing.T) {
	mockListener := new(MockBasicListener)

	mockListener.On("Accept").Return(nil, errors.New("accept error"))

	listener := NewListener(mockListener)
	conn, err := listener.Accept()

	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if conn != nil {
		t.Errorf("expected nil connection, got %v", conn)
	}
	if err == nil || err.Error() != "accept error" {
		t.Errorf("expected error 'accept error', got %v", err)
	}
	mockListener.AssertExpectations(t)
}
