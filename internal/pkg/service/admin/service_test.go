package admin

import "testing"

func TestDummy(t *testing.T) {}

func TestNewService(t *testing.T) {
	srv := NewService(nil)
	if srv == nil {
		t.Fatal("NewService(nil) returned nil")
	}
}
