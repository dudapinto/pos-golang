// main.go
package main

import (
	"log"
	"net/http"

	"github.com/dudapinto/pos-golang/rate-limiter/config"
	"github.com/dudapinto/pos-golang/rate-limiter/internal/limiter"
	"github.com/dudapinto/pos-golang/rate-limiter/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Carregar a configuração
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	// Verificar o que foi carregado
	log.Printf("Loaded config: STORAGE_TYPE = %s", cfg.StorageType)

	// Criar um ponteiro para a configuração
	cfgPtr := &cfg

	log.Printf("STORAGE_TYPE antes de chamar o rate limiter: %s", cfg.StorageType)
	// Criar um novo rate limiter
	rl := limiter.NewLimiter(cfgPtr)

	// Criar um novo router
	r := chi.NewRouter()

	// Adicionar o middleware de rate limit
	r.Use(middleware.RateLimitMiddleware(rl))

	// Adicionar as rotas
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	// Iniciar o servidor
	log.Fatal(http.ListenAndServe(":8080", r))
}
