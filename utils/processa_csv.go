package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"desafio-b3/models"
)

func ProcessarPasta(caminhoPasta string) {
	arqs, err := os.ReadDir("cotacoes/")
	if err != nil {
		fmt.Println("Erro ao ler a pasta:", err)
		return
	}

	var wg sync.WaitGroup

	for _, arq := range arqs {
		if !arq.IsDir() {
			wg.Add(2)
			caminhoArquivo := filepath.Join("cotacoes/", arq.Name())
			fmt.Println("Processando arquivo:", caminhoArquivo)
			ProcessarCSV(caminhoArquivo)
		}
	}
	wg.Wait()

}

func ProcessarCSV(caminhoArquivo string) {
	db := ConectarBanco()

	arquivo, err := os.Open(caminhoArquivo)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer arquivo.Close()

	scanner := bufio.NewScanner(arquivo)
	var linhas [][]string

	for scanner.Scan() {
		linha := scanner.Text()
		valores := strings.Split(linha, ";")
		linhas = append(linhas, valores)
	}

	if len(linhas) > 0 {
		linhas = linhas[1:]
	}

	var negocios []models.Negocio
	for _, l := range linhas {
		if len(l) < 10 {
			continue
		}

		preco := converteFloat(l[3])
		qtd := converteInt(l[4])

		negocio := models.Negocio{
			HoraFechamento:      l[5],
			DataNegocio:         l[8],
			CodigoInstrumento:   l[1],
			PrecoNegocio:        preco,
			QuantidadeNegociada: qtd,
		}
		negocios = append(negocios, negocio)

	}

	if len(negocios) > 0 {
		db.CreateInBatches(negocios, 1000) // mil registros por lote
	}
}

func converteFloat(valor string) float64 {
	valor = strings.Replace(valor, ",", ".", 1)
	f, _ := strconv.ParseFloat(valor, 64)
	return f
}

func converteInt(valor string) int {
	i, _ := strconv.Atoi(valor)
	return i
}
