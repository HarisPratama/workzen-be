# workzen

Docker command

- docker compose up -d --build app-dev
- docker compose --env-file .env.dev up --build
- docker exec -it workzen-db psql -U postgres
- docker compose -f docker-compose.yml -f docker-compose.local.yml up -d

---

**Last Deployment:** 2026-03-26 13:52 UTC (CI/CD Test)

**Status:** ✅ GitHub Actions CI/CD Integration Active
- Auto-deploy on push to `main` or `develop`
- Local Docker build on VPS
- Telegram notifications enabled
