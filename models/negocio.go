package models

import "gorm.io/gorm"

type Negocio struct {
	gorm.Model
	HoraFechamento      string  `json:"hora_fechamento"`
	DataNegocio         string  `json:"data_negocio"`
	CodigoInstrumento   string  `json:"codigo_instrumento"`
	PrecoNegocio        float64 `json:"preco_negocio"`
	QuantidadeNegociada int     `json:"quantidade_negociadas"`
}
