// rate_limiter_middleware_test.go
package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dudapinto/pos-golang/rate-limiter/config"
	"github.com/dudapinto/pos-golang/rate-limiter/internal/limiter"
	"github.com/dudapinto/pos-golang/rate-limiter/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {
	// Configuração
	cfg := &config.Config{
		LimitPerSecond: 1, // Limite baixo para facilitar o teste
		BlockDuration:  1 * time.Second,
	}

	var rateLimit *limiter.RateLimit

	// Criando o MockStore
	mockStore := &limiter.MockStore{
		GetFn: func(key string) (*limiter.RateLimit, error) {
			// Retornar um mock de RateLimit com RateLimitData
			return rateLimit, nil
		},
		SetFn: func(key string, rl *limiter.RateLimit) error {
			// Simular sucesso na configuração
			return nil
		},
		IncrementFn: func(key string) error {
			// Simular incremento
			rateLimit.Requests++
			return nil
		},
	}

	rateLimit = &limiter.RateLimit{
		Requests: 0,
	}

	// Criar o RateLimiter com o MockStore
	rl := limiter.NewRateLimiter(mockStore, cfg)
	mw := middleware.RateLimitMiddleware(rl)

	// Manipulador de teste
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Requisição de teste
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1"
	rec := httptest.NewRecorder()

	// Primeira requisição deve ser permitida
	mw(next).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code, "A primeira requisição deveria ser permitida")

	// Segunda requisição deve ser bloqueada
	rec = httptest.NewRecorder()
	mw(next).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code, "A segunda requisição deveria ser bloqueada")
}
