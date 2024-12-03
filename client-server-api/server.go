package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3" // Driver SQLite
	"github.com/valyala/fastjson"
)

const (
	url     = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	timeout = 200 * time.Millisecond
	dbPath  = "cotacoes.db" // Arquivo do banco de dados SQLite
)

func createTable(db *sql.DB) error {
	// Cria a tabela cotacoes
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS cotacoes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			rate TEXT
		)
	`)
	return err
}

func insertExchangeRate(ctx context.Context, db *sql.DB, rate string) error {
	// Usa o contexto com prazo de 10 ms para a operação no banco de dados
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	_, err := db.ExecContext(ctx, "INSERT INTO cotacoes (rate) VALUES (?)", rate)
	return err
}

func fetchExchangeRate(ctx context.Context, db *sql.DB) (string, error) {
	startGet := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta em um slice de bytes
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Faz o parse da resposta usando fastjson
	parser := fastjson.Parser{}
	value, err := parser.ParseBytes(body)
	if err != nil {
		return "", err
	}

	// Extrai o valor "bid"
	bid := value.GetStringBytes("USDBRL", "bid")
	elapsedGet := time.Since(startGet)
	log.Printf("Tempo de execução do GET Bid: %s", elapsedGet)

	// Insere a cotação no banco de dados
	startInsert := time.Now()
	if err := insertExchangeRate(ctx, db, string(bid)); err != nil {
		log.Printf("Erro ao inserir no banco de dados: %v", err)
	}
	elapsedInsert := time.Since(startInsert)
	log.Printf("Tempo de execução do INSERT Bid: %s", elapsedInsert)
	log.Println("------------------------------------------")

	return string(bid), nil
}

func exchangeRateHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	rate, err := fetchExchangeRate(ctx, db)
	if err != nil {
		http.Error(w, "Erro ao buscar a cotação", http.StatusInternalServerError)
		return
	}

	// Define o tipo de conteúdo da resposta como JSON
	w.Header().Set("Content-Type", "application/json")
	// Escreve a resposta JSON
	fmt.Fprintf(w, `{"bid": "%s"}`, rate)
}

func main() {
	// Inicializa o banco de dados SQLite
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Erro ao abrir o banco de dados: %v", err)
	}
	defer db.Close()

	// Cria a tabela cotacoes
	if err := createTable(db); err != nil {
		log.Fatalf("Erro ao criar a tabela: %v", err)
	}

	// Registra o handler
	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		exchangeRateHandler(w, r, db)
	})

	log.SetOutput(os.Stdout)
	log.Println("Servidor ativo na porta :8080")
	log.Println("------------------------------------------")

	// Inicia o servidor
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Erro ao iniciar o servidor: %v", err)
	}
}
