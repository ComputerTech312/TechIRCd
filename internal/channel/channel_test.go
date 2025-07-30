package main

import (
	"testing"
)

func TestNewChannel(t *testing.T) {
	ch := NewChannel("#test")

	if ch.name != "#test" {
		t.Errorf("Expected channel name #test, got %s", ch.name)
	}

	if ch.clients == nil {
		t.Error("Expected clients map to be initialized")
	}

	if ch.operators == nil {
		t.Error("Expected operators map to be initialized")
	}

	if ch.quietList == nil {
		t.Error("Expected quietList slice to be initialized")
	}
}

func TestChannelModes(t *testing.T) {
	ch := NewChannel("#test")

	// Test setting mode
	ch.SetMode('m', true)
	if !ch.HasMode('m') {
		t.Error("Expected channel to have mode +m")
	}

	// Test unsetting mode
	ch.SetMode('m', false)
	if ch.HasMode('m') {
		t.Error("Expected channel to not have mode +m")
	}
}

func TestChannelPermissions(t *testing.T) {
	ch := NewChannel("#test")

	// Mock client for testing
	// This would need actual client implementation
	// TODO: Add proper client mocking
}

// Benchmark tests
func BenchmarkChannelCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewChannel("#test")
	}
}

func BenchmarkModeSet(b *testing.B) {
	ch := NewChannel("#test")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ch.SetMode('m', true)
		ch.SetMode('m', false)
	}
}
