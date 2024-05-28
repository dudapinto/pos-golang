package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func fetchCotacao() (string, error) {
	url := "http://localhost:8080/cotacao" // Substituir pelo endpoint real de onde for publicado.

	// Crie um contexto com timeout de 300 ms
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("A API retornou um código de status diferente de 200: %d", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func writeToTextFile(data string) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	cotacao, err := fetchCotacao()
	if err != nil {
		log.Fatal("Erro ao buscar a cotação:", err)
	}

	err = writeToTextFile(cotacao)
	if err != nil {
		log.Fatal("Erro ao gravar em cotacao.txt:", err)
	}

	fmt.Println("Cotação gravada em cotacao.txt com sucesso!")
}
