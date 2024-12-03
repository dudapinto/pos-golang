package limiter

import (
	"errors"
	"time"
)

type MockStore struct {
	store       map[string]*RateLimit
	GetFn       func(key string) (*RateLimit, error)
	SetFn       func(key string, rl *RateLimit) error
	IncrementFn func(key string) error
}

// Get retorna o RateLimit para uma chave
func (m *MockStore) Get(key string) (*RateLimit, error) {
	if m.GetFn != nil {
		return m.GetFn(key)
	}
	// Verifica se a chave existe no armazenamento
	rl, exists := m.store[key]
	if !exists {
		// Se não existir, retorna um erro específico
		return nil, errors.New("rate limit not found")
	}
	return rl, nil
}

// Set armazena um RateLimit para uma chave
func (m *MockStore) Set(key string, rl *RateLimit) error {
	if m.SetFn != nil {
		return m.SetFn(key, rl)
	}
	// Verifica se o RateLimit a ser armazenado é válido
	if rl == nil {
		return errors.New("invalid rate limit object")
	}
	// Armazena o RateLimit na chave
	m.store[key] = rl
	return nil
}

// Increment incrementa o contador de requisições de uma chave
func (m *MockStore) Increment(key string) error {
	if m.IncrementFn != nil {
		return m.IncrementFn(key)
	}
	// Verifica se a chave existe antes de incrementar
	if rl, exists := m.store[key]; exists {
		rl.Requests++
	} else {
		// Se não existir, cria um novo RateLimit
		m.store[key] = &RateLimit{
			Requests:       1,
			LimitPerSecond: 1,
			LastRequest:    time.Now(),
		}
	}
	return nil
}
