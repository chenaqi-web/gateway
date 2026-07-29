package config

import (
	"testing"
	"time"
)

func TestAuthConfigRefreshDuration(t *testing.T) {
	duration, err := (AuthConfig{RefreshExpire: "168h"}).RefreshDuration()
	if err != nil {
		t.Fatalf("RefreshDuration() error = %v", err)
	}
	if duration != 7*24*time.Hour {
		t.Fatalf("RefreshDuration() = %v, want %v", duration, 7*24*time.Hour)
	}
}

func TestAuthConfigRefreshDurationRejectsInvalidValue(t *testing.T) {
	if _, err := (AuthConfig{RefreshExpire: "invalid"}).RefreshDuration(); err == nil {
		t.Fatal("RefreshDuration() expected error for invalid duration")
	}
}
