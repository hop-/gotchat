package network

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListener_Accept_Success(t *testing.T) {
	mockListener := NewMockBasicListener(t)
	mockConn := NewMockBasicConn(t)

	mockListener.On("Accept").Return(mockConn, nil)

	listener := NewListener(mockListener)
	conn, err := listener.Accept()

	assert.NoError(t, err)
	assert.NotNil(t, conn)
	mockListener.AssertExpectations(t)
}

func TestListener_Accept_Error(t *testing.T) {
	mockListener := new(MockBasicListener)

	mockListener.On("Accept").Return(nil, errors.New("accept error"))

	listener := NewListener(mockListener)
	conn, err := listener.Accept()

	assert.Error(t, err)
	assert.Nil(t, conn)
	assert.EqualError(t, err, "accept error")
	mockListener.AssertExpectations(t)
}
