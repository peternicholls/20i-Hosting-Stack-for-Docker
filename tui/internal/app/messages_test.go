// Project: 20i Stack Manager TUI
// File: messages_test.go
// Purpose: Unit tests for message type creation and validation
// Version: 0.1.0
// Updated: 2025-12-28

package app

import "testing"

func TestContainerActionMsg_ValidActions(t *testing.T) {
	actions := []string{"start", "stop", "restart"}
	for _, action := range actions {
		msg := ContainerActionMsg{
			Action:      action,
			ContainerID: "abc123",
		}
		if msg.Action != action {
			t.Fatalf("expected action %q, got %q", action, msg.Action)
		}
		if msg.ContainerID == "" {
			t.Fatal("expected container ID to be set")
		}
	}
}

func TestComposeActionMsg_ValidActions(t *testing.T) {
	actions := []string{"stop", "restart", "down"}
	for _, action := range actions {
		msg := ComposeActionMsg{Action: action}
		if msg.Action != action {
			t.Fatalf("expected action %q, got %q", action, msg.Action)
		}
	}
}

func TestContainerActionResultMsg_Fields(t *testing.T) {
	msg := ContainerActionResultMsg{
		Success: true,
		Message: "ok",
	}
	if !msg.Success {
		t.Fatal("expected success to be true")
	}
	if msg.Message == "" {
		t.Fatal("expected message to be set")
	}
}

func TestComposeActionResultMsg_Fields(t *testing.T) {
	msg := ComposeActionResultMsg{
		Success: false,
		Message: "failed",
	}
	if msg.Success {
		t.Fatal("expected success to be false")
	}
	if msg.Message == "" {
		t.Fatal("expected message to be set")
	}
}
