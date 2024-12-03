// tests/store_test.go
package tests

import (
	"testing"

	"github.com/dudapinto/pos-golang/rate-limiter/config"
	"github.com/dudapinto/pos-golang/rate-limiter/internal/limiter"
)

func TestMockStore(t *testing.T) {
	// Cria uma instância do MockStore com as funções mock
	mockStore := &limiter.MockStore{
		GetFn: func(key string) (*limiter.RateLimit, error) {
			// Simula um erro ou sucesso no Get
			return nil, nil
		},
		SetFn: func(key string, rl *limiter.RateLimit) error {
			// Simula sucesso no Set
			return nil
		},
		IncrementFn: func(key string) error {
			// Simula sucesso no Increment
			return nil
		},
	}

	// Inicializa o RateLimiter com o MockStore
	rateLimiter := &limiter.RateLimiter{
		Store: mockStore,
		// outros campos necessários
	}

	allowed := rateLimiter.AllowRequest("some-key")
	if !allowed { // Se a resposta for false, significa que a requisição foi bloqueada
		t.Errorf("Expected true, got false")
	}
}

// Testa a criação de Store com o MockStore através da função NewStore
func TestNewStoreWithMock(t *testing.T) {
	cfg := &config.Config{
		StorageType: "mock", // Configura o armazenamento para mock
	}

	store := limiter.NewStore(cfg)
	if store == nil {
		t.Fatalf("Expected store, got nil")
	}

	// Verifica se a instância criada é uma MockStore
	_, ok := store.(*limiter.MockStore)
	if !ok {
		t.Errorf("Expected MockStore, got %T", store)
	}
}
