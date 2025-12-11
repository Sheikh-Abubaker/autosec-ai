# AutoSec AI

AutoSec AI — Project for the AI Agents Assemble hackathon.

This repo contains:
- Go backend (API endpoints `/api/scan` and `/api/autofix-plan`)
- Kestra workflows for driving Syft/Grype/AI scans
- Docker Compose for local development
- Minimal frontend planned for Vercel deployment

## Quickstart (local)
1. `docker compose up --build -d`
2. `curl -X POST http://localhost:8080/api/scan -H "Content-Type: application/json" -d '{"repo_url":"https://github.com/example/repo.git"}'`
