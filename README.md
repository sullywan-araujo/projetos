# Projeto de Processamento de Cotações B3

Este projeto processa arquivos de cotações históricas da B3 e expõe uma API para consulta de dados relevantes.

## Tecnologias
- Golang
- SQLite (via GORM)

## Como rodar
1. Clone o repositório:
   git clone https://github.com/sullywan-araujo/projetos.git

2. Instale as dependências (se necessário).

3. Execute:
   go run main.go

## Rotas da API
- GET /dados?ticker=PETR4&data=20250101

Retorna:
```json
{
  "ticker": "PETR4",
  "max_range_value": 30.5,
  "max_daily_volume": 1000000
}
