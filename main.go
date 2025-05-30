package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"desafio-b3/models"
	"desafio-b3/utils"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	processar := flag.Bool("processar", false, "Processar a pasta cotacoes antes de iniciar a API")
	flag.Parse()

	if *processar {
		fmt.Println("Iniciando processamento dos CSVs...")
		utils.ProcessarPasta("cotacoes/")
		fmt.Println("Processamento concluído!")
	}

	db := utils.ConectarBanco()

	http.HandleFunc("/dados", func(w http.ResponseWriter, r *http.Request) {
		ticker := r.URL.Query().Get("ticker")
		data := r.URL.Query().Get("data")

		var negocios []models.Negocio
		query := db.Where("codigo_instrumento = ?", ticker)
		if data != "" {
			query = query.Where("data_negocio = ?", data)
		}
		query.Find(&negocios)

		var maxRange float64
		volumeDiario := make(map[string]int)
		for _, n := range negocios {
			if n.PrecoNegocio > maxRange {
				maxRange = n.PrecoNegocio
			}
			volumeDiario[n.DataNegocio] += n.QuantidadeNegociada
		}

		var maxVolume int
		for _, v := range volumeDiario {
			if v > maxVolume {
				maxVolume = v
			}
		}

		resp := map[string]interface{}{
			"ticker":           ticker,
			"max_range_value":  maxRange,
			"max_daily_volume": maxVolume,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("API disponível em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
