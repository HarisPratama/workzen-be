# workzen

Docker command

- docker compose up -d --build app-dev
- docker compose --env-file .env.dev up --build
- docker exec -it workzen-db psql -U postgres
- docker compose -f docker-compose.yml -f docker-compose.local.yml up -d