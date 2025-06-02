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

	"gorm.io/gorm"
)

func ProcessarPasta(db *gorm.DB, pasta string) error {
	var wg sync.WaitGroup

	err := filepath.Walk(pasta, func(caminho string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			wg.Add(1)
			go func(caminhoArquivo string) {
				defer wg.Done()
				ProcessarCSV(caminhoArquivo)
			}(caminho)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Erro ao percorrer a pasta:", err)
	}

	wg.Wait()

	if err != nil {
		fmt.Println("Erro ao percorrer a pasta:", err)
		return err
	}

	return nil
}

func ProcessarCSV(caminhoArquivo string) error {
	db := ConectarBanco()

	arquivo, err := os.Open(caminhoArquivo)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", caminhoArquivo, "-", err)
		return err
	}
	defer arquivo.Close()

	scanner := bufio.NewScanner(arquivo)
	var negocios []models.Negocio

	scanner.Scan()

	for scanner.Scan() {
		linha := scanner.Text()
		campos := strings.Split(linha, ";")
		if len(campos) < 10 {
			continue
		}

		preco := converteFloat(campos[3])
		qtd := converteInt(campos[4])

		negocio := models.Negocio{
			HoraFechamento:      campos[5],
			DataNegocio:         campos[8],
			CodigoInstrumento:   campos[1],
			PrecoNegocio:        preco,
			QuantidadeNegociada: qtd,
		}
		negocios = append(negocios, negocio)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro ao ler o arquivo:", caminhoArquivo, "-", err)
		return err
	}

	if len(negocios) > 0 {
		result := db.CreateInBatches(negocios, 1000)
		if result.Error != nil {
			fmt.Println("Erro ao inserir registros no banco:", result.Error)
			return result.Error
		}
		fmt.Printf("Arquivo processado: %s - Registros inseridos: %d\n", caminhoArquivo, len(negocios))
	} else {
		fmt.Printf("Arquivo processado: %s - Nenhum registro válido encontrado.\n", caminhoArquivo)
	}

	return nil
}

func converteFloat(valor string) float64 {
	valor = strings.ReplaceAll(valor, ",", ".")
	f, _ := strconv.ParseFloat(valor, 64)
	return f
}

func converteInt(valor string) int {
	i, _ := strconv.Atoi(valor)
	return i
}
