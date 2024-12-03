package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	url         = flag.String("url", "", "URL do serviço a ser testado")
	requests    = flag.Int("requests", 0, "Número total de requests")
	concurrency = flag.Int("concurrency", 0, "Número de chamadas simultâneas")
)

func main() {
	flag.Parse()

	if *url == "" || *requests == 0 || *concurrency == 0 {
		log.Fatal("Parâmetros obrigatórios faltando")
	}

	startTime := time.Now()

	var wg sync.WaitGroup
	var mu sync.Mutex

	statusCodes := make(map[int]int)
	totalRequests := 0

	// Exibe o título e limpa a tela
	fmt.Print("\033[H\033[2J") // Limpa o terminal
	fmt.Println("===================================")
	fmt.Println("Iniciando o Stress Test")
	fmt.Println("===================================")
	fmt.Printf("Testando URL: %s\n", *url)
	fmt.Printf("Número de requisições: %d\n", *requests)
	fmt.Printf("Concorrência: %d\n", *concurrency)
	fmt.Println("===================================")

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < *requests/(*concurrency); j++ {
				resp, err := http.Get(*url)
				if err != nil {
					log.Printf("Erro ao realizar request: %v", err)
					continue
				}
				defer resp.Body.Close()

				// Atualiza o status code e o contador total
				statusCode := resp.StatusCode
				mu.Lock()
				statusCodes[statusCode]++
				totalRequests++

				// Exibe a evolução do contador de requests a cada 10 requests
				if totalRequests%10 == 0 {
					fmt.Printf("\rRequests realizados: %d/%d", totalRequests, *requests)
				}
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Finaliza o teste e exibe os resultados
	endTime := time.Now()
	totalTime := endTime.Sub(startTime)

	fmt.Println("\n===================================")
	fmt.Printf("Tempo total gasto: %v\n", totalTime)
	fmt.Printf("Quantidade total de requests realizados: %d\n", totalRequests)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", statusCodes[200])
	for statusCode, count := range statusCodes {
		if statusCode != 200 {
			fmt.Printf("Quantidade de requests com status HTTP %d: %d\n", statusCode, count)
		}
	}

	// Exibe a finalização do teste
	log.Println("Teste finalizado...")
	os.Exit(0)
}
