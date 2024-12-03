// rate_limiter_mock_test.go
package tests

import (
	"testing"
	"time"

	"github.com/dudapinto/pos-golang/rate-limiter/config"
	"github.com/dudapinto/pos-golang/rate-limiter/internal/limiter"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiterWithMockStore(t *testing.T) {
	var rateLimit *limiter.RateLimit

	rateLimit = &limiter.RateLimit{
		LimitPerSecond: 5,
		BlockDuration:  time.Second,
		LastRequest:    time.Now().Add(-time.Second), // Simulando um pedido anterior dentro do limite de tempo
		Requests:       4,                            // 4 requisições já feitas
	}

	mockStore := &limiter.MockStore{
		GetFn: func(key string) (*limiter.RateLimit, error) {
			// Retornando um RateLimit com 4 requisições já feitas
			return rateLimit, nil
		},
		SetFn: func(key string, rl *limiter.RateLimit) error {
			// Simulando persistência de dados para o mock
			return nil
		},
		IncrementFn: func(key string) error {
			// Usando fechamento para acessar a instância do mockStore
			rateLimit.Requests++ // Incrementando o número de requisições
			return nil           // Atualizando a instância mockStore
		},
	}

	// Criando a configuração (do pacote config)
	cfg := &config.Config{
		LimitPerSecond: 5, // Limite de 5 requisições por segundo
		BlockDuration:  time.Second,
	}

	// Criando o rateLimiter usando a função NewRateLimiter do limiter
	rateLimiter := limiter.NewRateLimiter(mockStore, cfg)

	// Definindo a chave de teste
	key := "test_key"

	// Primeira requisição - deve ser permitida
	allowed := rateLimiter.AllowRequest(key)
	assert.True(t, allowed, "A primeira requisição deve ser permitida")

	// Segunda requisição - deve ser permitida, já que o limite de 4 foi atingido, mas o limite é 5
	allowed = rateLimiter.AllowRequest(key)
	assert.True(t, allowed, "A segunda requisição deve ser permitida")

	// Definindo os dados de RateLimit para refletir que 5 requisições já foram feitas
	rateLimiter.Limits[key] = limiter.RateLimitData{
		LimitPerSecond: 5,
		LastRequest:    time.Now(),
		Requests:       5, // Já atingiu o limite
	}

	// Terceira requisição - deve ser bloqueada, pois o limite de 5 requisições foi atingido
	allowed = rateLimiter.AllowRequest(key)
	assert.False(t, allowed, "A requisição deve ser bloqueada quando o limite for atingido")

	// Verificando o mockStore após a atualização
	rateLimit, err := mockStore.Get(key)
	assert.NoError(t, err, "Get deve funcionar sem erros")
	assert.Equal(t, 6, rateLimit.Requests, "O número de requisições deve ser atualizado corretamente para 6")
}
