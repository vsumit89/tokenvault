package tokenvault

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// Mock TokenGeneratorFunc
func mockTokenGenerator() (interface{}, error) {
	return "mock-token", nil
}

func TestNewTokenManager(t *testing.T) {
	// Test case: Create a new TokenManager
	duration := 10 * time.Second
	tm := NewTokenManager("test", &duration, mockTokenGenerator)

	if tm == nil {
		t.Errorf("NewTokenManager returned nil")
		return
	}

	if tm.name != "test" {
		t.Errorf("Incorrect name. Expected: 'test', Got: '%s'", tm.name)
	}

	if tm.duration != duration {
		t.Errorf("Incorrect duration. Expected: 10s, Got: %v", tm.duration)
	}

	if tm.generatorFunc == nil {
		t.Error("generatorFunc is nil")
	}

	if tm.hasGeneratedFirstToken {
		t.Error("hasGeneratedFirstToken should be false initially")
	}
}

func TestGetToken(t *testing.T) {
	// Test case: Get token before first token is generated
	duration := 10 * time.Second
	tm := NewTokenManager("test", &duration, mockTokenGenerator)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		_ = tm.GetToken()
	}()

	select {
	case <-doneCh:
		// Test ran without any issues within the timeout
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			t.Log("GetToken blocked as expected before first token generation")
		} else {
			t.Errorf("Unexpected error: %v", ctx.Err())
		}
	}

	// Test case: Get token after first token is generated
	tm.UpdateToken()
	token := tm.GetToken()

	if token != "mock-token" {
		t.Errorf("Incorrect token. Expected: 'mock-token', Got: '%v'", token)
	}
}

func TestUpdateToken(t *testing.T) {
	// Test case: Update token successfully
	duration := 1 * time.Second
	tm := NewTokenManager("test", &duration, mockTokenGenerator)

	go tm.RunTokenGenerator()

	time.Sleep(5 * time.Second)

	token := tm.GetToken()

	fmt.Println("token", token)
	if token != "mock-token" {
		t.Errorf("Incorrect token. Expected: 'mock-token', Got: '%v'", token)
	}

	if !tm.hasGeneratedFirstToken {
		t.Error("hasGeneratedFirstToken should be true after first token generation")
	}
}

func TestRunTokenGenerator(t *testing.T) {
	// Test case: Run token generator
	duration := 500 * time.Millisecond
	tm := NewTokenManager("test", &duration, mockTokenGenerator)

	// Start the token generator in a separate goroutine
	go tm.RunTokenGenerator()

	// Wait for the first token generation
	<-tm.blockFirstTime

	// Check if the first token was generated successfully
	token := tm.GetToken()
	if token != "mock-token" {
		t.Errorf("Incorrect token after first generation. Expected: 'mock-token', Got: '%v'", token)
	}

}
