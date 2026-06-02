# OhayouBot

WhatsApp morning message bot with AI-generated content.

## Stack
- **Backend**: Go, sqlc, golang-migrate
- **Frontend**: Next.js 14, TypeScript, shadcn/ui
- **Database**: PostgreSQL
- **Infra**: Docker, Nginx

## Getting Started

```bash
cp apps/api/.env.example apps/api/.env
# edit .env with your values
make up
```

## Commands
| Command | Description |
|---|---|
| `make up` | Start all services |
| `make down` | Stop all services |
| `make migrate` | Run DB migrations |
| `make sqlc` | Generate DB code |
