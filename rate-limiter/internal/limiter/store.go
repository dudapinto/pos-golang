// store.go
package limiter

import (
	"fmt"
	"time"

	"github.com/dudapinto/pos-golang/rate-limiter/config"
)

type RateLimit struct {
	LimitPerSecond int           // Limite de requisições por segundo
	BlockDuration  time.Duration // Duração do bloqueio
	LastRequest    time.Time     // Data da última requisição
	Requests       int           // Número de requisições feitas
	Store          Store         // Armazenamento que será utilizado para persistir os dados
}

type Store interface {
	Get(key string) (*RateLimit, error)
	Set(key string, rl *RateLimit) error
	Increment(key string) error
}

// NewStore cria uma nova instância de armazenamento com base no tipo configurado
func NewStore(cfg *config.Config) Store {
	switch cfg.StorageType {
	case "redis":
		return NewRedisStore(cfg.RedisAddr, cfg.RedisPassword)
	case "mock":
		// Aqui configuramos a MockStore com funções de mock
		return &MockStore{
			GetFn: func(key string) (*RateLimit, error) {
				// Implementação mockada de Get
				fmt.Println("Mocking Get function")
				return nil, nil // Simula que não encontrou nada
			},
			SetFn: func(key string, rl *RateLimit) error {
				// Implementação mockada de Set
				fmt.Println("Mocking Set function")
				return nil // Simula sucesso no set
			},
			IncrementFn: func(key string) error {
				// Implementação mockada de Increment
				fmt.Println("Mocking Increment function")
				return nil // Simula sucesso no incremento
			},
		}
	default:
		fmt.Println("Tipo de armazenamento não suportado")
		return nil
	}
}
