set -e

cleanup() {
  echo "Stopping database..."
  docker compose stop db
  docker compose stop redis
}
trap cleanup EXIT INT TERM

docker compose up -d db
docker compose run --rm migrate
docker compose up -d redis

set -a
. ./.env
export DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${DATABASE_PORT}/${POSTGRES_DB}?sslmode=disable"
export REDIS_ADDR="localhost:${REDIS_PORT}"
set +a

air