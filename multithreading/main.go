package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
}

type BrasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ApiResponse struct {
	Response interface{}
	Time     time.Duration
	Source   string
}

func consultarAPI(ctx context.Context, apiURL, cep, source string, ch chan<- ApiResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	apiForRequest := apiURL + cep
	via := strings.Index(apiURL, "viacep")
	if via != -1 {
		apiForRequest += "/json/"
	}

	start := time.Now()

	req, err := http.NewRequest("GET", apiForRequest, nil)
	if err != nil {
		fmt.Println("Erro ao criar requisição:", err)
		return
	}

	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer requisição:", err)
		return
	}
	defer resp.Body.Close()

	var response interface{}
	if via == -1 { // BrasilAPI
		var brasilAPIResp BrasilAPIResponse
		err = json.NewDecoder(resp.Body).Decode(&brasilAPIResp)
		response = brasilAPIResp
	} else { // ViaCep
		var viaCEPResp ViaCEPResponse
		err = json.NewDecoder(resp.Body).Decode(&viaCEPResp)
		response = viaCEPResp
	}

	if err != nil {
		fmt.Println("Erro ao decodificar resposta:", err)
		return
	}

	elapsed := time.Since(start)
	ch <- ApiResponse{
		Response: response,
		Time:     elapsed,
		Source:   source,
	}
	fmt.Println("Resposta recebida de", source)
}

func main() {
	cep := "01153000"
	api1 := "https://brasilapi.com.br/api/cep/v1/"
	api2 := "https://viacep.com.br/ws/"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := make(chan ApiResponse)
	var wg sync.WaitGroup

	// ============================================================================
	// 	ATENÇÃO: Se ao testar, vc comentar uma das APIs, não se esqueça de
	//           alterar tb wg.Add(2) <== para 1 ou wg.wait() nunca vai completar
	//	  		 e vai parecer que travou e não houve response, ok?  ;)
	// ============================================================================
	wg.Add(2)

	go consultarAPI(ctx, api1, cep, "BrasilAPI", ch, &wg)
	go consultarAPI(ctx, api2, cep, "ViaCEP", ch, &wg)

	go func() {
		wg.Wait()
		close(ch)
	}()

	var fastestAPI string
	var fastestTime time.Duration
	var fastestResponse interface{}

	for resp := range ch {
		fmt.Printf("Recebido de %s em %s\n", resp.Source, resp.Time)
		if fastestTime == 0 || resp.Time < fastestTime {
			fastestTime = resp.Time
			fastestAPI = resp.Source
			fastestResponse = resp.Response
		}
	}

	if fastestAPI != "" {
		fmt.Printf("API mais rápida: %s\nTempo de resposta: %s\nDados respondidos: %+v\n", fastestAPI, fastestTime, fastestResponse)
	} else {
		fmt.Println("Timeout ou erro ao consultar as APIs.")
	}
}
