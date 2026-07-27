set -e

cleanup() {
  echo "Stopping database..."
  docker compose stop db
}
trap cleanup EXIT INT TERM

docker compose up -d db
docker compose run --rm migrate

set -a
. ./.env
export DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${DATABASE_PORT}/${POSTGRES_DB}?sslmode=disable"
set +a

air