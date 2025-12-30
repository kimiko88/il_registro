# Registro Elettronico Backend

Backend in Go per il registro elettronico scolastico.

## Struttura
Microservizi modulari (monorepo) con architettura Clean Architecture / Domain-Driven Design layout.

- `cmd/`: Entry points
- `internal/`: Business logic per dominio (auth, users, grades, etc.)
- `pkg/`: Librerie condivise
- `docker/`: Configurazioni Docker

## Setup

1. **Requisiti**: Go 1.21+, Docker.
2. **Env**: Copia `.env.example` in `.env`.
3. **Run**: 
   - Locale: `make run`
   - Docker: `make docker-run`

## Testing
`make test`
