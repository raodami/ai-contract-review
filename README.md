# AI Contract Review

AI-powered contract review and risk analysis platform for SMBs and freelancers.

## Tech Stack
- **Backend**: Go + Gin + SQLite (modernc)
- **Frontend**: Next.js 14 + TypeScript + Stripe UI
- **AI Service**: DeepSeek (NLP analysis)

## Quick Start

### 1. Clone and Setup
```bash
cd D:/ai-contract-review
cp .env.example .env
# Add your API keys to .env
```

### 2. Backend
```bash
go mod tidy
go test ./...
go build -o bin/server ./cmd/server
./bin/server
```

### 3. Frontend
```bash
cd web
npm install
npm run dev
```

Open http://localhost:3000

### 4. Docker
```bash
docker-compose up -d
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/auth/register | Register new user |
| POST | /api/auth/login | Login |
| GET | /api/auth/me | Get current user |
| GET | /api/user/usage | Check usage quota |
| POST | /api/contracts/upload | Upload contract |
| GET | /api/contracts/{id} | Get review result |

## Pricing
- **Free**: 3 reviews
- **Pro**: $19.9/mo (unlimited)
- **Business**: $49.9/mo (team + API)

## Environment Variables
- `DEEPSEEK_API_KEY` — AI contract analysis
- `STRIPE_SECRET_KEY` — Payment processing
- `ADMIN_KEY` — Admin authentication

## Project Structure
\`\`\`
ai-contract-review/
├── cmd/server/           # Main entry point
├── internal/
│   ├── api/             # HTTP handlers & routes
│   ├── llm/             # DeepSeek client
│   ├── nlp/             # PDF/DOCX parsing
│   ├── store/           # SQLite database
│   └── webhook/         # Stripe webhooks
├── web/                 # Next.js frontend
├── docs/PRD.md          # Product requirements
├── Dockerfile
└── render.yaml          # Render deployment config
\`\`\`
