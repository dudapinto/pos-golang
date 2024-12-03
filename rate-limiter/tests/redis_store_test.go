package tests

import (
	"testing"
	"time"

	"github.com/dudapinto/pos-golang/rate-limiter/internal/limiter"

	"github.com/stretchr/testify/assert"
)

func TestRedisStore(t *testing.T) {
	// Configuração do RedisStore
	store := limiter.NewRedisStore("localhost:6379", "")

	// Chave e RateLimit para teste
	key := "test_key"
	rateLimit := &limiter.RateLimit{
		LimitPerSecond: 10,
		BlockDuration:  60 * time.Second,
		LastRequest:    time.Now(),
		Requests:       1,
	}

	// Teste de Set
	err := store.Set(key, rateLimit)
	assert.NoError(t, err, "Set deveria funcionar sem erros")

	// Teste de Get
	retrieved, err := store.Get(key)
	assert.NoError(t, err, "Get deveria funcionar sem erros")
	assert.Equal(t, rateLimit.LimitPerSecond, retrieved.LimitPerSecond, "Os limites devem coincidir")
	assert.Equal(t, rateLimit.Requests, retrieved.Requests, "O número de requisições deve coincidir")

	// Teste de Increment
	err = store.Increment(key)
	assert.NoError(t, err, "Increment deveria funcionar sem erros")

	// Verifique se o número de requisições aumentou
	retrieved, err = store.Get(key)
	assert.NoError(t, err, "Get deveria funcionar após o Increment")
	assert.Equal(t, rateLimit.Requests+1, retrieved.Requests, "O número de requisições deve ter aumentado")
}
