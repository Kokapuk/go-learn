set -e

api_pid=""
worker_pid=""

cleanup() {
  echo "Stopping worker and local services..."
  
  [ -n "$api_pid" ] && kill "$api_pid" 2>/dev/null || true
  [ -n "$worker_pid" ] && kill "$worker_pid" 2>/dev/null || true
  [ -n "$api_pid" ] && wait "$api_pid" 2>/dev/null || true
  [ -n "$worker_pid" ] && wait "$worker_pid" 2>/dev/null || true
  
  docker compose stop db
  docker compose stop redis
  docker compose stop rabbitmq
}
trap cleanup EXIT INT TERM

docker compose up -d --wait db
docker compose run --rm migrate
docker compose up -d --wait redis rabbitmq

set -a
. ./.env
export DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${DATABASE_PORT}/${POSTGRES_DB}?sslmode=disable"
export REDIS_ADDR="localhost:${REDIS_PORT}"
export RABBITMQ_URL="amqp://${RABBITMQ_DEFAULT_USER}:${RABBITMQ_DEFAULT_PASS}@localhost:${RABBITMQ_PORT}"
set +a

air -c .air.worker.toml &
worker_pid=$!

air -c .air.api.toml &
api_pid=$!

wait "$api_pid"
