// limiter.go
package limiter

import (
	"log"
	"time"

	"github.com/dudapinto/pos-golang/rate-limiter/config"
)

type RateLimitData struct {
	LimitPerSecond int
	LastRequest    time.Time
	Requests       int
	Count          int
}

type RateLimiter struct {
	Store  Store
	config *config.Config
	Limits map[string]RateLimitData // Change from limits to Limits (uppercase)
}

func NewLimiter(cfg *config.Config) *RateLimiter {
	return &RateLimiter{
		config: cfg,
		Store:  NewStore(cfg),
		Limits: make(map[string]RateLimitData), // Use the new exported field
	}
}

func NewRateLimiter(store Store, config *config.Config) *RateLimiter {
	return &RateLimiter{
		Store:  store,
		config: config,
		Limits: make(map[string]RateLimitData), // Use the new exported field
	}
}

func (rl *RateLimiter) AllowRequest(key string) bool {
	// Obtém dados do Store
	data, err := rl.Store.Get(key)
	if err != nil {
		// Cria novo limite se não existir
		data = &RateLimit{
			LimitPerSecond: rl.config.LimitPerSecond,
			BlockDuration:  rl.config.BlockDuration,
			LastRequest:    time.Now(),
			Requests:       0,
		}
		if err := rl.Store.Set(key, data); err != nil {
			return false
		}
	}

	// Bloqueio por duração
	if time.Since(data.LastRequest) < data.BlockDuration {
		return false
	}

	// Verificação de limite
	if data.Requests >= data.LimitPerSecond {
		return false
	}

	// Atualiza dados
	data.LastRequest = time.Now()
	data.Requests++
	if err := rl.Store.Set(key, data); err != nil {
		log.Printf("Failed to save data for key %s: %v", key, err)
		return false
	}

	return true
}
