package handlers

import (
	"server/models"
	"testing"

	"github.com/gorilla/websocket"
)

func TestGetHandler_ReturnsCorrectFunction(t *testing.T) {
	handler := GetHandler("register")
	if handler == nil {
		t.Fatalf("Expected handler for 'register', got nil")
	}
}

func TestHandlerFunctionGetsCalled(t *testing.T) {
	// Setup
	called := false
	fakeHandler := func(conn *websocket.Conn, msg models.ReceiveMessage) error {
		called = true
		return nil
	}

	// Replace the actual handler temporarily
	handlersBackup := handlers
	defer func() { handlers = handlersBackup }() // Restore after test
	handlers = map[string]models.Handler{
		"test_type": fakeHandler,
	}

	// Call
	handlerFunc := GetHandler("test_type")
	if handlerFunc == nil {
		t.Fatal("Expected handler function, got nil")
	}

	err := handlerFunc(nil, models.ReceiveMessage{Type: "test_type", Content: "hi"})
	if err != nil {
		t.Fatalf("Handler returned unexpected error: %v", err)
	}

	// Assert
	if !called {
		t.Fatal("Expected handler function to be called")
	}
}
