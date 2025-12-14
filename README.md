# AutoSec AI
# Project for the AI Agents Assemble hackathon.

AutoSec AI — Autonomous Supply Chain Security Agent

An AI agent that scans GitHub repos and finds vulnerabilities

This repo contains:
- Go backend (API endpoints `/api/scan` and `/api/autofix-plan`)
- Kestra workflows for driving Syft/Grype/AI scans
- Docker Compose for local development
- Next.js based frontend Vercel deployment

## Quickstart (local)
1. `docker compose up --build -d`
2. Run an API scan:
   ```bash
    curl -X POST http://localhost:8080/api/scan \
    -H "Content-Type: application/json" \
    -d '{"repo_url":"https://github.com/example/repo.git"}`
   ```
