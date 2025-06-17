package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"desafio-b3/models"
	"desafio-b3/utils"
)

type RespostaAPI struct {
	Ticker                string  `json:"ticker"`
	MaxRangeValue         float64 `json:"max_range_value"`
	MaxDailyVolume        int     `json:"max_daily_volume"`
	DateMaxDailyVolume    string  `json:"date_max_daily_volume"`
	TotalNegotiatedVolume int     `json:"total_negotiated_volume"`
}

func main() {
	db := utils.ConectarBanco()

	http.HandleFunc("/dados", func(w http.ResponseWriter, r *http.Request) {
		ticker := r.URL.Query().Get("ticker")
		data := r.URL.Query().Get("data")

		if ticker == "" {
			http.Error(w, "Ticker não fornecido", http.StatusBadRequest)
			return
		}

		var negocios []models.Negocio
		query := db.Where("codigo_instrumento = ?", ticker)
		if data != "" {
			query = query.Where("data_negocio = ?", data)
		}
		query.Find(&negocios)

		var maxPrice float64
		err := db.Model(&models.Negocio{}).
			Where("codigo_instrumento = ?", ticker).
			Select("MAX(preco_negocio)").
			Scan(&maxPrice).Error
		if err != nil {
			http.Error(w, "Erro ao consultar preço máximo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var volumeData struct {
			DataNegocio string
			VolumeTotal int
		}
		err = db.Raw(`
			SELECT data_negocio, SUM(quantidade_negociada) AS volume_total
			FROM negocios
			WHERE codigo_instrumento = ?
			GROUP BY data_negocio
			ORDER BY volume_total DESC
			LIMIT 1
		`, ticker).Scan(&volumeData).Error
		if err != nil {
			http.Error(w, "Erro ao consultar volume máximo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp := map[string]interface{}{
			"ticker":           ticker,
			"max_range_value":  maxPrice,
			"max_daily_volume": volumeData.VolumeTotal,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("API disponível em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
