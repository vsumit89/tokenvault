package tokenvault

import (
	"fmt"
	"sync/atomic"
	"time"
)

type TokenGeneratorFunc func() (interface{}, error)

type TokenManager struct {
	token                  atomic.Value // Safely handle atomic load and store of a string
	duration               time.Duration
	generatorFunc          TokenGeneratorFunc
	name                   string
	hasGeneratedFirstToken bool
	blockFirstTime         chan bool
}

func NewTokenManager(name string, duration *time.Duration, generatorFunc TokenGeneratorFunc) *TokenManager {
	tm := &TokenManager{
		duration:               *duration,
		generatorFunc:          generatorFunc,
		name:                   name,
		blockFirstTime:         make(chan bool, 2), // keeping the size 2 so that only reading is blocked and writing is unblocked
		hasGeneratedFirstToken: false,
	}

	tm.token.Store("") // Initialize with an empty token
	return tm
}

func (t *TokenManager) GetToken() interface{} {
	if !t.hasGeneratedFirstToken {
		<-t.blockFirstTime
	}

	return t.token.Load()
}

func (tm *TokenManager) RunTokenGenerator() {
	ticker := time.NewTicker(tm.duration)
	defer ticker.Stop()

	// Immediately generate the first token
	tm.UpdateToken()

	for range ticker.C {
		tm.UpdateToken()
	}
}

// UpdateToken safely updates the token stored in the manager
func (tm *TokenManager) UpdateToken() {
	// Wait for first token to be generated
	fmt.Println("Generating new session token for connecting to " + tm.name)

	response, err := tm.generatorFunc()
	if err != nil {
		fmt.Println("Error updating token:", err)
		return
	}

	tm.token.Store(response)

	if !tm.hasGeneratedFirstToken {
		tm.hasGeneratedFirstToken = true
		tm.blockFirstTime <- true
	}

	fmt.Println("Successfully generated new token for " + tm.name)
}
